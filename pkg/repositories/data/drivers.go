package data

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repoFilter "github.com/Jacobbrewer1/f1-data/pkg/repositories/data/filters"
	"github.com/Jacobbrewer1/pagefilter"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ErrDriverChampionshipNotFound is returned when a driver championship is not found.
	ErrDriverChampionshipNotFound = errors.New("driver championship not found")

	// ErrDriversNotFound is returned when drivers are not found.
	ErrDriversNotFound = errors.New("drivers not found")
)

func (r *repository) GetDriversChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversChampionshipFilters) (*PaginationResponse[models.DriverChampionship], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_drivers_championship"))
	defer t.ObserveDuration()

	mf := r.getDriversChampionshipFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "driver_championship", "id", mf)

	if err := pg.SetDetails(paginationDetails, "id", "position"); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriverChampionshipNotFound
		default:
			return nil, fmt.Errorf("set paginator details: %w", err)
		}
	}

	pvt, err := pg.Pivot()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriverChampionshipNotFound
		default:
			return nil, fmt.Errorf("failed to pivot: %w", err)
		}
	}

	items := make([]*models.DriverChampionship, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriverChampionshipNotFound
		default:
			return nil, fmt.Errorf("failed to retrieve: %w", err)
		}
	}

	var total int64 = 0
	err = pg.Counts(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count: %w", err)
	}

	resp := &PaginationResponse[models.DriverChampionship]{
		Items: items,
		Total: total,
	}

	return resp, nil
}

func (r *repository) getDriversChampionshipFilters(filters *GetDriversChampionshipFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	if filters == nil {
		return mf
	}

	if filters.SeasonYear != nil {
		f := repoFilter.NewDriverChampSeasonYear(*filters.SeasonYear)
		mf.Add(f)
	}

	if filters.DriverName != nil {
		f := repoFilter.NewDriverChampNameLike(*filters.DriverName)
		mf.Add(f)
	}

	if filters.DriverTag != nil {
		f := repoFilter.NewDriverChampTagLike(*filters.DriverTag)
		mf.Add(f)
	}

	if filters.Team != nil {
		f := repoFilter.NewDriverChampTeamLike(*filters.Team)
		mf.Add(f)
	}

	return mf
}

func (r *repository) GetDrivers(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversFilters) (*PaginationResponse[models.DriverChampionship], error) {
	t := prometheus.NewTimer(models.DatabaseLatency.WithLabelValues("get_drivers"))
	defer t.ObserveDuration()

	mf := r.getDriversFilters(filters)
	pg := pagefilter.NewPaginator(r.db, "driver_championship", "id", mf)

	if err := pg.SetDetails(paginationDetails, "id", "name"); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriversNotFound
		default:
			return nil, fmt.Errorf("set paginator details: %w", err)
		}
	}

	pvt, err := pg.Pivot()
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriversNotFound
		default:
			return nil, fmt.Errorf("failed to pivot: %w", err)
		}
	}

	items := make([]*models.DriverChampionship, 0)
	err = pg.Retrieve(pvt, &items)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriversNotFound
		default:
			return nil, fmt.Errorf("failed to retrieve: %w", err)
		}
	}

	resp := &PaginationResponse[models.DriverChampionship]{
		Items: items,
		Total: int64(len(items)),
	}

	return resp, nil
}

func (r *repository) getDriversFilters(filters *GetDriversFilters) *pagefilter.MultiFilter {
	mf := pagefilter.NewMultiFilter()
	mf.Add(repoFilter.NewDriverChampNameGrouper())
	if filters == nil {
		return mf
	}

	if filters.Name != nil {
		f := repoFilter.NewDriverChampNameLike(*filters.Name)
		mf.Add(f)
	}

	if filters.Tag != nil {
		f := repoFilter.NewDriverChampTagLike(*filters.Tag)
		mf.Add(f)
	}

	if filters.Team != nil {
		f := repoFilter.NewDriverChampTeamLike(*filters.Team)
		mf.Add(f)
	}

	if filters.Nationality != nil {
		f := repoFilter.NewDriverChampNationalityLike(*filters.Nationality)
		mf.Add(f)
	}

	return mf
}
