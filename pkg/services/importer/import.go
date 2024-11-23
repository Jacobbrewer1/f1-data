package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/jacobbrewer1/f1-data/pkg/models"
	repo "github.com/jacobbrewer1/f1-data/pkg/repositories/importer"
	"github.com/jacobbrewer1/workerpool"
)

func (s *service) Import(from, to int) error {
	workers := workerpool.New(
		workerpool.WithDelayedStart(),
	)

	wg := new(sync.WaitGroup)
	for i := from; i <= to; i++ {
		slog.Debug(fmt.Sprintf("Importing season %d", i))

		_, err := s.r.GetSeasonByYear(i)
		if err != nil && !errors.Is(err, repo.ErrNoSeasonFound) {
			return fmt.Errorf("error getting season by year: %w", err)
		} else if errors.Is(err, repo.ErrNoSeasonFound) {
			season := &models.Season{
				Year: i,
			}
			err = s.r.SaveSeason(season)
			if err != nil {
				return fmt.Errorf("error saving season: %w", err)
			}
		}

		wg.Add(1)
		seasonRacesTask := newTask(func() {
			defer wg.Done()
			err := s.ImportSeasonRaces(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported races for season %d", i))
		})

		wg.Add(1)
		driverChampsTask := newTask(func() {
			defer wg.Done()
			err := s.ImportSeasonDriversChamps(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d drivers championship: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported drivers championship for season %d", i))
		})

		wg.Add(1)
		constructorChampsTask := newTask(func() {
			defer wg.Done()
			err := s.ImportSeasonConstructorsChamps(i)
			if err != nil {
				slog.Error(fmt.Sprintf("error importing season %d constructors championship: %v", i, err))
				return
			}

			slog.Info(fmt.Sprintf("Imported constructors championship for season %d", i))
		})

		if err := workers.Schedule(seasonRacesTask); err != nil {
			return fmt.Errorf("error scheduling season: %w", err)
		}

		if err := workers.Schedule(driverChampsTask); err != nil {
			return fmt.Errorf("error scheduling season: %w", err)
		}

		if err := workers.Schedule(constructorChampsTask); err != nil {
			return fmt.Errorf("error scheduling season: %w", err)
		}
	}

	wg.Wait()

	return nil
}
