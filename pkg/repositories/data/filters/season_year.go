package filters

type SeasonYear struct {
	y int
}

// NewSeasonYear creates an instance of the filter.
func NewSeasonYear(year int) *SeasonYear {
	return &SeasonYear{y: year}
}

// Where creates the where query.
func (s *SeasonYear) Where() (string, []interface{}) {
	return `AND t.year = ?`, []interface{}{s.y}
}
