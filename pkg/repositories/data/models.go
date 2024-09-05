package data

import (
	"time"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
)

type PaginationResponse[T comparable] struct {
	Items []*T  `json:"items"`
	Total int64 `json:"total"`
}

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
	DriverTag    string  `db:"driver_tag"`
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

type driverChampionship struct {
	Id          int     `db:"id"`
	SeasonId    int     `db:"season_id"`
	Position    int     `db:"position"`
	Driver      string  `db:"driver"`
	DriverTag   string  `db:"driver_tag"`
	Nationality string  `db:"nationality"`
	Team        string  `db:"team"`
	Points      float64 `db:"points"`
}

func (d *driverChampionship) AsModel() *models.DriverChampionship {
	return &models.DriverChampionship{
		Id:          d.Id,
		SeasonId:    d.SeasonId,
		Position:    d.Position,
		Driver:      d.Driver,
		DriverTag:   d.DriverTag,
		Nationality: d.Nationality,
		Team:        d.Team,
		Points:      d.Points,
	}
}

type constructorChampionship struct {
	Id       int     `db:"id"`
	SeasonId int     `db:"season_id"`
	Position int     `db:"position"`
	Name     string  `db:"name"`
	Points   float64 `db:"points"`
}

func (c *constructorChampionship) AsModel() *models.ConstructorChampionship {
	return &models.ConstructorChampionship{
		Id:       c.Id,
		SeasonId: c.SeasonId,
		Position: c.Position,
		Name:     c.Name,
		Points:   c.Points,
	}
}
