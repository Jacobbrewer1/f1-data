package filters

type SeasonYearMax struct {
	ym int
}

// NewSeasonYearMax creates an instance of the filter.
func NewSeasonYearMax(yearMax int) *SeasonYearMax {
	return &SeasonYearMax{ym: yearMax}
}

// Where creates the where query.
func (s *SeasonYearMax) Where() (string, []interface{}) {
	return `t.year <= ?`, []interface{}{s.ym}
}
