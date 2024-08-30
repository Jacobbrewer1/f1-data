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

func (s *service) GetRaceResults(w http.ResponseWriter, r *http.Request, raceId api.PathRaceId, params api.GetRaceResultsParams) {
	var sortDir *common.SortDirection
	if params.SortDir != nil {
		sortDir = (*common.SortDirection)(params.SortDir)
	}
	sortBy := params.SortBy
	if sortBy != nil && *sortBy == "position" {
		// As we store the position as a string to store non-classified results as "NC", we need to sort by position+0
		// to get the correct order.
		sortBy = utils.Ptr("position+0")
	}
	paginationDetails := pagefilter.GetPaginatorDetails(params.Limit, params.LastVal, params.LastId, sortBy, sortDir)

	// If the limit is not set, remove it from the pagination details.
	if params.Limit == nil {
		paginationDetails.RemoveLimit()
	}

	filts, err := s.getRaceResultFilters(&raceId)
	if err != nil {
		slog.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	raceResults, err := s.r.GetRaceResults(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoRacesFound):
			raceResults = make([]*models.RaceResult, 0)
		default:
			slog.Error("Error getting season races", slog.String(logging.KeyError, err.Error()))
		}
	}

	resp := make([]*api.RaceResult, len(raceResults))
	for i, rr := range raceResults {
		resp[i] = s.modelAsApiRaceResult(rr)
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding season to user", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getRaceResultFilters(
	raceId *api.PathRaceId,
) (*repo.GetRaceResultsFilters, error) {
	f := new(repo.GetRaceResultsFilters)

	if raceId != nil {
		f.RaceID = utils.Ptr(int(*raceId))
	}

	return f, nil
}

func (s *service) modelAsApiRaceResult(m *models.RaceResult) *api.RaceResult {
	return &api.RaceResult{
		DriverName:    utils.Ptr(m.Driver),
		DriverNumber:  utils.Ptr(int64(m.DriverNumber)),
		DriverTag:     utils.Ptr(m.DriverTag),
		Id:            utils.Ptr(int64(m.Id)),
		LapsCompleted: utils.Ptr(int64(m.Laps)),
		Points:        utils.Ptr(int64(m.Points)),
		Position:      utils.Ptr(m.Position),
		TeamName:      utils.Ptr(m.Team),
		TimeOrRetired: utils.Ptr(m.TimeRetired),
	}
}
