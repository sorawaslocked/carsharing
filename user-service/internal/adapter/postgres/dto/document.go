package dto

import (
	"fmt"

	"carsharing/user-service/internal/model"
)

func WhereClausesFromDocumentFilter(filter model.DocumentFilter, args []any, argNumber int) ([]string, []any, int) {
	var clauses []string
	if args == nil {
		args = []any{}
	}

	if filter.UserID != nil {
		clauses = append(clauses, fmt.Sprintf("user_id = $%d", argNumber))
		args = append(args, *filter.UserID)
		argNumber++
	}
	if filter.ExcludeStatus != nil {
		clauses = append(clauses, fmt.Sprintf("status != $%d", argNumber))
		args = append(args, filter.ExcludeStatus.String())
		argNumber++
	}

	return clauses, args, argNumber
}

func SetClausesFromDocumentUpdate(update model.DocumentUpdate) ([]string, []any, int) {
	var clauses []string
	var args []any
	argNumber := 1

	if update.Status != nil {
		clauses = append(clauses, fmt.Sprintf("status = $%d", argNumber))
		args = append(args, update.Status.String())
		argNumber++
	}
	if update.Error != nil {
		clauses = append(clauses, fmt.Sprintf("error = $%d", argNumber))
		args = append(args, *update.Error)
		argNumber++
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argNumber))
	args = append(args, update.UpdatedAt)
	argNumber++

	return clauses, args, argNumber
}
