package filters

type Year struct {
	y int
}

// NewYear creates an instance of the filter.
func NewYear(year int) *Year {
	return &Year{y: year}
}

// Where creates the where query.
func (s *Year) Where() (string, []interface{}) {
	return `AND t.year = ?`, []interface{}{s.y}
}
