package data

import "github.com/Jacobbrewer1/f1-data/pkg/models"

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
