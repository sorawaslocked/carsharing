package dto

import (
	"database/sql"
	"fmt"
	"github.com/sorawaslocked/car-rental-car-service/internal/model"
)

type scanner interface {
	Scan(dest ...any) error
}

type ArgsBuilder struct {
	Args []any
}

func newArgsBuilder() *ArgsBuilder {
	return &ArgsBuilder{
		Args: []any{},
	}
}

func (b *ArgsBuilder) Add(arg any) string {
	b.Args = append(b.Args, arg)

	return fmt.Sprintf("$%d", len(b.Args))
}

func BuildPagination(b *ArgsBuilder, p model.Pagination) string {
	clause := ""

	if p.Limit != nil {
		clause += " LIMIT " + b.Add(*p.Limit)
	}
	if p.Offset != nil {
		clause += " OFFSET " + b.Add(*p.Offset)
	}

	return clause
}

func column(tableAlias, name string) string {
	if tableAlias == "" {
		return name
	}

	return fmt.Sprintf("%s.%s", tableAlias, name)
}

func NullableFloat32(v *float32) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: float64(*v), Valid: true}
}
