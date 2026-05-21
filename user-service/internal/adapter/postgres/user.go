package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	sharedmodel "carsharing/shared/model"
	pkglog "carsharing/shared/pkg/log"
	"carsharing/shared/pkg/utils"
	pgdto "carsharing/user-service/internal/adapter/postgres/dto"
	"carsharing/user-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

func NewUserRepository(log *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		log:  pkglog.WithComponent(log, "repo.UserRepository"),
		pool: pool,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

const userSelect = `
    SELECT u.id, u.email, u.phone_number, u.first_name, u.last_name, u.birth_date,
           u.password_hash, u.profile_image_key,
           ARRAY(SELECT role_name FROM user_roles WHERE user_id = u.id ORDER BY role_name) AS roles,
           u.is_document_verified, u.is_email_verified, u.is_suspended,
           u.created_at, u.updated_at
    FROM users u`

func scanUser(rs rowScanner) (model.User, error) {
	var u model.User
	var phoneNumber *string
	var profileImageKey *string
	var roleStrings []string

	if err := rs.Scan(
		&u.ID, &u.Email, &phoneNumber, &u.FirstName, &u.LastName, &u.BirthDate,
		&u.PasswordHash, &profileImageKey, &roleStrings,
		&u.IsDocumentVerified, &u.IsEmailVerified, &u.IsSuspended,
		&u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return model.User{}, err
	}

	u.PhoneNumber = phoneNumber
	if profileImageKey != nil {
		u.ProfileImage = &model.Image{Key: *profileImageKey}
	}

	u.Roles = make([]sharedmodel.Role, len(roleStrings))
	for i, s := range roleStrings {
		u.Roles[i] = sharedmodel.Role(s)
	}

	return u, nil
}

func (r *UserRepository) handlePGErr(logger *slog.Logger, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.ConstraintName {
		case "users_email_key":
			return model.ErrDuplicateEmail
		case "users_phone_number_key":
			return model.ErrDuplicatePhone
		}
	}
	logger.Error("unexpected postgres error", pkglog.Err(err))
	return model.ErrSql
}

func (r *UserRepository) Insert(ctx context.Context, user model.User) (string, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Insert"), utils.MetadataFromCtx(ctx))

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		logger.Error("beginning transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}
	defer tx.Rollback(ctx)

	var profileImageKey *string
	if user.ProfileImage != nil {
		profileImageKey = &user.ProfileImage.Key
	}

	var id string
	err = tx.QueryRow(ctx, `
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
		return "", r.handlePGErr(logger, err)
	}

	for _, role := range user.Roles {
		if _, err = tx.Exec(ctx,
			`INSERT INTO user_roles (user_id, role_name) VALUES ($1, $2)`,
			id, role.String(),
		); err != nil {
			logger.Error("inserting user role", slog.String("role", role.String()), pkglog.Err(err))
			return "", model.ErrSql
		}
	}

	if err := tx.Commit(ctx); err != nil {
		logger.Error("committing transaction", pkglog.Err(err))
		return "", model.ErrSqlTransaction
	}

	return id, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (model.User, error) {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "FindByID"), utils.MetadataFromCtx(ctx))

	user, err := scanUser(r.pool.QueryRow(ctx, userSelect+" WHERE u.id = $1", id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	user, err := scanUser(r.pool.QueryRow(ctx, query, args...))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

	rows, err := r.pool.Query(ctx, query, args...)
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

func (r *UserRepository) Update(ctx context.Context, id string, update model.UserUpdate) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Update"), utils.MetadataFromCtx(ctx))

	setClauses, args, nextArg := pgdto.SetClausesFromUpdate(update)
	hasRoles := len(update.Roles) > 0

	userQuery := "UPDATE users SET " + strings.Join(setClauses, ", ") +
		fmt.Sprintf(" WHERE id = $%d", nextArg)
	userArgs := append(args, id)

	if hasRoles {
		tx, err := r.pool.Begin(ctx)
		if err != nil {
			logger.Error("beginning transaction", pkglog.Err(err))
			return model.ErrSqlTransaction
		}
		defer tx.Rollback(ctx)

		tag, err := tx.Exec(ctx, userQuery, userArgs...)
		if err != nil {
			return r.handlePGErr(logger, err)
		}
		if tag.RowsAffected() == 0 {
			return model.ErrNotFound
		}

		if _, err = tx.Exec(ctx, `DELETE FROM user_roles WHERE user_id = $1`, id); err != nil {
			logger.Error("deleting user roles", pkglog.Err(err))
			return model.ErrSql
		}
		for _, role := range update.Roles {
			if _, err = tx.Exec(ctx,
				`INSERT INTO user_roles (user_id, role_name) VALUES ($1, $2)`,
				id, role.String(),
			); err != nil {
				logger.Error("inserting user role", slog.String("role", role.String()), pkglog.Err(err))
				return model.ErrSql
			}
		}

		if err := tx.Commit(ctx); err != nil {
			logger.Error("committing transaction", pkglog.Err(err))
			return model.ErrSqlTransaction
		}
		return nil
	}

	tag, err := r.pool.Exec(ctx, userQuery, userArgs...)
	if err != nil {
		return r.handlePGErr(logger, err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	logger := pkglog.WithMetadata(pkglog.WithMethod(r.log, "Delete"), utils.MetadataFromCtx(ctx))

	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		logger.Error("deleting user", pkglog.Err(err))
		return model.ErrSql
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}
	return nil
}
