package validation

import sharedmodel "carsharing/shared/model"

type Pagination struct {
	Limit  int64 `validate:"min=1,max=200"`
	Offset int64 `validate:"min=0"`
}

func DefaultPagination() *Pagination {
	return &Pagination{Limit: sharedmodel.DefaultPaginationLimit, Offset: sharedmodel.DefaultPaginationOffset}
}
