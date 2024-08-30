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

type GetRaceResultsFilters struct {
	// RaceID
	RaceID *int

	// SeasonYear is the year of the season.
	SeasonYear *int `json:"season_year"`
}

type GetDriversChampionshipFilters struct {
	// SeasonYear is the year of the season.
	SeasonYear *int `json:"season_year"`

	// DriverName is the name of the driver.
	DriverName *string `json:"driver_name"`

	// DriverTag is the tag of the driver.
	DriverTag *string `json:"driver_tag"`

	// Team is the team of the driver.
	Team *string `json:"team"`
}
