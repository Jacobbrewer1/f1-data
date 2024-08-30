package data

type GetSeasonsFilters struct {
	// Year is the year of the season.
	Year *int `json:"year"`

	// YearMin is the minimum year of the season.
	YearMin *int `json:"year_min"`

	// YearMax is the maximum year of the season.
	YearMax *int `json:"year_max"`
}

type GetSeasonRacesFilters struct {
	// SeasonYear is the year of the season.
	SeasonYear *int `json:"season_year"`
}
