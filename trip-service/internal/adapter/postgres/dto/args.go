package dto

import (
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
