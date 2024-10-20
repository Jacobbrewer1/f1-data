package importer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jacobbrewer1/f1-data/pkg/models"
)

var (
	// ErrDriverNotFound is returned when a driver is not found.
	ErrDriverNotFound = errors.New("driver not found")
)

func (r *repository) GetDriverByName(seasonId int, name string) (*models.DriverChampionship, error) {
	sqlStmt := `SELECT id FROM driver_championship WHERE season_id = ? AND driver LIKE ?`

	id := 0
	err := r.db.Get(&id, sqlStmt, seasonId, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrDriverNotFound
		default:
			return nil, fmt.Errorf("error getting driver by name: %w", err)
		}
	}

	return models.DriverChampionshipById(r.db, id)
}

func (r *repository) SaveDriver(driver *models.DriverChampionship) error {
	err := driver.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return fmt.Errorf("error saving driver: %w", err)
		}
	}

	return nil
}
