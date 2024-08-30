package data

import (
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
)

type Repository interface {
	// GetSeasons returns all seasons
	GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) ([]*models.Season, error)

	// GetSeasonRaces returns all races for a season
	GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) ([]*models.Race, error)

	// GetRaceResults returns all results for a specific race
	GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) ([]*models.RaceResult, error)
}
