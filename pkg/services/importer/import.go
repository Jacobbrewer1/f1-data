package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"
)

func (s *service) Import(from, to int) error {
	wg := new(sync.WaitGroup)
	for i := from; i <= to; i++ {
		slog.Debug(fmt.Sprintf("Importing season %d", i))

		season, err := s.r.GetSeasonByYear(i)
		if err != nil && !errors.Is(err, repo.ErrNoSeasonFound) {
			return fmt.Errorf("error getting season by year: %w", err)
		} else if errors.Is(err, repo.ErrNoSeasonFound) {
			season = &models.Season{
				Year: i,
			}
			err = s.r.SaveSeason(season)
			if err != nil {
				return fmt.Errorf("error saving season: %w", err)
			}
		}

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := s.ImportSeasonRaces(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported races for season %d", i))
		}(i)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := s.ImportSeasonDriversChamps(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d drivers championship: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported drivers championship for season %d", i))
		}(i)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err := s.ImportSeasonConstructorsChamps(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d constructors championship: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported constructors championship for season %d", i))
		}(i)
	}

	wg.Wait()

	return nil
}
