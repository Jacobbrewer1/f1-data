package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBTransactioner is the interface that database connections that can utilise
// transactions should implement.
type DBTransactioner interface {
	DB
	Transactioner
}

// DB is the common interface for database operations
//
// This should work with database/sql.DB and database/sql.Tx.
type DB interface {
	Exec(string, ...any) (sql.Result, error)
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	// Additional sqlx methods we like
	Get(dest any, query string, args ...interface{}) error
	Select(dest any, query string, args ...interface{}) error
}

// Transactioner is the interface that a database connection that can start
// a transaction should implement.
type Transactioner interface {
	Beginx() (*sqlx.Tx, error)
}

// XODB is a compat alias to DB
type XODB = DB

// DBLog provides the log func used by generated queries.
var DBLog = func(string, ...any) {}

// XOLog is a compat shim for DBLog
var XOLog = func(msg string, args ...any) {
	DBLog(msg, args...)
}
