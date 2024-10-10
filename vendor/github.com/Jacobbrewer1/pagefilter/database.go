package pagefilter

// DB is the common interface for database operations
//
// This should work with database/sql.DB and database/sql.Tx.
type DB interface {
	Get(dest any, query string, args ...any) error
	Select(dest any, query string, args ...any) error
}
