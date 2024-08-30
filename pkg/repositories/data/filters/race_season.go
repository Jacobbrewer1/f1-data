package filters

type RaceYear struct {
	ry int
}

// NewRaceYear creates an instance of the filter.
func NewRaceYear(year int) *RaceYear {
	return &RaceYear{ry: year}
}

// Join creates the join query.
func (s *RaceYear) Join() (string, []interface{}) {
	return `
		INNER JOIN season s ON s.id = t.season_id
`, nil
}

// Where creates the where query.
func (s *RaceYear) Where() (string, []interface{}) {
	return `AND s.year = ?`, []interface{}{s.ry}
}
