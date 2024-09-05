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

var (
	// ErrConstructorChampionshipNotFound is returned when a constructor championship is not found
	ErrConstructorChampionshipNotFound = errors.New("constructor championship not found")
)

func (r *repository) GetConstructorsChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetConstructorsChampionshipFilters) (*PaginationResponse[models.ConstructorChampionship], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_constructors_championship"))
	defer t.ObserveDuration()

	mf := r.getConstructorsChampionshipFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "constructor_championship", "id", mf)

	if err := pg.SetDetails(paginationDetails, "id", "position"); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrConstructorChampionshipNotFound
		default:
			return nil, fmt.Errorf("set paginator details: %w", err)
		}
	}

	pvt, err := pg.Pivot()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrConstructorChampionshipNotFound
		default:
			return nil, fmt.Errorf("failed to pivot: %w", err)
		}
	}

	items := make([]constructorChampionship, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrConstructorChampionshipNotFound
		default:
			return nil, fmt.Errorf("failed to retrieve: %w", err)
		}
	}

	returnItems := make([]*models.ConstructorChampionship, len(items))
	for i, item := range items {
		returnItems[i] = item.AsModel()
	}

	var total int64 = 0
	err = pg.Counts(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	resp := &PaginationResponse[models.ConstructorChampionship]{
		Items: returnItems,
		Total: total,
	}

	return resp, nil
}

func (r *repository) getConstructorsChampionshipFilters(filters *GetConstructorsChampionshipFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.SeasonYear != nil {
		f := repoFilter.NewDriverChampSeasonYear(*filters.SeasonYear)
		mf.Add(f)
	}

	if filters.ConstructorName != nil {
		f := repoFilter.NewDriverChampNameLike(*filters.ConstructorName)
		mf.Add(f)
	}

	return mf
}
