package filters

type RaceResultSeasonYear struct {
	y int
}

// NewRaceResultSeasonYear creates a new RaceResultSeasonYear.
func NewRaceResultSeasonYear(y int) *RaceResultSeasonYear {
	return &RaceResultSeasonYear{
		y: y,
	}
}

// Join returns the join clause for the filter.
func (r *RaceResultSeasonYear) Join() (string, []interface{}) {
	return `
		JOIN race r ON r.id = t.race_id
		JOIN season s ON s.id = r.season_id
`, nil
}

// Where returns the where clause for the filter.
func (r *RaceResultSeasonYear) Where() (string, []interface{}) {
	return `s.year = ?`, []interface{}{r.y}
}
