package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/logger"
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
		r.log.Error("beginning transaction", logger.Err(err))

		return 0, err
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
		r.log.Error("inserting user", logger.Err(err))

		return 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		`
		INSERT INTO user_roles
		(user_id, role_id)
		VALUES ($1, $2)`,
		userID,
		int32(user.Role),
	)
	if err != nil {
		r.log.Error("inserting user role", logger.Err(err))

		return 0, err
	}

	return uint64(userID), tx.Commit()
}

func (r *UserRepository) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	query := `
        SELECT u.id, u.email, u.phone_number, u.first_name, u.last_name, 
               u.birth_date, u.password_hash, u.is_active, u.is_confirmed,
               u.created_at, u.updated_at, ur.role_id
        FROM users u
    `

	var whereClauses []string
	var args []any
	argID := 1

	if filter.ID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.id = $%d", argID))
		args = append(args, *filter.ID)
		argID++
	}
	if filter.Email != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.email = $%d", argID))
		args = append(args, *filter.Email)
		argID++
	}
	if filter.PhoneNumber != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.phone_number = $%d", argID))
		args = append(args, *filter.PhoneNumber)
		argID++
	}
	if filter.FirstName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.first_name = $%d", argID))
		args = append(args, *filter.FirstName)
		argID++
	}
	if filter.LastName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.last_name = $%d", argID))
		args = append(args, *filter.LastName)
		argID++
	}
	if filter.IsActive != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argID))
		args = append(args, *filter.IsActive)
		argID++
	}
	if filter.IsConfirmed != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_confirmed = $%d", argID))
		args = append(args, *filter.IsConfirmed)
		argID++
	}
	if filter.Role != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("ur.role_id = $%d", argID))
		args = append(args, int32(*filter.Role))
		argID++
	}

	if len(whereClauses) > 0 {
		query += ` INNER JOIN user_roles ur ON u.id = ur.user_id`
		query += " WHERE " + strings.Join(whereClauses, " AND ") + " LIMIT 1"
	}

	var u model.User

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&u.ID, &u.Email, &u.PhoneNumber, &u.FirstName, &u.LastName,
		&u.BirthDate, &u.PasswordHash, &u.IsActive, &u.IsConfirmed,
		&u.CreatedAt, &u.UpdatedAt, &u.Role,
	)
	if err != nil {
		r.log.Error("finding user", logger.Err(err))

		return model.User{}, err
	}

	return u, nil
}

func (r *UserRepository) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	query := `
        SELECT u.id, u.email, u.phone_number, u.first_name, u.last_name, 
               u.birth_date, u.password_hash, u.is_active, u.is_confirmed,
               u.created_at, u.updated_at, ur.role_id
        FROM users u
    `

	var whereClauses []string
	var args []any
	argID := 1

	if filter.ID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.id = $%d", argID))
		args = append(args, *filter.ID)
		argID++
	}
	if filter.Email != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.email = $%d", argID))
		args = append(args, *filter.Email)
		argID++
	}
	if filter.PhoneNumber != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.phone_number = $%d", argID))
		args = append(args, *filter.PhoneNumber)
		argID++
	}
	if filter.FirstName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.first_name = $%d", argID))
		args = append(args, *filter.FirstName)
		argID++
	}
	if filter.LastName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.last_name = $%d", argID))
		args = append(args, *filter.LastName)
		argID++
	}
	if filter.IsActive != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argID))
		args = append(args, *filter.IsActive)
		argID++
	}
	if filter.IsConfirmed != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_confirmed = $%d", argID))
		args = append(args, *filter.IsConfirmed)
		argID++
	}
	if filter.Role != nil {
		query += ` INNER JOIN user_roles ur ON u.id = ur.user_id`
		whereClauses = append(whereClauses, fmt.Sprintf("ur.role_id = $%d", argID))
		args = append(args, int32(*filter.Role))
		argID++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User

		err := rows.Scan(
			&u.ID, &u.Email, &u.PhoneNumber, &u.FirstName, &u.LastName,
			&u.BirthDate, &u.PasswordHash, &u.IsActive, &u.IsConfirmed,
			&u.CreatedAt, &u.UpdatedAt, &u.Role,
		)
		if err != nil {
			r.log.Error("scanning user row", logger.Err(err))

			return nil, err
		}

		users = append(users, u)
	}

	return users, rows.Err()
}

func (r *UserRepository) Update(ctx context.Context, filter model.UserFilter, update model.UserUpdateData) error {
	query := "UPDATE users SET "
	var args []any
	argID := 1
	var setClauses []string

	if update.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argID))
		args = append(args, *update.Email)
		argID++
	}
	if update.PhoneNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argID))
		args = append(args, *update.PhoneNumber)
		argID++
	}
	if update.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argID))
		args = append(args, *update.FirstName)
		argID++
	}
	if update.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argID))
		args = append(args, *update.LastName)
		argID++
	}
	if update.BirthDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birth_date = $%d", argID))
		args = append(args, *update.BirthDate)
		argID++
	}
	if update.PasswordHash != nil {
		setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", argID))
		args = append(args, *update.PasswordHash)
		argID++
	}
	if update.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argID))
		args = append(args, *update.IsActive)
		argID++
	}
	if update.IsConfirmed != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_confirmed = $%d", argID))
		args = append(args, *update.IsConfirmed)
		argID++
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.log.Error("beginning transaction", logger.Err(err))

		return err
	}
	defer tx.Rollback()

	if update.Role != nil {
		_, err = tx.ExecContext(
			ctx,
			`
			UPDATE user_roles
			SET user_id = $1
				role_id = $2`,
			*filter.ID,
			int32(*update.Role),
		)

		if len(setClauses) == 0 {
			return tx.Commit()
		}
	}

	if len(setClauses) == 0 {
		return ErrNoUpdateFields
	}
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, update.UpdatedAt)
	argID++

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argID)
	args = append(args, *filter.ID)

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		r.log.Error("updating user", logger.Err(err))

		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.log.Error("getting affected rows", logger.Err(err))

		return err
	}

	if rowsAffected == 0 {
		return model.ErrNotFound
	}

	return tx.Commit()
}

func (r *UserRepository) Delete(ctx context.Context, filter model.UserFilter) error {
	return nil
}
