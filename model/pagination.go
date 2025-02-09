package model

// Pagination struct
type Pagination struct {
	TotalCount int    `json:"totalCount"`
	PageSize   int    `json:"pageSize" default:"10"`
	PageNumber int    `json:"pageNumber" default:"1"`
	Sort       string `json:"sort" default:"asc"`
	SortBy     string `json:"sortBy"`
}
