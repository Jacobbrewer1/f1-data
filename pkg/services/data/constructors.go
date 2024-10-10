package data

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/common"
	api "github.com/Jacobbrewer1/f1-data/pkg/codegen/apis/data"
	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/data"
	"github.com/Jacobbrewer1/f1-data/pkg/utils"
	uhttp "github.com/Jacobbrewer1/f1-data/pkg/utils/http"
	"github.com/Jacobbrewer1/pagefilter"
)

func (s *service) GetConstructorsChampionship(w http.ResponseWriter, r *http.Request, year api.PathYear, params api.GetConstructorsChampionshipParams) {
	var sortDir *common.SortDirection
	if params.SortDir != nil {
		sortDir = (*common.SortDirection)(params.SortDir)
	}
	paginationDetails := pagefilter.GetPaginatorDetails(params.Limit, params.LastVal, params.LastId, params.SortBy, sortDir)

	// If the limit is not set, remove it from the pagination details.
	if params.Limit == nil {
		paginationDetails.RemoveLimit()
	}

	filts, err := s.getConstructorsChampionshipFilters(&year, params.Name)
	if err != nil {
		slog.Error("Failed to parse filters", slog.String(logging.KeyError, err.Error()))
		uhttp.SendErrorMessageWithStatus(w, http.StatusBadRequest, "failed to parse filters", err)
		return
	}

	constructorChampionship, err := s.r.GetConstructorsChampionship(paginationDetails, filts)
	if err != nil {
		switch {
		case errors.Is(err, repo.ErrConstructorChampionshipNotFound):
			constructorChampionship = &repo.PaginationResponse[models.ConstructorChampionship]{
				Items: make([]*models.ConstructorChampionship, 0),
				Total: 0,
			}
		default:
			slog.Error("Error getting drivers championship", slog.String(logging.KeyError, err.Error()))
			uhttp.SendErrorMessageWithStatus(w, http.StatusInternalServerError, "error getting drivers championship", err)
			return
		}
	}

	respArray := make([]api.ConstructorChampionship, len(constructorChampionship.Items))
	for i, driverChampionship := range constructorChampionship.Items {
		respArray[i] = *s.modelAsApiConstructor(driverChampionship)
	}

	resp := &api.ConstructorChampionshipResponse{
		Constructors: respArray,
		Total:        constructorChampionship.Total,
	}

	err = uhttp.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("Error encoding drivers championship to user", slog.String(logging.KeyError, err.Error()))
		return
	}
}

func (s *service) getConstructorsChampionshipFilters(
	year *api.PathYear,
	name *string,
) (*repo.GetConstructorsChampionshipFilters, error) {
	filters := new(repo.GetConstructorsChampionshipFilters)

	if year != nil {
		filters.SeasonYear = utils.Ptr(int(*year))
	}

	if name != nil {
		filters.ConstructorName = name
	}

	return filters, nil
}

func (s *service) modelAsApiConstructor(c *models.ConstructorChampionship) *api.ConstructorChampionship {
	return &api.ConstructorChampionship{
		Id:       utils.Ptr(int64(c.Id)),
		Name:     utils.Ptr(c.Name),
		Points:   utils.Ptr(float32(c.Points)),
		Position: utils.Ptr(int64(c.Position)),
	}
}
