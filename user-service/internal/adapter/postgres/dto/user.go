package dto

import (
	"fmt"
	"github.com/sorawaslocked/car-rental-user-service/internal/model"
)

func WhereClausesFromFilter(filter model.UserFilter, args []any, argNumber int) ([]string, []any) {
	var whereClauses []string
	if args == nil {
		args = []any{}
	}

	if filter.ID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.id = $%d", argNumber))
		args = append(args, *filter.ID)
		argNumber++
	}
	if filter.Email != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.email = $%d", argNumber))
		args = append(args, *filter.Email)
		argNumber++
	}
	if filter.PhoneNumber != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.phone_number = $%d", argNumber))
		args = append(args, *filter.PhoneNumber)
		argNumber++
	}
	if filter.FirstName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.first_name = $%d", argNumber))
		args = append(args, *filter.FirstName)
		argNumber++
	}
	if filter.LastName != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.last_name = $%d", argNumber))
		args = append(args, *filter.LastName)
		argNumber++
	}
	if filter.IsActive != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_active = $%d", argNumber))
		args = append(args, *filter.IsActive)
		argNumber++
	}
	if filter.IsConfirmed != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("u.is_confirmed = $%d", argNumber))
		args = append(args, *filter.IsConfirmed)
		argNumber++
	}

	return whereClauses, args
}

func SetClausesFromUpdateData(update model.UserUpdateData) ([]string, []any, int) {
	var setClauses []string
	var args []any
	argNumber := 1

	if update.Email != nil {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argNumber))
		args = append(args, *update.Email)
		argNumber++
	}
	if update.PhoneNumber != nil {
		setClauses = append(setClauses, fmt.Sprintf("phone = $%d", argNumber))
		args = append(args, *update.PhoneNumber)
		argNumber++
	}
	if update.FirstName != nil {
		setClauses = append(setClauses, fmt.Sprintf("first_name = $%d", argNumber))
		args = append(args, *update.FirstName)
		argNumber++
	}
	if update.LastName != nil {
		setClauses = append(setClauses, fmt.Sprintf("last_name = $%d", argNumber))
		args = append(args, *update.LastName)
		argNumber++
	}
	if update.BirthDate != nil {
		setClauses = append(setClauses, fmt.Sprintf("birth_date = $%d", argNumber))
		args = append(args, *update.BirthDate)
		argNumber++
	}
	//if update.PasswordHash != nil {
	//	setClauses = append(setClauses, fmt.Sprintf("password_hash = $%d", argNumber))
	//	args = append(args, *update.PasswordHash)
	//	argNumber++
	//}
	if update.IsActive != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_active = $%d", argNumber))
		args = append(args, *update.IsActive)
		argNumber++
	}
	if update.IsConfirmed != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_confirmed = $%d", argNumber))
		args = append(args, *update.IsConfirmed)
		argNumber++
	}
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argNumber))
	args = append(args, update.UpdatedAt)
	argNumber++

	return setClauses, args, argNumber
}
