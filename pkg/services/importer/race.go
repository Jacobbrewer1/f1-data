package importer

import (
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/models"
	repo "github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"
	"github.com/gocolly/colly/v2"
)

func (s *service) processRace(raceId int, url *url.URL) error {
	c := colly.NewCollector()

	// The order that the headers come through in the HTML.
	//
	// 0: "Pos"
	// 1: "No"
	// 2: "Driver"
	// 3: "Car"
	// 4: "Laps"
	// 5: "Time/Retired"
	// 6: "Pts"

	i := 0

	c.OnHTML("tbody", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			raceResult := new(models.RaceResult)

			el.ForEach("p", func(_ int, elm *colly.HTMLElement) {
				// Get the value in the span tags.
				child := elm.ChildTexts("span")

				text := elm.Text
				if len(child) > 0 {
					text = child[0] + " " + child[1]
				}

				if len(child) >= 2 {
					raceResult.DriverTag = child[2]
				}

				var err error
				switch i {
				case 0:
					raceResult.Position = text
				case 1:
					raceResult.DriverNumber, err = strconv.Atoi(text)
					if err == nil {
						foundRaceResult, rrErr := s.r.GetRaceResultByRaceIdAndDriverNumber(raceId, raceResult.DriverNumber)
						if rrErr != nil && !errors.Is(rrErr, repo.ErrNoRaceResultFound) {
							slog.Error("Error getting existing race result", slog.String(logging.KeyError, rrErr.Error()))
							break
						} else if errors.Is(rrErr, repo.ErrNoRaceResultFound) {
							err = nil
							break
						}

						raceResult.Id = foundRaceResult.Id
					}
				case 2:
					raceResult.Driver = text
				case 3:
					raceResult.Team = text
				case 4:
					if text == "" {
						text = "-1"
					}
					raceResult.Laps, err = strconv.Atoi(text)
				case 5:
					raceResult.TimeRetired = text
				case 6:
					raceResult.Points, err = strconv.ParseFloat(text, 64)
				}

				if err != nil {
					slog.Error("Error parsing race result", slog.String(logging.KeyError, err.Error()), slog.Int("index", i))
				}

				i++
			})
			i = 0

			raceResult.RaceId = raceId

			err := s.r.SaveRaceResult(raceResult)
			if err != nil {
				slog.Error("Error saving race result", slog.String(logging.KeyError, err.Error()))
			}
		})
	})

	err := c.Visit(url.String())
	if err != nil {
		return fmt.Errorf("error visiting URL: %w", err)
	}

	return nil
}
