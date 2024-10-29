package importer

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jacobbrewer1/f1-data/pkg/models"
)

var (
	// ErrNoRaceResultFound is returned when no race result is found.
	ErrNoRaceResultFound = errors.New("race result not found")
)

func (r *repository) SaveRaceResult(raceResult *models.RaceResult) error {
	raceResult.UpdatedAt = time.Now().UTC()
	err := raceResult.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return fmt.Errorf("failed to save race result: %w", err)
		}
	}

	return nil
}

func (r *repository) GetRaceResultByRaceIdAndDriverNumber(raceId int, driverNumber int) (*models.RaceResult, error) {
	sqlStmt := `SELECT id from race_result WHERE race_id = ? AND driver_number = ?`

	id := 0
	err := r.db.Get(&id, sqlStmt, raceId, driverNumber)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRaceResultFound
		default:
			return nil, fmt.Errorf("failed to get id for race result: %w", err)
		}
	}

	return models.RaceResultById(r.db, id)
}
