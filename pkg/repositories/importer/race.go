package importer

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jacobbrewer1/f1-data/pkg/models"
)

var (
	// ErrRaceNotFound is returned when the race is not found.
	ErrRaceNotFound = errors.New("race not found")
)

func (r *repository) SaveRace(race *models.Race) error {
	race.UpdatedAt = time.Now().UTC()
	err := race.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return err
		}
	}

	return nil
}

func (r *repository) GetRaceBySeasonIdAndGrandPrix(seasonId int, grandPrix string) (*models.Race, error) {
	sqlStmt := `SELECT id FROM race WHERE season_id = ? AND grand_prix LIKE ?`

	id := 0
	err := r.db.Get(&id, sqlStmt, seasonId, grandPrix)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRaceNotFound
		default:
			return nil, err
		}
	}

	return models.RaceById(r.db, id)
}
