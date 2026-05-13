package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lib/pq"
	pgdto "github.com/sorawaslocked/car-rental-user-service/internal/adapter/postgres/dto"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
	pkglog "github.com/sorawaslocked/car-rental-user-service/internal/pkg/log"
	"github.com/sorawaslocked/car-rental-user-service/internal/pkg/utils"
)

type UserRepository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserRepository(log *slog.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		log: pkglog.WithComponent(log, "repo.UserRepository"),
		db:  db,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

// Roles are fetched via a correlated subquery so every SELECT path
// returns them without extra round-trips.
const userSelect = `
    SELECT u.id, u.email, u.phone_number, u.first_name, u.last_name, u.birth_date,
           u.password_hash, u.profile_image_key,
           ARRAY(SELECT role_name FROM user_roles WHERE user_id = u.id ORDER BY role_name) AS roles,
           u.is_document_verified, u.is_email_verified, u.is_suspended,
           u.created_at, u.updated_at
    FROM users u`

func scanUser(rs rowScanner) (model.User, error) {
	var u model.User
	var phoneNumber sql.NullString
	var profileImageKey sql.NullString
	var roleStrings pq.StringArray

	if err := rs.Scan(
		&u.ID, &u.Email, &phoneNumber, &u.FirstName, &u.LastName, &u.BirthDate,
		&u.PasswordHash, &profileImageKey, &roleStrings,
		&u.IsDocumentVerified, &u.IsEmailVerified, &u.IsSuspended,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return model.User{}, err
	}

	if phoneNumber.Valid {
		u.PhoneNumber = &phoneNumber.String
	}
	if profileImageKey.Valid {
		u.ProfileImage = &model.Image{Key: profileImageKey.String}
	}

	u.Roles = make([]model.Role, len(roleStrings))
	for i, s := range roleStrings {
		u.Roles[i] = model.Role(s)
	}

	return u, nil
}

// handlePQErr maps known constraint violations to domain errors. Any other
// failure is logged and returned as ErrSql.
func (r *UserRepository) handlePQErr(logger *slog.Logger, err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Constraint {
		case "users_email_key":
			return model.ErrDuplicateEmail
		case "users_phone_number_key":
			return model.ErrDuplicatePhone
		}
	}
	logger.Error("unexpected sql error", pkglog.Err(err))
	return model.ErrSql
}

func (r *UserRepository) Insert(ctx context.Context, user model.User) (string, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Error("beginning transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}
	defer tx.Rollback()

	var profileImageKey *string
	if user.ProfileImage != nil {
		profileImageKey = &user.ProfileImage.Key
	}

	var id string
	err = tx.QueryRowContext(ctx, `
        INSERT INTO users
            (email, phone_number, first_name, last_name, birth_date, password_hash,
             profile_image_key, is_document_verified, is_email_verified, is_suspended,
             created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id`,
		user.Email,
		user.PhoneNumber,
		user.FirstName,
		user.LastName,
		user.BirthDate,
		user.PasswordHash,
		profileImageKey,
		user.IsDocumentVerified,
		user.IsEmailVerified,
		user.IsSuspended,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return "", r.handlePQErr(logger, err)
	}

	for _, role := range user.Roles {
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO user_roles (user_id, role_name) VALUES ($1, $2)`,
			id, role.String(),
		); err != nil {
			logger.Error("inserting user role", slog.String("role", role.String()), pkglog.Err(err))
			return "", model.ErrSql
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("committing transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}

	return id, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	user, err := scanUser(r.db.QueryRowContext(ctx, userSelect+" WHERE u.id = $1", id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrNotFound
		}
		logger.Error("scanning user", pkglog.Err(err))
		return model.User{}, model.ErrSql
	}
	return user, nil
}

func (r *UserRepository) FindOne(ctx context.Context, filter model.UserFilter) (model.User, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindOne"), utils.MetadataFromCtx(ctx))

	query := userSelect
	clauses, args, _ := pgdto.WhereClausesFromFilter(filter, nil, 1)
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	user, err := scanUser(r.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrNotFound
		}
		logger.Error("scanning user", pkglog.Err(err))
		return model.User{}, model.ErrSql
	}
	return user, nil
}

func (r *UserRepository) Find(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Find"), utils.MetadataFromCtx(ctx))

	query := userSelect
	clauses, args, nextArg := pgdto.WhereClausesFromFilter(filter, nil, 1)
	if len(clauses) > 0 {
		query += " WHERE " + strings.Join(clauses, " AND ")
	}

	if filter.Pagination != nil {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", nextArg, nextArg+1)
		args = append(args, filter.Pagination.Limit, filter.Pagination.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Error("querying users", pkglog.Err(err))
		return nil, model.ErrSql
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			logger.Error("scanning user row", pkglog.Err(err))
			return nil, model.ErrSql
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logger.Error("iterating user rows", pkglog.Err(err))
		return nil, model.ErrSql
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, update model.UserRepoUpdate) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, nextArg := pgdto.SetClausesFromRepoUpdate(update)
	hasRoles := len(update.Roles) > 0

	userQuery := "UPDATE users SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE id = $%d", nextArg)
	userArgs := append(args, id)

	if hasRoles {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			logger.Error("beginning transaction", pkglog.Err(err))
			return model.ErrSqlTransaction
		}
		defer tx.Rollback()

		res, err := tx.ExecContext(ctx, userQuery, userArgs...)
		if err != nil {
			return r.handlePQErr(logger, err)
		}
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			return model.ErrNotFound
		}

		if _, err = tx.ExecContext(ctx, `DELETE FROM user_roles WHERE user_id = $1`, id); err != nil {
			logger.Error("deleting user roles", pkglog.Err(err))
			return model.ErrSql
		}
		for _, role := range update.Roles {
			if _, err = tx.ExecContext(ctx,
				`INSERT INTO user_roles (user_id, role_name) VALUES ($1, $2)`,
				id, role.String(),
			); err != nil {
				logger.Error("inserting user role", slog.String("role", role.String()), pkglog.Err(err))
				return model.ErrSql
			}
		}

		if err := tx.Commit(); err != nil {
			logger.Error("committing transaction", pkglog.Err(err))
			return model.ErrSqlTransaction
		}
		return nil
	}

	res, err := r.db.ExecContext(ctx, userQuery, userArgs...)
	if err != nil {
		return r.handlePQErr(logger, err)
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	res, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		logger.Error("deleting user", pkglog.Err(err))
		return model.ErrSql
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logger.Error("reading rows affected", pkglog.Err(err))
		return model.ErrSql
	}
	if rowsAffected == 0 {
		return model.ErrNotFound
	}
	return nil
}
