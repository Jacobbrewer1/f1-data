package filters

type SeasonYearMin struct {
	ym int
}

// NewSeasonYearMin creates an instance of the filter.
func NewSeasonYearMin(yearMin int) *SeasonYearMin {
	return &SeasonYearMin{ym: yearMin}
}

// Where creates the where query.
func (s *SeasonYearMin) Where() (string, []interface{}) {
	return `t.year >= ?`, []interface{}{s.ym}
}
