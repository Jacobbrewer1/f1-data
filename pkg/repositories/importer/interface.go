package importer

import "github.com/Jacobbrewer1/f1-data/pkg/models"

type Repository interface {
	// GetSeasonByYear returns a season by year
	GetSeasonByYear(year int) (*models.Season, error)

	// GetRaceBySeasonIdAndGrandPrix returns a race by season id and grand prix
	GetRaceBySeasonIdAndGrandPrix(seasonId int, grandPrix string) (*models.Race, error)

	// GetRaceResultByRaceIdAndDriverNumber returns a race result
	GetRaceResultByRaceIdAndDriverNumber(raceId int, driverNumber int) (*models.RaceResult, error)

	// SaveSeason saves a season
	SaveSeason(season *models.Season) error

	// SaveRace saves a race
	SaveRace(race *models.Race) error

	// SaveRaceResult saves a race result
	SaveRaceResult(raceResult *models.RaceResult) error
}
