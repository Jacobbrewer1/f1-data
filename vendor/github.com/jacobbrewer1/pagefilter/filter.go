package pagefilter

import "strings"

type Filter interface {
	Joiner
	Wherer
}

type MultiFilter struct {
	joinSQL   strings.Builder
	joinArgs  []any
	whereSQL  strings.Builder
	whereArgs []any
	groupCols []string
}

func NewMultiFilter() *MultiFilter {
	return &MultiFilter{
		joinArgs:  make([]any, 0),
		whereArgs: make([]any, 0),
		groupCols: make([]string, 0),
	}
}

func (m *MultiFilter) Add(f any) {
	if j, ok := f.(Joiner); ok {
		joinSQL, joinArgs := j.Join()
		if joinArgs == nil {
			joinArgs = make([]any, 0)
		}
		m.joinSQL.WriteString(strings.TrimSpace(joinSQL))
		m.joinSQL.WriteString("\n")
		m.joinArgs = append(m.joinArgs, joinArgs...)
	}

	if wt, okWt := f.(WhereTyper); okWt {
		whereSQL, whereArgs := wt.Where()
		if whereArgs == nil {
			whereArgs = make([]any, 0)
		}
		wtStr := WhereTypeAnd
		if wt.WhereType().IsValid() {
			wtStr = wt.WhereType()
		}
		m.whereSQL.WriteString(string(wtStr) + " ")
		m.whereSQL.WriteString(strings.TrimSpace(whereSQL))
		m.whereSQL.WriteString("\n")
		m.whereArgs = append(m.whereArgs, whereArgs...)
	} else if w, ok := f.(Wherer); ok {
		whereSQL, whereArgs := w.Where()
		if whereArgs == nil {
			whereArgs = []any{}
		}
		m.whereSQL.WriteString(string(WhereTypeAnd) + " ") // default to AND
		m.whereSQL.WriteString(strings.TrimSpace(whereSQL))
		m.whereSQL.WriteString("\n")
		m.whereArgs = append(m.whereArgs, whereArgs...)
	}

	if g, ok := f.(Grouper); ok {
		m.groupCols = append(m.groupCols, g.Group()...)
	}
}

func (m *MultiFilter) Join() (string, []any) {
	return strings.TrimSpace(m.joinSQL.String()), m.joinArgs
}

func (m *MultiFilter) Where() (string, []any) {
	return strings.TrimSpace(m.whereSQL.String()), m.whereArgs
}

func (m *MultiFilter) Group() []string {
	return m.groupCols
}
