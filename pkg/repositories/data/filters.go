package data

type GetSeasonsFilters struct {
	// Year is the year of the season.
	Year *int

	// YearMin is the minimum year of the season.
	YearMin *int

	// YearMax is the maximum year of the season.
	YearMax *int
}

type GetSeasonRacesFilters struct {
	// SeasonYear is the year of the season.
	SeasonYear *int
}

type GetRaceResultsFilters struct {
	// RaceID
	RaceID *int

	// SeasonYear is the year of the season.
	SeasonYear *int
}

type GetDriversChampionshipFilters struct {
	// SeasonYear is the year of the season.
	SeasonYear *int

	// DriverName is the name of the driver.
	DriverName *string

	// DriverTag is the tag of the driver.
	DriverTag *string

	// Team is the team of the driver.
	Team *string
}

type GetConstructorsChampionshipFilters struct {
	// SeasonYear is the year of the season.
	SeasonYear *int

	// ConstructorName is the name of the constructor.
	ConstructorName *string
}

type GetDriversFilters struct {
	// Name is the name of the driver.
	Name *string

	// Tag is the tag of the driver.
	Tag *string

	// Team is the team of the driver.
	Team *string

	// Nationality is the nationality of the driver.
	Nationality *string
}
