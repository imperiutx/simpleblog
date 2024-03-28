package db

import (
	"simpleblog/util"
	"strings"
)

type Filters struct {
	Page           int32
	PageSize       int32
	Sort           string
	SortSafeValues []string
}

func (f Filters) sortColumn() string {
	for _, sortSafeValue := range f.SortSafeValues {
		if f.Sort == sortSafeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	// A sensible failsafe to help stop a SQL injection attack.
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int32 {
	return f.PageSize
}

func (f Filters) offset() int32 {
	return (f.Page - 1) * f.PageSize
}

func ValidateFilters(v *util.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(util.In(f.Sort, f.SortSafeValues...), "sort", "invalid sort value")
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	NextPage     int `json:"next_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}
	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		NextPage:     page + 1,
		LastPage:     (totalRecords-1)/pageSize + 1,
		TotalRecords: totalRecords,
	}
}
