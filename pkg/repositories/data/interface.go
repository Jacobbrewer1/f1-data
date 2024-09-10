package data

import (
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
)

type Repository interface {
	// GetSeasons returns all seasons
	GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) (*PaginationResponse[models.Season], error)

	// GetSeasonRaces returns all races for a season
	GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) (*PaginationResponse[models.Race], error)

	// GetRaceResults returns all results for a specific race
	GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) (*PaginationResponse[models.RaceResult], error)

	// GetDriversChampionship returns the drivers championship for a specific season
	GetDriversChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversChampionshipFilters) (*PaginationResponse[models.DriverChampionship], error)

	// GetConstructorsChampionship returns the constructors' championship for a specific season
	GetConstructorsChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetConstructorsChampionshipFilters) (*PaginationResponse[models.ConstructorChampionship], error)

	// GetDrivers returns all drivers as per the filters
	GetDrivers(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversFilters) (*PaginationResponse[models.DriverChampionship], error)
}
