package filters

// DriverChampNameLike is a filter that filters by driver championship name like.
type DriverChampNameLike struct {
	nl string
}

// NewDriverChampNameLike creates an instance of the filter.
func NewDriverChampNameLike(nameLike string) *DriverChampNameLike {
	return &DriverChampNameLike{nl: nameLike}
}

// Where creates the where query.
func (s *DriverChampNameLike) Where() (string, []interface{}) {
	return `AND t.driver LIKE ?`, []interface{}{"%" + s.nl + "%"}
}
