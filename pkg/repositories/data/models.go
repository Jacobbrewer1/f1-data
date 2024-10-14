package data

type PaginationResponse[T comparable] struct {
	Items []*T  `json:"items"`
	Total int64 `json:"total"`
}
