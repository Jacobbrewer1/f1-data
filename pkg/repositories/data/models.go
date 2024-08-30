package data

import (
	"time"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
)

// TODO: This is extremely gross and should be changed.

type season struct {
	Id   int `db:"id"`
	Year int `db:"year"`
}

func (s *season) AsModel() *models.Season {
	return &models.Season{
		Id:   s.Id,
		Year: s.Year,
	}
}

type race struct {
	Id        int       `db:"id"`
	SeasonId  int       `db:"season_id"`
	GrandPrix string    `db:"grand_prix"`
	Date      time.Time `db:"date"`
}

func (r *race) AsModel() *models.Race {
	return &models.Race{
		Id:        r.Id,
		SeasonId:  r.SeasonId,
		GrandPrix: r.GrandPrix,
		Date:      r.Date,
	}
}

type raceResult struct {
	Id           int     `db:"id"`
	RaceId       int     `db:"race_id"`
	Position     string  `db:"position"`
	DriverNumber int     `db:"driver_number"`
	Driver       string  `db:"driver"`
	Team         string  `db:"team"`
	Laps         int     `db:"laps"`
	TimeRetired  string  `db:"time_retired"`
	Points       float64 `db:"points"`
}

func (r *raceResult) AsModel() *models.RaceResult {
	return &models.RaceResult{
		Id:           r.Id,
		RaceId:       r.RaceId,
		Position:     r.Position,
		DriverNumber: r.DriverNumber,
		Driver:       r.Driver,
		Team:         r.Team,
		Laps:         r.Laps,
		TimeRetired:  r.TimeRetired,
		Points:       r.Points,
	}
}
