package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"
	"github.com/gocolly/colly/v2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func (s *service) ImportSeasonRaces(year int) error {
	season, err := s.r.GetSeasonByYear(year)
	if err != nil {
		return fmt.Errorf("error getting season by year: %w", err)
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

		newRace := new(models.Race)
		existingRace, err := s.r.GetRaceBySeasonIdAndGrandPrix(season.Id, e.Text)
		if err != nil && !errors.Is(err, repo.ErrRaceNotFound) {
			slog.Error("Error getting race", slog.String(logging.KeyError, err.Error()))
			return
		} else if existingRace != nil {
			newRace.Id = existingRace.Id
		}

		newRace.SeasonId = season.Id
		newRace.GrandPrix = e.Text

		<-datesCollected

		if raceDate, ok := raceDates[e.Text]; ok {
			newRace.Date = *raceDate
		} else {
			slog.Error("Error getting race date from map, assuming the hasn't happened or didnt happen", slog.String("grand_prix", e.Text))
			return
		}

		opts := cmpopts.IgnoreFields(models.Race{}, "Id")
		if diff := cmp.Diff(existingRace, newRace, opts); diff == "" {
			slog.Debug("Races are the same",
				slog.String("grand_prix", e.Text),
				slog.String("year", fmt.Sprintf("%d", year)),
				slog.String("diff", diff),
			)
		} else {
			slog.Debug("Races are different",
				slog.String("grand_prix", e.Text),
				slog.String("year", fmt.Sprintf("%d", year)),
				slog.String("diff", diff),
			)

			err = s.r.SaveRace(newRace)
			if err != nil {
				slog.Error("Error saving race", slog.String(logging.KeyError, err.Error()))
				return
			}
		}

		raceUrl, uErr := url.Parse(s.baseUrl + link)
		if uErr != nil {
			slog.Error("Error parsing URL", slog.String(logging.KeyError, err.Error()))
		}

		err = s.processRace(newRace.Id, raceUrl)
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

func (s *service) ImportSeasonDriversChamps(year int) error {
	season, err := s.r.GetSeasonByYear(year)
	if err != nil {
		return fmt.Errorf("error getting season by year: %w", err)
	}

	urlFmt := fmt.Sprintf("%s/en/results/%d/drivers", s.baseUrl, year)
	u, err := url.Parse(urlFmt)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	c := colly.NewCollector()

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			children := el.ChildTexts("td")
			if len(children) < 5 {
				return
			}

			driverNames := el.ChildTexts("span")
			if driverNames == nil || len(driverNames) < 2 {
				slog.Error("Error getting driver names")
				return
			}
			driverName := driverNames[0] + " " + driverNames[1]

			newDriver := new(models.DriverChampionship)
			existingDriver, err := s.r.GetDriverByName(season.Id, driverName)
			if err != nil && !errors.Is(err, repo.ErrDriverNotFound) {
				slog.Error("Error getting driver", slog.String(logging.KeyError, err.Error()))
				return
			} else if existingDriver != nil {
				newDriver.Id = existingDriver.Id
			}

			newDriver.SeasonId = season.Id
			newDriver.Driver = driverName
			newDriver.DriverTag = driverNames[2]
			newDriver.Nationality = children[2]
			newDriver.Team = children[3]

			position, err := strconv.Atoi(children[0])
			if err != nil {
				slog.Error("Error converting position to int", slog.String(logging.KeyError, err.Error()))
				return
			}
			newDriver.Position = position

			points, err := strconv.ParseFloat(children[4], 64)
			if err != nil {
				slog.Error("Error converting points to int", slog.String(logging.KeyError, err.Error()))
				return
			}
			newDriver.Points = points

			opts := cmpopts.IgnoreFields(models.DriverChampionship{}, "Id")
			if diff := cmp.Diff(existingDriver, newDriver, opts); diff == "" {
				slog.Debug("Drivers are the same",
					slog.String("driver", driverName),
					slog.String("year", fmt.Sprintf("%d", year)),
					slog.String("diff", diff),
				)
			} else {
				slog.Debug("Drivers are different",
					slog.String("driver", driverName),
					slog.String("year", fmt.Sprintf("%d", year)),
					slog.String("diff", diff),
				)

				err = s.r.SaveDriver(newDriver)
				if err != nil {
					slog.Error("Error saving driver", slog.String(logging.KeyError, err.Error()))
					return
				}
			}
		})
	})

	err = c.Visit(u.String())
	if err != nil {
		return fmt.Errorf("error visiting URL: %w", err)
	}

	return nil
}

func (s *service) ImportSeasonConstructorsChamps(year int) error {
	season, err := s.r.GetSeasonByYear(year)
	if err != nil {
		return fmt.Errorf("error getting season by year: %w", err)
	}

	urlFmt := fmt.Sprintf("%s/en/results/%d/team", s.baseUrl, year)
	u, err := url.Parse(urlFmt)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	c := colly.NewCollector()

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			children := el.ChildTexts("td")
			if len(children) < 3 {
				return
			}

			newConstructor := new(models.ConstructorChampionship)
			existingConstructor, err := s.r.GetConstructorByName(season.Id, children[1])
			if err != nil && !errors.Is(err, repo.ErrConstructorNotFound) {
				slog.Error("Error getting driver", slog.String(logging.KeyError, err.Error()))
				return
			} else if existingConstructor != nil {
				newConstructor.Id = existingConstructor.Id
			}

			newConstructor.SeasonId = season.Id
			newConstructor.Name = children[1]

			position, err := strconv.Atoi(children[0])
			if err != nil {
				slog.Error("Error converting position to int", slog.String(logging.KeyError, err.Error()))
				return
			}
			newConstructor.Position = position

			points, err := strconv.ParseFloat(children[2], 64)
			if err != nil {
				slog.Error("Error converting points to int", slog.String(logging.KeyError, err.Error()))
				return
			}
			newConstructor.Points = points

			opts := cmpopts.IgnoreFields(models.ConstructorChampionship{}, "Id")
			if diff := cmp.Diff(existingConstructor, newConstructor, opts); diff == "" {
				slog.Debug("Constructors are the same",
					slog.String("constructor", children[1]),
					slog.String("year", fmt.Sprintf("%d", year)),
					slog.String("diff", diff),
				)
			} else {
				slog.Debug("Constructors are different",
					slog.String("constructor", children[1]),
					slog.String("year", fmt.Sprintf("%d", year)),
					slog.String("diff", diff),
				)

				err = s.r.SaveConstructor(newConstructor)
				if err != nil {
					slog.Error("Error saving driver", slog.String(logging.KeyError, err.Error()))
					return
				}
			}
		})
	})

	err = c.Visit(u.String())
	if err != nil {
		return fmt.Errorf("error visiting URL: %w", err)
	}

	return nil
}
