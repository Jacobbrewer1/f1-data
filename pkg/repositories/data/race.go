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
	ErrNoRacesFound = errors.New("no races found")
)

func (r *repository) GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) (*PaginationResponse[models.Race], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_season_races"))
	defer t.ObserveDuration()

	mf := r.getSeasonRacesFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "race", "date", mf)

	if err := pg.SetDetails(paginationDetails, "date", "date"); err != nil {
		return nil, fmt.Errorf("set paginator details: %w", err)
	}

	pvt, err := pg.Pivot()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRacesFound
		default:
			return nil, fmt.Errorf("failed to pivot: %w", err)
		}
	}

	items := make([]*models.Race, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRacesFound
		default:
			return nil, fmt.Errorf("failed to retrieve: %w", err)
		}
	}

	var total int64 = 0
	err = pg.Counts(&total)
	if err != nil {
		return nil, fmt.Errorf("get total count: %w", err)
	}

	resp := &PaginationResponse[models.Race]{
		Items: items,
		Total: total,
	}

	return resp, nil
}

func (r *repository) getSeasonRacesFilters(filters *GetSeasonRacesFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.SeasonYear != nil {
		f := repoFilter.NewRaceYear(*filters.SeasonYear)
		mf.Add(f)
	}

	return mf
}
