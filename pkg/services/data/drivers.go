package data

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common"
	api "github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/data"
	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	"github.com/Jacobbrewer1/f1-data/pkg/pagefilter"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/data"
	"github.com/Jacobbrewer1/f1-data/pkg/utils"
	uhttp "github.com/Jacobbrewer1/f1-data/pkg/utils/http"
)

func (s *service) GetDriversChampionship(w http.ResponseWriter, r *http.Request, year api.PathYear, params api.GetDriversChampionshipParams) {
	var sortDir *common.SortDirection
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
		slog.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	driversChampionship, err := s.r.GetDriversChampionship(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrDriverChampionshipNotFound):
			driversChampionship = make([]*models.DriverChampionship, 0)
		default:
			slog.Error("Error getting drivers championship", slog.String(logging.KeyError, err.Error()))
		}
	}

	resp := make([]*api.Driver, len(driversChampionship))
	for i, driverChampionship := range driversChampionship {
		resp[i] = s.modelAsApiDriver(driverChampionship)
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding season to user", slog.String(logging.KeyError, err.Error()))
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

func (s *service) modelAsApiDriver(m *models.DriverChampionship) *api.Driver {
	return &api.Driver{
		Id:          utils.Ptr(int64(m.Id)),
		Name:        utils.Ptr(m.Driver),
		Nationality: utils.Ptr(m.Nationality),
		Points:      utils.Ptr(int64(m.Points)),
		Position:    utils.Ptr(int64(m.Position)),
		Tag:         utils.Ptr(m.DriverTag),
		Team:        utils.Ptr(m.Team),
	}
}
