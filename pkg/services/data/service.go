package data

import (
	api "github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/data"
	"github.com/Jacobbrewer1/f1-data/pkg/repositories/data"
)

type service struct {
	// r is the repository used by the service.
	r data.Repository
}

// NewService creates a new service.
func NewService(r data.Repository) api.ServerInterface {
	return &service{
		r: r,
	}
}
