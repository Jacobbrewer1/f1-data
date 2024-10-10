package filters

type DriverChampSeasonYear struct {
	y int
}

// NewDriverChampSeasonYear creates an instance of the filter.
func NewDriverChampSeasonYear(year int) *DriverChampSeasonYear {
	return &DriverChampSeasonYear{y: year}
}

// Join creates the join query.
func (s *DriverChampSeasonYear) Join() (string, []interface{}) {
	return `
		JOIN season s ON s.id = t.season_id
`, []interface{}{}
}

// Where creates the where query.
func (s *DriverChampSeasonYear) Where() (string, []interface{}) {
	return `s.year = ?`, []interface{}{s.y}
}
