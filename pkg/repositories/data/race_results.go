package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
	repoFilter "github.com/Jacobbrewer1/f1-data/pkg/repositories/data/filters"
	"github.com/prometheus/client_golang/prometheus"
)

func (r *repository) GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) (*PaginationResponse[models.RaceResult], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_race_results"))
	defer t.ObserveDuration()

	mf := r.getRaceResultFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "race_result", "id", mf)

	if err := pg.SetDetails(paginationDetails, "id", "position+0"); err != nil {
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

	items := make([]raceResult, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoSeasonsFound
		default:
			return nil, fmt.Errorf("retrieve seasons: %w", err)
		}
	}

	returnItems := make([]*models.RaceResult, len(items))
	for i, item := range items {
		returnItems[i] = item.AsModel()
	}

	var total int64 = 0
	err = pg.Counts(&total)
	if err != nil {
		return nil, fmt.Errorf("get total count: %w", err)
	}

	resp := &PaginationResponse[models.RaceResult]{
		Items: returnItems,
		Total: total,
	}

	return resp, nil
}

func (r *repository) getRaceResultFilters(filters *GetRaceResultsFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.SeasonYear != nil {
		f := repoFilter.NewRaceResultSeasonYear(*filters.SeasonYear)
		mf.Add(f)
	}

	if filters.RaceID != nil {
		f := repoFilter.NewRaceResultRaceID(*filters.RaceID)
		mf.Add(f)
	}

	return mf
}
