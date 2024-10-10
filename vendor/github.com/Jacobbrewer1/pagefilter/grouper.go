package pagefilter

// Grouper represents something which can provide a group by to filter a sql query.
// Grouper support should only ever be used with a filter that adds the table/column refs
// that it needs to function, otherwise you will likely have a bad time.
type Grouper interface {
	Group() []string
}
