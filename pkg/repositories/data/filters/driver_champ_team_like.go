package filters

// DriverChampTeamLike is a filter that filters by driver championship team like.
type DriverChampTeamLike struct {
	tl string
}

// NewDriverChampTeamLike creates an instance of the filter.
func NewDriverChampTeamLike(teamLike string) *DriverChampTeamLike {
	return &DriverChampTeamLike{tl: teamLike}
}

// Where creates the where query.
func (s *DriverChampTeamLike) Where() (string, []interface{}) {
	return `AND t.team LIKE ?`, []interface{}{"%" + s.tl + "%"}
}
