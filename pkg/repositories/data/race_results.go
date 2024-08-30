package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
	repoFilter "github.com/Jacobbrewer1/f1-data/pkg/repositories/data/filters"
)

func (r *repository) GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) ([]*models.RaceResult, error) {
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

	return returnItems, nil
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
