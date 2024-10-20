package data

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jacobbrewer1/f1-data/pkg/codegen/apis/common"
	api "github.com/jacobbrewer1/f1-data/pkg/codegen/apis/data"
	"github.com/jacobbrewer1/f1-data/pkg/logging"
	"github.com/jacobbrewer1/f1-data/pkg/models"
	repo "github.com/jacobbrewer1/f1-data/pkg/repositories/data"
	"github.com/jacobbrewer1/f1-data/pkg/utils"
	"github.com/jacobbrewer1/pagefilter"
	"github.com/jacobbrewer1/uhttp"
)

func (s *service) GetDriversChampionship(w http.ResponseWriter, r *http.Request, year api.PathYear, params api.GetDriversChampionshipParams) {
	l := logging.LoggerFromRequest(r)

	sortDir := new(common.SortDirection)
	if params.SortDir != nil {
		sortDir = (*common.SortDirection)(params.SortDir)
	}
	paginationDetails := pagefilter.GetPaginatorDetails(params.Limit, params.LastVal, params.LastId, params.SortBy, sortDir)

	// If the limit is not set, remove it from the pagination details.
	if params.Limit == nil {
		paginationDetails.RemoveLimit()
	}

	filts, err := s.getDriversChampionshipFilters(&year, params.Name, params.Tag, params.Team)
	if err != nil {
		l.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	driversChampionship, err := s.r.GetDriversChampionship(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDriverChampionshipNotFound):
			driversChampionship = &repo.PaginationResponse[models.DriverChampionship]{
				Items: make([]*models.DriverChampionship, 0),
				Total: 0,
			}
		default:
			slog.Error("Error getting drivers championship", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting drivers championship", err)
			return
		}
	}

	respArray := make([]api.DriverChampionship, len(driversChampionship.Items))
	for i, driverChampionship := range driversChampionship.Items {
		respArray[i] = *s.modelAsApiDriverChampionship(driverChampionship)
	}

	resp := &api.DriverChampionshipResponse{
		Drivers: respArray,
		Total:   driversChampionship.Total,
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		l.Error("Error encoding season to user", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getDriversChampionshipFilters(
	year *api.PathYear,
	name *api.QueryName,
	tag *api.QueryTag,
	team *api.QueryTeam,
) (*repo.GetDriversChampionshipFilters, error) {
	filters := new(repo.GetDriversChampionshipFilters)

	if year != nil {
		filters.SeasonYear = utils.Ptr(int(*year))
	}

	if name != nil {
		filters.DriverName = name
	}

	if tag != nil {
		filters.DriverTag = tag
	}

	if team != nil {
		filters.Team = team
	}

	return filters, nil
}

func (s *service) modelAsApiDriverChampionship(m *models.DriverChampionship) *api.DriverChampionship {
	return &api.DriverChampionship{
		Id:          utils.Ptr(int64(m.Id)),
		Name:        utils.Ptr(m.Driver),
		Nationality: utils.Ptr(m.Nationality),
		Points:      utils.Ptr(float32(m.Points)),
		Position:    utils.Ptr(int64(m.Position)),
		Tag:         utils.Ptr(m.DriverTag),
		Team:        utils.Ptr(m.Team),
	}
}

func (s *service) GetDrivers(w http.ResponseWriter, r *http.Request, params api.GetDriversParams) {
	l := logging.LoggerFromRequest(r)

	sortDir := new(common.SortDirection)
	if params.SortDir != nil {
		sortDir = (*common.SortDirection)(params.SortDir)
	}
	paginationDetails := pagefilter.GetPaginatorDetails(params.Limit, params.LastVal, params.LastId, params.SortBy, sortDir)

	// If the limit is not set, remove it from the pagination details.
	if params.Limit == nil {
		paginationDetails.RemoveLimit()
	}

	filts, err := s.getDriversFilters(params.Name, params.Tag, params.Team, params.Nationality)
	if err != nil {
		l.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	drivers, err := s.r.GetDrivers(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDriversNotFound):
			drivers = &repo.PaginationResponse[models.DriverChampionship]{
				Items: make([]*models.DriverChampionship, 0),
				Total: 0,
			}
		default:
			l.Error("Error getting drivers", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting drivers", err)
			return
		}
	}

	respArray := make([]api.Driver, len(drivers.Items))
	for i, driver := range drivers.Items {
		respArray[i] = *s.modelAsApiDriver(driver)
	}

	resp := &api.DriverResponse{
		Drivers: respArray,
		Total:   drivers.Total,
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		l.Error("Error encoding drivers to user", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getDriversFilters(
	name *api.QueryName,
	tag *api.QueryTag,
	team *api.QueryTeam,
	nationality *api.QueryNationality,
) (*repo.GetDriversFilters, error) {
	filters := new(repo.GetDriversFilters)

	if name != nil {
		filters.Name = name
	}

	if tag != nil {
		filters.Tag = tag
	}

	if team != nil {
		filters.Team = team
	}

	if nationality != nil {
		filters.Nationality = nationality
	}

	return filters, nil
}

func (s *service) modelAsApiDriver(m *models.DriverChampionship) *api.Driver {
	return &api.Driver{
		Name:        utils.Ptr(m.Driver),
		Nationality: utils.Ptr(m.Nationality),
		Tag:         utils.Ptr(m.DriverTag),
	}
}
