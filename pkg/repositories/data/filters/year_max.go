package filters

type YearMax struct {
	ym int
}

// NewYearMax creates an instance of the filter.
func NewYearMax(yearMax int) *YearMax {
	return &YearMax{ym: yearMax}
}

// Where creates the where query.
func (s *YearMax) Where() (string, []interface{}) {
	return `AND t.year <= ?`, []interface{}{s.ym}
}
