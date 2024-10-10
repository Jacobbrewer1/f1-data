// Code generated by mockery. DO NOT EDIT.

package data

import (
	models "github.com/Jacobbrewer1/f1-data/pkg/models"
	mock "github.com/stretchr/testify/mock"

	pagefilter "github.com/Jacobbrewer1/pagefilter"
)

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// GetConstructorsChampionship provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetConstructorsChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetConstructorsChampionshipFilters) (*PaginationResponse[models.ConstructorChampionship], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetConstructorsChampionship")
	}

	var r0 *PaginationResponse[models.ConstructorChampionship]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetConstructorsChampionshipFilters) (*PaginationResponse[models.ConstructorChampionship], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetConstructorsChampionshipFilters) *PaginationResponse[models.ConstructorChampionship]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.ConstructorChampionship])
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetConstructorsChampionshipFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDrivers provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetDrivers(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversFilters) (*PaginationResponse[models.DriverChampionship], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetDrivers")
	}

	var r0 *PaginationResponse[models.DriverChampionship]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetDriversFilters) (*PaginationResponse[models.DriverChampionship], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetDriversFilters) *PaginationResponse[models.DriverChampionship]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.DriverChampionship])
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetDriversFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDriversChampionship provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetDriversChampionship(paginationDetails *pagefilter.PaginatorDetails, filters *GetDriversChampionshipFilters) (*PaginationResponse[models.DriverChampionship], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetDriversChampionship")
	}

	var r0 *PaginationResponse[models.DriverChampionship]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetDriversChampionshipFilters) (*PaginationResponse[models.DriverChampionship], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetDriversChampionshipFilters) *PaginationResponse[models.DriverChampionship]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.DriverChampionship])
		}
	}

	if rf, ok := ret.Get(1).(func(*pagefilter.PaginatorDetails, *GetDriversChampionshipFilters) error); ok {
		r1 = rf(paginationDetails, filters)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRaceResults provides a mock function with given fields: paginationDetails, filters
func (_m *MockRepository) GetRaceResults(paginationDetails *pagefilter.PaginatorDetails, filters *GetRaceResultsFilters) (*PaginationResponse[models.RaceResult], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetRaceResults")
	}

	var r0 *PaginationResponse[models.RaceResult]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetRaceResultsFilters) (*PaginationResponse[models.RaceResult], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetRaceResultsFilters) *PaginationResponse[models.RaceResult]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.RaceResult])
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
func (_m *MockRepository) GetSeasonRaces(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonRacesFilters) (*PaginationResponse[models.Race], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetSeasonRaces")
	}

	var r0 *PaginationResponse[models.Race]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonRacesFilters) (*PaginationResponse[models.Race], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonRacesFilters) *PaginationResponse[models.Race]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.Race])
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
func (_m *MockRepository) GetSeasons(paginationDetails *pagefilter.PaginatorDetails, filters *GetSeasonsFilters) (*PaginationResponse[models.Season], error) {
	ret := _m.Called(paginationDetails, filters)

	if len(ret) == 0 {
		panic("no return value specified for GetSeasons")
	}

	var r0 *PaginationResponse[models.Season]
	var r1 error
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonsFilters) (*PaginationResponse[models.Season], error)); ok {
		return rf(paginationDetails, filters)
	}
	if rf, ok := ret.Get(0).(func(*pagefilter.PaginatorDetails, *GetSeasonsFilters) *PaginationResponse[models.Season]); ok {
		r0 = rf(paginationDetails, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PaginationResponse[models.Season])
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
