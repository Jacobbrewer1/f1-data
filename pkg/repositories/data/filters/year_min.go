package filters

type YearMin struct {
	ym int
}

// NewYearMin creates an instance of the filter.
func NewYearMin(yearMin int) *YearMin {
	return &YearMin{ym: yearMin}
}

// Where creates the where query.
func (s *YearMin) Where() (string, []interface{}) {
	return `AND t.year >= ?`, []interface{}{s.ym}
}
