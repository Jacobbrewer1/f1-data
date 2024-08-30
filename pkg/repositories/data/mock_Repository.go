// Code generated by mockery v2.45.0. DO NOT EDIT.

package data

import (
	models "github.com/Jacobbrewer1/f1-data/pkg/models"
	mock "github.com/stretchr/testify/mock"

	pagefilter "github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// GetRaceResults provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) ([]*models.RaceResult, error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetRaceResults")
	}

	var r0 []*models.RaceResult
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetRaceResultsFilters) ([]*models.RaceResult, error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetRaceResultsFilters) []*models.RaceResult); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.RaceResult)
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetRaceResultsFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSeasonRaces provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) ([]*models.Race, error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetSeasonRaces")
	}

	var r0 []*models.Race
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonRacesFilters) ([]*models.Race, error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonRacesFilters) []*models.Race); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Race)
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetSeasonRacesFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSeasons provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) ([]*models.Season, error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetSeasons")
	}

	var r0 []*models.Season
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonsFilters) ([]*models.Season, error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonsFilters) []*models.Season); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Season)
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetSeasonsFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockRepository creates a new instance of MockRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockRepository {
	mock := &MockRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
