package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
	repoFilter "github.com/Jacobbrewer1/f1-data/pkg/repositories/data/filters"
)

var (
	ErrNoRacesFound = errors.New("no races found")
)

func (r *repository) GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) ([]*models.Race, error) {
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

	items := make([]race, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRacesFound
		default:
			return nil, fmt.Errorf("failed to retrieve: %w", err)
		}
	}

	returnItems := make([]*models.Race, len(items))
	for i, item := range items {
		returnItems[i] = item.AsModel()
	}

	return returnItems, nil
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
