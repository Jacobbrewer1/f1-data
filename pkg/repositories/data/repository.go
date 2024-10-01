package data

import (
	"github.com/Jacobbrewer1/vaulty/pkg/repositories"
)

type repository struct {
	// db is the database used by the repository.
	db *repositories.Database
}

// NewRepository creates a new repository.
func NewRepository(db *repositories.Database) Repository {
	return &repository{
		db: db,
	}
}
