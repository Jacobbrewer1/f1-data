package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

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

	raceDates := make(map[string]*time.Time)
	datesCollected := make(chan struct{}) // Used to enforce order of operations.

	c := colly.NewCollector()

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			children := el.ChildTexts("td")
			if len(children) < 2 {
				return
			}

			t, err := parseDate(children[1])
			if err != nil {
				slog.Error("Error parsing date", slog.String(logging.KeyError, err.Error()))
				return
			}

			raceDates[children[0]] = t
		})

		slog.Debug("Dates collected, closing channel")
		close(datesCollected)
	})

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

		<-datesCollected

		if raceDate, ok := raceDates[e.Text]; ok {
			race.Date = *raceDate
		} else {
			slog.Error("Error getting race date from map, assuming the hasn't happened or didnt happen", slog.String("grand_prix", e.Text))
			return
		}

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

func parseDate(date string) (*time.Time, error) {
	t, err := time.Parse("02 Jan 2006", date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}
	return &t, nil
}
