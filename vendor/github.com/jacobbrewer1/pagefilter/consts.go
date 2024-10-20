package pagefilter

const (
	// Equal is the equal comparison operator
	Equal = "eq"

	// LessThan is the less than comparison operator
	LessThan = "lt"

	// GreaterThan is the greater than comparison operator
	GreaterThan = "gt"

	// Like is the partial match operator
	Like = "like"

	defaultPageLimit = 100
	maxLimit         = 20000

	orderAsc  = "asc"
	orderDesc = "desc"

	sqlComparatorAsc  = "ASC"
	sqlComparatorDesc = "DESC"
	sqlOperatorAsc    = ">"
	sqlOperatorDesc   = "<"

	// DefaultSortBy is the default sort by key
	DefaultSortBy = "id"

	QueryLastVal = "last_val"
	QueryLastID  = "last_id"
	QuerySortBy  = "sort_by"
	QuerySortDir = "sort_dir"
	QueryLimit   = "limit"

	dbTagAutoIncrement = "autoinc"
	dbTagPrimaryKey    = "pk"
	dbTagDefault       = "default"
)
