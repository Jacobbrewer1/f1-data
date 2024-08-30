package filters

// DriverChampTagLike is a filter that filters by driver championship name like.
type DriverChampTagLike struct {
	t string
}

// NewDriverChampTagLike creates an instance of the filter.
func NewDriverChampTagLike(tagLike string) *DriverChampTagLike {
	return &DriverChampTagLike{t: tagLike}
}

// Where creates the where query.
func (s *DriverChampTagLike) Where() (string, []interface{}) {
	return `AND t.driver_tag LIKE ?`, []interface{}{"%" + s.t + "%"}
}
