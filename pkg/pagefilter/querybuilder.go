package pagefilter

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/jmoiron/sqlx"
)

const (
	defaultPageLimit = 100
	maxLimit         = 20000

	orderAsc  = "asc"
	orderDesc = "desc"

	sqlComparatorAsc  = "ASC"
	sqlComparatorDesc = "DESC"
	sqlOperatorAsc    = ">"
	sqlOperatorDesc   = "<"
)

// Paginator is the struct that provides the paging.
type Paginator struct {
	db      models.DB
	idKey   string
	table   string
	filter  Filter
	details *PaginatorDetails
}

// NewPaginator creates a new paginator
func NewPaginator(db models.DB, table, idk string, f Filter) *Paginator {
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
		LastVal: q.Get("last_val"),
		LastID:  q.Get("last_id"),
		SortBy:  q.Get("sort_by"),
		SortDir: q.Get("sort_dir"),
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
	sql := fmt.Sprintf(`
		SELECT t.%s
		FROM %s t
		%s
		WHERE 1
		%s
		%s
		ORDER BY t.%s %s, t.%s ASC
		LIMIT 1`, p.details.SortBy, p.table, jSQL, wSQL, gSQL, p.details.SortBy, p.details.sortComparator, p.idKey)
	args := append(jArgs, wArgs...)

	var err error
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
	sql := fmt.Sprintf(`
		SELECT t.%s
		FROM %s t
		%s
		WHERE (t.%s = ? AND t.%s >= ?)
		%s
		%s
		LIMIT 1
	`, p.details.SortBy, p.table, jSQL, p.details.SortBy, p.idKey, wSQL, gSQL)
	args := append(jArgs, p.details.LastVal, p.details.LastID)
	args = append(args, wArgs...)

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
func (p *Paginator) Retrieve(pivot string, dest interface{}) error {
	// Gracefully locate all the columns to load.
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("unexpected type %s (expected pointer)", t.Kind())
	}
	if t = t.Elem(); t.Kind() != reflect.Slice {
		return fmt.Errorf("unexpected type %s (expected slice)", t.Kind())
	}
	if t = t.Elem(); t.Kind() != reflect.Struct {
		return fmt.Errorf("unexpected type %s (expected struct)", t.Kind())
	}

	var cols strings.Builder
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
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
	sql := fmt.Sprintf(`
		SELECT %s
		FROM %s t
		%s
		WHERE (t.%s %s ? OR (t.%s = ? AND t.%s > ?))
		%s
		%s
		ORDER BY t.%s %s, t.%s ASC
		LIMIT ?
	`, cols.String(), p.table, jSQL, p.details.SortBy, p.details.sortOperator, p.details.SortBy, p.idKey, wSQL, gSQL, p.details.SortBy, p.details.sortComparator, p.idKey)
	args := append(jArgs, pivot, pivot, p.details.LastID)
	args = append(args, wArgs...)
	args = append(args, p.details.Limit)

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
	sql := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM %s t
		%s
		WHERE (1=1)
		%s
		%s
	`, p.table, jSQL, wSQL, gSQL)
	args := jArgs
	args = append(args, wArgs...)

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
