package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jacobbrewer1/f1-data/pkg/models"
	repoFilter "github.com/jacobbrewer1/f1-data/pkg/repositories/data/filters"
	"github.com/jacobbrewer1/pagefilter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ErrNoSeasonsFound = errors.New("no seasons found")
)

func (r *repository) GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) (*PaginationResponse[models.Season], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_seasons"))
	defer t.ObserveDuration()

	mf := r.getSeasonsFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "season", "year", mf)

	if err := pg.SetDetails(paginationDetails, "year", "year"); err != nil {
		return nil, fmt.Errorf("set paginator details: %w", err)
	}

	pvt, err := pg.Pivot()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoSeasonsFound
		default:
			return nil, fmt.Errorf("pivot paginator: %w", err)
		}
	}

	items := make([]*models.Season, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoSeasonsFound
		default:
			return nil, fmt.Errorf("retrieve seasons: %w", err)
		}
	}

	var total int64 = 0
	err = pg.Counts(&total)
	if err != nil {
		return nil, fmt.Errorf("get total count: %w", err)
	}

	resp := &PaginationResponse[models.Season]{
		Items: items,
		Total: total,
	}

	return resp, nil
}

func (r *repository) getSeasonsFilters(filters *GetSeasonsFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.Year != nil {
		f := repoFilter.NewSeasonYear(*filters.Year)
		mf.Add(f)
	}

	if filters.YearMin != nil {
		f := repoFilter.NewSeasonYearMin(*filters.YearMin)
		mf.Add(f)
	}

	if filters.YearMax != nil {
		f := repoFilter.NewSeasonYearMax(*filters.YearMax)
		mf.Add(f)
	}

	return mf
}
