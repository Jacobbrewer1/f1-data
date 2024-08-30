package api

import (
	api "github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/data"
)

type service struct{}

// NewService creates a new service.
func NewService() api.ServerInterface {
	return &service{}
}
