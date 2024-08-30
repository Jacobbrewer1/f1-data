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
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func (s *service) GetSeasonsYearRaces(w http.ResponseWriter, r *http.Request, year api.PathYear, params api.GetSeasonsYearRacesParams) {
	var sortDir *common.SortDirection
	if params.SortDir != nil {
		sortDir = (*common.SortDirection)(params.SortDir)
	}
	paginationDetails := pagefilter.GetPaginatorDetails(params.Limit, params.LastVal, params.LastId, params.SortBy, sortDir)

	// If the limit is not set, remove it from the pagination details.
	if params.Limit == nil {
		paginationDetails.RemoveLimit()
	}

	filts, err := s.getSeasonYearRacesFilters(&year)
	if err != nil {
		slog.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	seasonRaces, err := s.r.GetSeasonRaces(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoRacesFound):
			seasonRaces = make([]*models.Race, 0)
		default:
			slog.Error("Error getting season races", slog.String(logging.KeyError, err.Error()))
		}
	}

	resp := make([]*api.Race, len(seasonRaces))
	for i, seasonRace := range seasonRaces {
		resp[i] = s.modelAsApiRace(seasonRace)
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding season to user", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getSeasonYearRacesFilters(
	year *api.PathYear,
) (*repo.GetSeasonRacesFilters, error) {
	f := new(repo.GetSeasonRacesFilters)

	if year != nil {
		f.SeasonYear = utils.Ptr(int(*year))
	}

	return f, nil
}

func (s *service) modelAsApiRace(m *models.Race) *api.Race {
	date := openapi_types.Date{
		Time: m.Date,
	}
	return &api.Race{
		Date: utils.Ptr(date),
		Id:   utils.Ptr(int64(m.Id)),
		Name: utils.Ptr(m.GrandPrix),
	}
}
