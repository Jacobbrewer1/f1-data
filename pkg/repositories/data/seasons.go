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
	ErrNoSeasonsFound = errors.New("no seasons found")
)

func (r *repository) GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) ([]*models.Season, error) {
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

	items := make([]season, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoSeasonsFound
		default:
			return nil, fmt.Errorf("retrieve seasons: %w", err)
		}
	}

	returnItems := make([]*models.Season, len(items))
	for i, item := range items {
		returnItems[i] = item.AsModel()
	}

	return returnItems, nil
}

func (r *repository) getSeasonsFilters(filters *GetSeasonsFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.Year != nil {
		f := repoFilter.NewYear(*filters.Year)
		mf.Add(f)
	}

	if filters.YearMin != nil {
		f := repoFilter.NewYearMin(*filters.YearMin)
		mf.Add(f)
	}

	if filters.YearMax != nil {
		f := repoFilter.NewYearMax(*filters.YearMax)
		mf.Add(f)
	}

	return mf
}
