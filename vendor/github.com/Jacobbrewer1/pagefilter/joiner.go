package pagefilter

// Joiner represents something which can provide joins for an SQL query.
type Joiner interface {
	Join() (string, []any)
}
