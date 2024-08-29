package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"
	"github.com/gocolly/colly/v2"
)

func (s *service) ImportSeason(year int) error {
	season, err := s.r.GetSeasonByYear(year)
	if err != nil && !errors.Is(err, repo.ErrNoSeasonFound) {
		return fmt.Errorf("error getting season by year: %w", err)
	} else if errors.Is(err, repo.ErrNoSeasonFound) {
		season = &models.Season{
			Year: year,
		}
		err = s.r.SaveSeason(season)
		if err != nil {
			return fmt.Errorf("error saving season: %w", err)
		}
	}

	urlFmt := fmt.Sprintf("%s/en/results/%d/races", s.baseUrl, year)
	u, err := url.Parse(urlFmt)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	c := colly.NewCollector()

	c.OnHTML("a", func(e *colly.HTMLElement) {
		// Get the href attribute of the element.
		link := e.Attr("href")

		// Only continue if the link contains "race-results".
		if link == "" {
			return
		} else if !strings.HasSuffix(link, "/race-result") || !strings.HasPrefix(link, "/en/results/") {
			return
		} else if strings.Contains(e.Text, "Season") || strings.Contains(e.Text, "Archive") {
			return
		}

		slog.Debug(fmt.Sprintf("Importing Grand Prix %s", e.Text))

		race, err := s.r.GetRaceBySeasonIdAndGrandPrix(season.Id, e.Text)
		if err != nil && !errors.Is(err, repo.ErrRaceNotFound) {
			slog.Error("Error getting race", slog.String(logging.KeyError, err.Error()))
			return
		} else if errors.Is(err, repo.ErrRaceNotFound) {
			race = new(models.Race)
		}
		race.SeasonId = season.Id
		race.GrandPrix = e.Text

		err = s.r.SaveRace(race)
		if err != nil {
			slog.Error("Error saving race", slog.String(logging.KeyError, err.Error()))
			return
		}

		raceUrl, err := url.Parse(s.baseUrl + link)
		err = s.processRace(race.Id, raceUrl)
		if err != nil {
			slog.Error("Error processing race results", slog.String(logging.KeyError, err.Error()))
			return
		}

		slog.Debug(fmt.Sprintf("Import complete for Grand Prix %s", e.Text))
	})

	err = c.Visit(u.String())
	if err != nil {
		return fmt.Errorf("error visiting URL: %w", err)
	}

	return nil
}
