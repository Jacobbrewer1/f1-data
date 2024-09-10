package filters

// DriverChampNameGrouper is a filter that filters by driver championship name like.
type DriverChampNameGrouper struct{}

// NewDriverChampNameGrouper creates an instance of the filter.
func NewDriverChampNameGrouper() *DriverChampNameGrouper {
	return &DriverChampNameGrouper{}
}

// Group creates the group query.
func (s *DriverChampNameGrouper) Group() []string {
	return []string{"t.driver"}
}
