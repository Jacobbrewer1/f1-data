package pagefilter

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	// ErrNoDestination is returned when the destination is nil
	ErrNoDestination = errors.New("destination is nil")
)

// Paginator is the struct that provides the paging.
type Paginator struct {
	db      DB
	idKey   string
	table   string
	filter  Filter
	details *PaginatorDetails
}

// NewPaginator creates a new paginator
func NewPaginator(db DB, table, idk string, f Filter) *Paginator {
	if f == nil {
		f = NewMultiFilter()
	}
	return &Paginator{
		db:     db,
		idKey:  idk,
		table:  table,
		filter: f,
	}
}

// ParseRequest parses the request to handle retrieving all the pagination and sorting parameters
func (p *Paginator) ParseRequest(req *http.Request, sortColumns ...string) error {
	pd, err := DetailsFromRequest(req)
	if err != nil {
		return err
	}

	return p.SetDetails(pd, sortColumns...)
}

// DetailsFromRequest retrieves the paginator details from the request.
func DetailsFromRequest(req *http.Request) (*PaginatorDetails, error) {
	q := req.URL.Query()

	limit, err := getLimit(q)
	if err != nil {
		return nil, fmt.Errorf("%v: %w", err, ErrInvalidPaginatorDetails)
	}

	return &PaginatorDetails{
		Limit:   limit,
		LastVal: q.Get(QueryLastVal),
		LastID:  q.Get(QueryLastID),
		SortBy:  q.Get(QuerySortBy),
		SortDir: q.Get(QuerySortDir),
	}, nil
}

// SetDetails sets paginator details from the passed in arguments.
func (p *Paginator) SetDetails(paginatorDetails *PaginatorDetails, sortColumns ...string) error {
	p.details = paginatorDetails

	wantedSort := p.details.SortBy
	p.details.SortBy = ""
	if wantedSort != "" {
		for _, v := range sortColumns {
			if v == wantedSort {
				p.details.SortBy = v
				break
			}
		}
		if p.details.SortBy == "" {
			return fmt.Errorf("invalid sort %q", wantedSort)
		}
	}

	if p.details.SortBy == "" {
		// We have no specified sort so use the id key
		p.details.SortBy = p.idKey
	}

	sort := strings.ToLower(p.details.SortDir)
	// Define sql from constants to ensure sql query / user input separation
	switch sort {
	case "", orderAsc:
		p.details.sortComparator = sqlComparatorAsc
		p.details.sortOperator = sqlOperatorAsc
	case orderDesc:
		p.details.sortComparator = sqlComparatorDesc
		p.details.sortOperator = sqlOperatorDesc
	default:
		return fmt.Errorf("invalid sort direction %q", sort)
	}
	return nil
}

// First is used when no details are provided which could give us the pivot point
// It will pick a start depending on the provided sort and filters.
func (p *Paginator) First() (string, error) {
	jSQL, jArgs := p.filter.Join()
	wSQL, wArgs := p.filter.Where()
	var gSQL string
	if g, ok := p.filter.(Grouper); ok && len(g.Group()) > 0 {
		gSQL = fmt.Sprintf("GROUP BY %s", strings.Join(g.Group(), ", "))
	}

	// Be aware of SQL injection if modifying the below SQL. Any parameters in the sprintf
	// MUST not be allowed to be created by external input.
	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("SELECT t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" \n")
	sqlBuilder.WriteString("FROM ")
	sqlBuilder.WriteString(p.table)
	sqlBuilder.WriteString(" t \n")

	if jSQL != "" {
		sqlBuilder.WriteString(jSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("WHERE (1 = 1) \n")
	if wSQL != "" {
		sqlBuilder.WriteString("AND (\n")
		sqlBuilder.WriteString(trimWherePrefix(wSQL))
		sqlBuilder.WriteString("\n)\n")
	}

	if gSQL != "" {
		sqlBuilder.WriteString(gSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("ORDER BY t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" ")
	sqlBuilder.WriteString(p.details.sortComparator)
	sqlBuilder.WriteString(", t.")
	sqlBuilder.WriteString(p.idKey)
	sqlBuilder.WriteString(" ASC \n")
	sqlBuilder.WriteString("LIMIT 1")

	args := append(jArgs, wArgs...)
	var err error
	sql := sqlBuilder.String()
	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		return "", fmt.Errorf("first sql in: %w", err)
	}

	var pivot string
	err = p.db.Get(&pivot, sql, args...)
	if err != nil {
		return "", fmt.Errorf("first select: %w", err)
	}

	return pivot, nil
}

// Pivot finds the pivot point in the data.
func (p *Paginator) Pivot() (string, error) {
	// We were given no information about where to pivot from, pivot from the first value
	if p.details.LastID == "" && p.details.LastVal == "" {
		return p.First()
	}

	jSQL, jArgs := p.filter.Join()
	wSQL, wArgs := p.filter.Where()
	var gSQL string
	if g, ok := p.filter.(Grouper); ok && len(g.Group()) > 0 {
		gSQL = fmt.Sprintf("GROUP BY %s", strings.Join(g.Group(), ", "))
	}

	// Be aware of SQL injection if modifying the below SQL. Any parameters in the sprintf
	// MUST not be allowed to be created by external input.
	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("SELECT t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" \n")
	sqlBuilder.WriteString("FROM ")
	sqlBuilder.WriteString(p.table)
	sqlBuilder.WriteString(" t \n")

	if jSQL != "" {
		sqlBuilder.WriteString(jSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("WHERE (t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" = ? AND t.")
	sqlBuilder.WriteString(p.idKey)
	sqlBuilder.WriteString(" >= ?) \n")

	if wSQL != "" {
		sqlBuilder.WriteString("AND (\n")
		sqlBuilder.WriteString(trimWherePrefix(wSQL))
		sqlBuilder.WriteString("\n)\n")
	}

	if gSQL != "" {
		sqlBuilder.WriteString(gSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("LIMIT 1")

	args := append(jArgs, p.details.LastVal, p.details.LastID)
	args = append(args, wArgs...)

	sql := sqlBuilder.String()
	var err error
	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		return "", fmt.Errorf("pivot sql in: %w", err)
	}

	var pivot string
	err = p.db.Get(&pivot, sql, args...)
	if err != nil {
		return "", fmt.Errorf("pivot select: %w", err)
	}

	return pivot, nil
}

// Retrieve pulls the next page given the pivot point and requires a destination *[]struct to load the data into.
func (p *Paginator) Retrieve(pivot string, dest any) error {
	if dest == nil {
		return ErrNoDestination
	}

	// Gracefully locate all the columns to load.
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("unexpected type %s (expected pointer)", t.Kind())
	}
	t = t.Elem()
	if t.Kind() != reflect.Slice {
		return fmt.Errorf("unexpected type %s (expected slice)", t.Kind())
	}
	elemType := t.Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected type %s (expected struct)", elemType.Kind())
	}

	cols := new(strings.Builder)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		dbTag := field.Tag.Get("db")
		switch dbTag {
		case "":
			dbTag = strings.ToLower(field.Name)
		case "-":
			continue
		}

		if cols.Len() > 0 {
			cols.WriteString(", ")
		}

		// If the db tag contains "autoinc" or "pk" then we need to use the first part of the db tag
		// as the column name and the second part as the alias. This is because the db tag is used to
		// generate the SQL query and the SQL query must be valid.
		//
		// e.g. `db:"id,autoinc,pk"` will generate "t.id 'id'" in the SQL query
		if structTags := strings.Split(dbTag, ","); len(structTags) > 1 {
			for _, tag := range structTags {
				switch tag {
				case dbTagAutoIncrement:
					dbTag = strings.ReplaceAll(dbTag, ","+dbTagAutoIncrement, "")
				case dbTagPrimaryKey:
					dbTag = strings.ReplaceAll(dbTag, ","+dbTagPrimaryKey, "")
				case dbTagDefault:
					dbTag = strings.ReplaceAll(dbTag, ","+dbTagDefault, "")
				}
			}
		}

		// In order for our db tag to remain compatible with the sql db tag mapper
		// we must use commas as separators. However, we want everything after the first comma to be
		// one argument, as the arbitrary SQL there may itself contain commas, hence the SplitN
		args := strings.SplitN(dbTag, ",", 2)
		switch len(args) {
		case 1:
			if len(strings.Split(args[0], ".")) == 2 {
				cols.WriteString(args[0] + " '" + args[0] + "'")
			} else {
				cols.WriteString("t." + args[0])
			}
		case 2:
			cols.WriteString(args[1] + " '" + args[0] + "'")
		}
	}

	jSQL, jArgs := p.filter.Join()
	wSQL, wArgs := p.filter.Where()
	var gSQL string
	if g, ok := p.filter.(Grouper); ok && len(g.Group()) > 0 {
		gSQL = fmt.Sprintf("GROUP BY %s", strings.Join(g.Group(), ", "))
	}

	// Be aware of SQL injection if modifying the below SQL. Any parameters in the sprintf
	// MUST not be allowed to be created by external input.
	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("SELECT ")
	sqlBuilder.WriteString(cols.String())
	sqlBuilder.WriteString(" \n")
	sqlBuilder.WriteString("FROM ")
	sqlBuilder.WriteString(p.table)
	sqlBuilder.WriteString(" t \n")

	if jSQL != "" {
		sqlBuilder.WriteString(jSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("WHERE (t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" ")
	sqlBuilder.WriteString(p.details.sortOperator)
	sqlBuilder.WriteString(" ? OR (t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" = ? AND t.")
	sqlBuilder.WriteString(p.idKey)
	sqlBuilder.WriteString(" > ?)) \n")

	if wSQL != "" {
		sqlBuilder.WriteString("AND (\n")
		sqlBuilder.WriteString(trimWherePrefix(wSQL))
		sqlBuilder.WriteString("\n)\n")
	}

	if gSQL != "" {
		sqlBuilder.WriteString(gSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("ORDER BY t.")
	sqlBuilder.WriteString(p.details.SortBy)
	sqlBuilder.WriteString(" ")
	sqlBuilder.WriteString(p.details.sortComparator)
	sqlBuilder.WriteString(", t.")
	sqlBuilder.WriteString(p.idKey)
	sqlBuilder.WriteString(" ASC \n")

	args := append(jArgs, pivot, pivot, p.details.LastID)
	args = append(args, wArgs...)

	if p.details.Limit > 0 {
		sqlBuilder.WriteString("LIMIT ?")
		args = append(args, p.details.Limit)
	}

	sql := sqlBuilder.String()
	var err error
	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		return fmt.Errorf("retrieve sql in: %w", err)
	}

	err = p.db.Select(dest, sql, args...)
	if err != nil {
		return fmt.Errorf("retrieve select: %w", err)
	}

	return nil
}

// Counts returns the total number of records in the table given the provided filters. This does not take into
// account of the current pivot or limit.
func (p *Paginator) Counts(dest *int64) error {
	jSQL, jArgs := p.filter.Join()
	wSQL, wArgs := p.filter.Where()
	var gSQL string
	if g, ok := p.filter.(Grouper); ok && len(g.Group()) > 0 {
		gSQL = fmt.Sprintf("GROUP BY %s", strings.Join(g.Group(), ", "))
	}

	// Be aware of SQL injection if modifying the below SQL. Any parameters in the sprintf
	// MUST not be allowed to be created by external input.
	sqlBuilder := new(strings.Builder)
	sqlBuilder.WriteString("SELECT COUNT(*) \n")
	sqlBuilder.WriteString("FROM ")
	sqlBuilder.WriteString(p.table)
	sqlBuilder.WriteString(" t \n")

	if jSQL != "" {
		sqlBuilder.WriteString(jSQL)
		sqlBuilder.WriteString(" \n")
	}

	sqlBuilder.WriteString("WHERE (1=1) \n")

	if wSQL != "" {
		sqlBuilder.WriteString("AND (\n")
		sqlBuilder.WriteString(trimWherePrefix(wSQL))
		sqlBuilder.WriteString("\n)\n")
	}

	if gSQL != "" {
		sqlBuilder.WriteString(gSQL)
		sqlBuilder.WriteString(" \n")
	}

	sql := sqlBuilder.String()

	args := append(jArgs, wArgs...)
	var err error
	sql, args, err = sqlx.In(sql, args...)
	if err != nil {
		return fmt.Errorf("counts sql in: %w", err)
	}

	err = p.db.Get(dest, sql, args...)
	if err != nil {
		return fmt.Errorf("counts select: %w", err)
	}

	return nil
}

func trimWherePrefix(w string) string {
	if strings.HasPrefix(w, string(WhereTypeAnd)) || strings.HasPrefix(w, string(WhereTypeOr)) {
		w = strings.TrimPrefix(w, string(WhereTypeAnd))
		w = strings.TrimPrefix(w, string(WhereTypeOr))
	}
	return strings.TrimSpace(w)
}
