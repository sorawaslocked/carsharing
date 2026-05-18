package dto

import (
	"database/sql"
	"fmt"

	"carsharing/trip-service/internal/model"
)

type ArgsBuilder struct {
	Args []any
}

func (b *ArgsBuilder) Add(arg any) string {
	b.Args = append(b.Args, arg)
	return fmt.Sprintf("$%d", len(b.Args))
}

func BuildPagination(p *model.Pagination, b *ArgsBuilder) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf(" LIMIT %s OFFSET %s", b.Add(p.Limit), b.Add(p.Offset))
}

func NullableFloat32(v *float32) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: float64(*v), Valid: true}
}

func NullableString(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *v, Valid: true}
}

func NullableInt64(v *int64) sql.NullInt64 {
	if v == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *v, Valid: true}
}
