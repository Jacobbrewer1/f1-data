package pagefilter

// PaginatedResponse is a response that contains a list of items and the total number of items. This is useful for
// when you want to return a paginated list of items.
type PaginatedResponse[T comparable] struct {
	// Items is the list of items that are returned.
	Items []*T `json:"items"`

	// Total is the total number of items that have been found in the source.
	Total int64 `json:"total"`
}
