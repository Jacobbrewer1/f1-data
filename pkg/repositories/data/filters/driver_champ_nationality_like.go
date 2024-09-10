package filters

// DriverChampNationalityLike is a filter that filters by driver championship by nationality.
type DriverChampNationalityLike struct {
	tl string
}

// NewDriverChampNationalityLike creates an instance of the filter.
func NewDriverChampNationalityLike(teamLike string) *DriverChampNationalityLike {
	return &DriverChampNationalityLike{tl: teamLike}
}

// Where creates the where query.
func (s *DriverChampNationalityLike) Where() (string, []interface{}) {
	return `AND t.nationality LIKE ?`, []interface{}{"%" + s.tl + "%"}
}
