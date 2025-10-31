package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"log/slog"
	"strings"
)

type UserRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserRepository(log *slog.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		log: log,
		db:  db,
	}
}

func (r *UserRepository) Insert(ctx context.Context, user model.User) (uint64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, model.ErrSqlTransaction
	}
	defer tx.Rollback()

	var userID int64
	err = tx.QueryRowContext(
		ctx,
		`
		INSERT INTO users
		(email, phone_number, first_name, last_name, birth_date, 
		 password_hash, is_active, is_confirmed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`,
		user.Email,
		user.PhoneNumber,
		user.FirstName,
		user.LastName,
		user.BirthDate,
		user.PasswordHash,
		user.IsActive,
		user.IsConfirmed,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&userID)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Constraint == "users_email_key" {
			return 0, model.ErrDuplicateEmail
		}
		if errors.As(err, &pqErr) && pqErr.Constraint == "users_phone_number_key" {
			return 0, model.ErrDuplicateEmail
		}

		return 0, model.ErrSql
	}

	roles := user.Roles
	for _, role := range roles {
		_, err = tx.ExecContext(
			ctx,
			`
			INSERT INTO user_roles
			(user_id, role_id)
			VALUES ($1, $2)`,
			userID,
			uint32(role),
		)
		if err != nil {
			return 0, model.ErrSql
		}
	}

	if tx.Commit() != nil {
		return 0, model.ErrSqlTransaction
	}

	return uint64(userID), nil
}

func (r *UserRepository) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	query := `
        SELECT id, email, phone_number, first_name, last_name, 
               birth_date, password_hash, is_active, is_confirmed,
               created_at, updated_at
        FROM users`

	whereClauses, args := dto.WhereClausesFromFilter(filter, nil, 1)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	var u model.User

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PhoneNumber, &u.FirstName, &u.LastName,
		&u.BirthDate, &u.PasswordHash, &u.IsActive, &u.IsConfirmed,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrNotFound
		}
		return model.User{}, model.ErrSql
	}

	roles, err := r.findRolesForUser(ctx, u.ID)
	if err != nil {
		return model.User{}, model.ErrSql
	}
	u.Roles = roles

	return u, nil
}

func (r *UserRepository) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	query := `
        SELECT id, email, phone_number, first_name, last_name, 
               birth_date, password_hash, is_active, is_confirmed,
               created_at, updated_at
        FROM users
    `

	whereClauses, args := dto.WhereClausesFromFilter(filter, nil, 1)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, model.ErrSql
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User

		err := rows.Scan(
			&u.ID, &u.Email, &u.PhoneNumber, &u.FirstName, &u.LastName,
			&u.BirthDate, &u.PasswordHash, &u.IsActive, &u.IsConfirmed,
			&u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, model.ErrSql
		}

		roles, err := r.findRolesForUser(ctx, u.ID)
		if err != nil {
			return nil, model.ErrSql
		}
		u.Roles = roles

		users = append(users, u)
	}

	if rows.Err() != nil {
		return nil, model.ErrSql
	}

	return users, nil
}

func (r *UserRepository) findRolesForUser(ctx context.Context, userId uint64) ([]model.Role, error) {
	query := `
		SELECT role_id
		FROM user_roles
		WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, model.ErrSql
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role

		err = rows.Scan(&role)
		if err != nil {
			return nil, model.ErrSql
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (r *UserRepository) Update(ctx context.Context, filter model.UserFilter, update model.UserUpdate) error {
	query := "UPDATE users SET "

	setClauses, args, argNumber := dto.SetClausesFromUpdateData(update)
	if len(setClauses) <= 1 && update.Roles == nil {
		return model.ErrNoUpdateFields
	}
	query += strings.Join(setClauses, ", ")

	whereClauses, args := dto.WhereClausesFromFilter(filter, args, argNumber)
	if len(whereClauses) == 0 {
		return model.ErrEmptyFilter
	}
	query += " WHERE " + strings.Join(whereClauses, " AND ")
	query += " RETURNING id"

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return model.ErrSqlTransaction
	}
	defer tx.Rollback()

	var userID uint64
	err = tx.QueryRowContext(ctx, query, args...).Scan(&userID)
	if err != nil {
		var pqErr *pq.Error

		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.ErrNotFound
		case errors.As(err, &pqErr) && pqErr.Constraint == "users_email_key":
			return model.ErrSqlTransaction
		default:
			return model.ErrSql
		}
	}

	if update.Roles != nil {
		query = "DELETE FROM user_roles WHERE user_id = $1"

		_, err = tx.ExecContext(ctx, query, userID)
		if err != nil {
			return model.ErrSql
		}

		for _, role := range *update.Roles {
			query = "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)"

			_, err = tx.ExecContext(ctx, query, userID, uint32(role))
			if err != nil {
				return model.ErrSql
			}
		}
	}

	if tx.Commit() != nil {
		return model.ErrSqlTransaction
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, filter model.UserFilter) error {
	query := `DELETE FROM users`

	whereClauses, args := dto.WhereClausesFromFilter(filter, nil, 1)
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return model.ErrSql
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return model.ErrSql
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return nil
}
