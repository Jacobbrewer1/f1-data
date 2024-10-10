package filters

// DriverChampNationalityLike is a filter that filters by driver championship by nationality.
type DriverChampNationalityLike struct {
	nl string
}

// NewDriverChampNationalityLike creates an instance of the filter.
func NewDriverChampNationalityLike(nationalityLike string) *DriverChampNationalityLike {
	return &DriverChampNationalityLike{nl: nationalityLike}
}

// Where creates the where query.
func (s *DriverChampNationalityLike) Where() (string, []interface{}) {
	return `t.nationality LIKE ?`, []interface{}{"%" + s.nl + "%"}
}
