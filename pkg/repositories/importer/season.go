package importer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
)

var (
	// ErrNoSeasonFound is returned when no season is found.
	ErrNoSeasonFound = errors.New("no season found")
)

func (r *repository) SaveSeason(season *models.Season) error {
	err := season.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return fmt.Errorf("failed to save season: %w", err)
		}
	}

	return nil
}

func (r *repository) GetSeasonByYear(year int) (*models.Season, error) {
	sqlStmt := `SELECT id FROM season WHERE year = ?`

	id := 0
	err := r.db.Get(&id, sqlStmt, year)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoSeasonFound
		default:
			return nil, fmt.Errorf("failed to get season by year: %w", err)
		}
	}

	return models.SeasonById(r.db, id)
}
