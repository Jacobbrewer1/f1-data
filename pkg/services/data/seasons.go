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

func (s *service) GetSeasons(w http.ResponseWriter, r *http.Request, params api.GetSeasonsParams) {
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

	filts, err := s.getSeasonsFilters(params.Year, params.YearMin, params.YearMax)
	if err != nil {
		l.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	seasons, err := s.r.GetSeasons(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrNoSeasonsFound):
			seasons = &repo.PaginationResponse[models.Season]{
				Items: make([]*models.Season, 0),
				Total: 0,
			}
		default:
			slog.Error("Error getting seasons", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting seasons", err)
			return
		}
	}

	respArray := make([]api.Season, len(seasons.Items))
	for i, season := range seasons.Items {
		respArray[i] = *s.modelAsApiSeason(season)
	}

	resp := &api.SeasonResponse{
		Seasons: respArray,
		Total:   seasons.Total,
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		l.Error("Error encoding seasons", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getSeasonsFilters(
	year *int64,
	yearMin *int64,
	yearMax *int64,
) (*repo.GetSeasonsFilters, error) {
	f := new(repo.GetSeasonsFilters)

	if year != nil {
		f.Year = utils.Ptr(int(*year))
	}

	if yearMin != nil {
		if year != nil {
			return nil, errors.New("cannot specify both year and year_min")
		}
		f.YearMin = utils.Ptr(int(*yearMin))
	}

	if yearMax != nil {
		if year != nil {
			return nil, errors.New("cannot specify both year and year_max")
		}
		f.YearMax = utils.Ptr(int(*yearMax))
	}

	if yearMin != nil && yearMax != nil {
		if *yearMin > *yearMax {
			return nil, errors.New("year_min cannot be greater than year_max")
		}
	}

	return f, nil
}

func (s *service) modelAsApiSeason(season *models.Season) *api.Season {
	return &api.Season{
		Id:   utils.Ptr(int64(season.Id)),
		Year: utils.Ptr(int64(season.Year)),
	}
}
