package model

const (
	DefaultPaginationLimit  int64 = 20
	DefaultPaginationOffset int64 = 0
)

type Pagination struct {
	Limit  int64
	Offset int64
}
