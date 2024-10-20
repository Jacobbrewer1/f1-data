package pagefilter

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/jacobbrewer1/pagefilter/common"
)

var (
	ErrInvalidPaginatorDetails = errors.New("invalid paginator details")
)

// PaginatorDetails contains pagination details
type PaginatorDetails struct {
	Limit          int
	LastVal        string
	LastID         string
	SortBy         string
	SortDir        string
	sortComparator string
	sortOperator   string
}

func getLimit(q url.Values) (limit int, err error) {
	limit = defaultPageLimit
	if limitStr := q.Get(QueryLimit); limitStr != "" {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			return -1, fmt.Errorf("invalid limit: %w", err)
		}
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if limit == 0 {
		limit = defaultPageLimit
	}
	return limit, nil
}

// GetPaginatorDetails loads paginator details from a request. Requests have each pagination detail determined
// separately by codegen.
func GetPaginatorDetails(
	limit *common.LimitParam,
	lastVal *common.LastValue,
	lastID *common.LastId,
	sortBy *common.SortBy,
	sortDir *common.SortDirection,
) *PaginatorDetails {
	d := new(PaginatorDetails)

	if lastID != nil {
		d.LastID = *lastID
	}
	if lastVal != nil {
		d.LastVal = *lastVal
	}
	if limit != nil {
		d.Limit, _ = strconv.Atoi(*limit)
	}
	if sortBy != nil {
		d.SortBy = *sortBy
	}
	if sortDir != nil {
		d.SortDir = *sortDir
	}

	if d.Limit <= 0 {
		d.Limit = defaultPageLimit
	}

	if d.Limit > maxLimit {
		d.Limit = maxLimit
	}

	return d
}

// RemoveLimit removes the limit from the paginator details.
func (p *PaginatorDetails) RemoveLimit() {
	p.Limit = -1
}
