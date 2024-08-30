package pagefilter

import (
	"strings"
)

const (
	// Equal is the equal comparison operator
	Equal = "eq"

	// LessThan is the less than comparison operator
	LessThan = "lt"

	// GreaterThan is the greater than comparison operator
	GreaterThan = "gt"

	// Like is the partial match operator
	Like = "like"
)

// Joiner represents something which can provide joins for an SQL query.
type Joiner interface {
	Join() (string, []any)
}

// Wherer represents something which can provide a where to filter an sql query.
type Wherer interface {
	Where() (string, []any)
}

// Filter represents something which can provide joins and wheres for an SQL query.
type Filter interface {
	Joiner
	Wherer
}

// MultiFilter is a utility filter to combine multiple filters.
type MultiFilter struct {
	jSQL  strings.Builder
	jArgs []any
	wSQL  strings.Builder
	wArgs []any
	gCols []string
}

// NewMultiFilter creates a new instance of a multi filter.
func NewMultiFilter() *MultiFilter {
	return &MultiFilter{
		jArgs: make([]any, 0, 8),
		wArgs: make([]any, 0, 8),
		gCols: make([]string, 0, 8),
	}
}

// Add adds to the filter
func (m *MultiFilter) Add(f any) {
	if j, ok := f.(Joiner); ok {
		fjSQL, fjArgs := j.Join()
		m.jSQL.WriteString(strings.TrimSpace(fjSQL))
		m.jSQL.WriteString("\n")
		m.jArgs = append(m.jArgs, fjArgs...)
	}

	if w, ok := f.(Wherer); ok {
		fwSQL, fwArgs := w.Where()
		m.wSQL.WriteString(strings.TrimSpace(fwSQL))
		m.wSQL.WriteString("\n")
		m.wArgs = append(m.wArgs, fwArgs...)
	}

	if g, ok := f.(Grouper); ok {
		m.gCols = append(m.gCols, g.Group()...)
	}
}

// Join provides the join for the sql query
func (m *MultiFilter) Join() (string, []any) {
	return strings.TrimSpace(m.jSQL.String()), m.jArgs
}

// Where provides the where for the sql query
func (m *MultiFilter) Where() (string, []any) {
	return strings.TrimSpace(m.wSQL.String()), m.wArgs
}

// Group provides the group by columns for the sql query
func (m *MultiFilter) Group() []string {
	return m.gCols
}

// Grouper represents something which can provide a group by to filter a sql query.
// Grouper support should only ever be used with a filter that adds the table/column refs
// that it needs to function, otherwise you will likely have a bad time.
type Grouper interface {
	Group() []string
}
