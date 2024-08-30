package data

import (
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
)

type Repository interface {
	// GetSeasons returns all seasons.
	GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) ([]*models.Season, error)
}
