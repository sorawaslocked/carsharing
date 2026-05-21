package dto

import (
	"fmt"

	"carsharing/user-service/internal/model"
)

// WhereClausesFromFilter builds parameterised WHERE clauses from a UserFilter.
// Returns the clauses, accumulated args, and the next available argument index.
func WhereClausesFromFilter(filter model.UserFilter, args []any, argNumber int) ([]string, []any, int) {
	var clauses []string
	if args == nil {
		args = []any{}
	}

	if filter.Email != nil {
		clauses = append(clauses, fmt.Sprintf("email = $%d", argNumber))
		args = append(args, *filter.Email)
		argNumber++
	}
	if filter.PhoneNumber != nil {
		clauses = append(clauses, fmt.Sprintf("phone_number = $%d", argNumber))
		args = append(args, *filter.PhoneNumber)
		argNumber++
	}
	if filter.FirstName != nil {
		clauses = append(clauses, fmt.Sprintf("first_name = $%d", argNumber))
		args = append(args, *filter.FirstName)
		argNumber++
	}
	if filter.LastName != nil {
		clauses = append(clauses, fmt.Sprintf("last_name = $%d", argNumber))
		args = append(args, *filter.LastName)
		argNumber++
	}
	if filter.IsDocumentVerified != nil {
		clauses = append(clauses, fmt.Sprintf("is_document_verified = $%d", argNumber))
		args = append(args, *filter.IsDocumentVerified)
		argNumber++
	}
	if filter.IsEmailVerified != nil {
		clauses = append(clauses, fmt.Sprintf("is_email_verified = $%d", argNumber))
		args = append(args, *filter.IsEmailVerified)
		argNumber++
	}
	if filter.IsSuspended != nil {
		clauses = append(clauses, fmt.Sprintf("is_suspended = $%d", argNumber))
		args = append(args, *filter.IsSuspended)
		argNumber++
	}

	return clauses, args, argNumber
}

// SetClausesFromUpdate builds parameterised SET clauses from a UserUpdate.
// Returns the clauses, accumulated args, and the next available argument index.
func SetClausesFromUpdate(update model.UserUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	argNumber := 1

	if update.Email != nil {
		clauses = append(clauses, fmt.Sprintf("email = $%d", argNumber))
		args = append(args, *update.Email)
		argNumber++
	}
	if update.PhoneNumber != nil {
		clauses = append(clauses, fmt.Sprintf("phone_number = $%d", argNumber))
		args = append(args, *update.PhoneNumber)
		argNumber++
	}
	if update.FirstName != nil {
		clauses = append(clauses, fmt.Sprintf("first_name = $%d", argNumber))
		args = append(args, *update.FirstName)
		argNumber++
	}
	if update.LastName != nil {
		clauses = append(clauses, fmt.Sprintf("last_name = $%d", argNumber))
		args = append(args, *update.LastName)
		argNumber++
	}
	if update.BirthDate != nil {
		clauses = append(clauses, fmt.Sprintf("birth_date = $%d", argNumber))
		args = append(args, *update.BirthDate)
		argNumber++
	}
	if len(update.PasswordHash) > 0 {
		clauses = append(clauses, fmt.Sprintf("password_hash = $%d", argNumber))
		args = append(args, update.PasswordHash)
		argNumber++
	}
	if update.ProfileImageKey != nil {
		clauses = append(clauses, fmt.Sprintf("profile_image_key = $%d", argNumber))
		args = append(args, *update.ProfileImageKey)
		argNumber++
	}
	if update.IsDocumentVerified != nil {
		clauses = append(clauses, fmt.Sprintf("is_document_verified = $%d", argNumber))
		args = append(args, *update.IsDocumentVerified)
		argNumber++
	}
	if update.IsEmailVerified != nil {
		clauses = append(clauses, fmt.Sprintf("is_email_verified = $%d", argNumber))
		args = append(args, *update.IsEmailVerified)
		argNumber++
	}
	if update.IsSuspended != nil {
		clauses = append(clauses, fmt.Sprintf("is_suspended = $%d", argNumber))
		args = append(args, *update.IsSuspended)
		argNumber++
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argNumber))
	args = append(args, update.UpdatedAt)
	argNumber++

	return clauses, args, argNumber
}
