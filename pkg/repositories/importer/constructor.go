package importer

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jacobbrewer1/f1-data/pkg/models"
)

var (
	// ErrConstructorNotFound is returned when a constructor is not found.
	ErrConstructorNotFound = errors.New("constructor not found")
)

func (r *repository) GetConstructorByName(seasonId int, name string) (*models.ConstructorChampionship, error) {
	sqlStmt := `SELECT id FROM constructor_championship WHERE season_id = ? AND name = ?`

	id := 0
	err := r.db.Get(&id, sqlStmt, seasonId, name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrConstructorNotFound
		default:
			return nil, fmt.Errorf("error getting constructor by name: %w", err)
		}
	}

	return models.ConstructorChampionshipById(r.db, id)
}

func (r *repository) SaveConstructor(constructor *models.ConstructorChampionship) error {
	err := constructor.SaveOrUpdate(r.db)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoAffectedRows):
			break
		default:
			return fmt.Errorf("error saving constructor: %w", err)
		}
	}

	return nil
}
