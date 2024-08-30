package importer

import "github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"

type Service interface {
	Import(from, to int) error
}

// service is the service used by the importer.
type service struct {
	// r is the repository used by the service.
	r importer.Repository

	// baseUrl is the base URL used to get the data.
	baseUrl string
}

// NewService creates a new service.
func NewService(r importer.Repository, baseUrl string) Service {
	return &service{
		r:       r,
		baseUrl: baseUrl,
	}
}
