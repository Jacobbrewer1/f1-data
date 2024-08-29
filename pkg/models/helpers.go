package models

import (
	"fmt"
	"log/slog"
	"reflect"
)

// Saveable is the interface implemented by types which can save themselves to the database.
type Saveable interface {
	Save(db DB) error
}

// PreSaveable is the interface implemented by types which run a pre save step.
type PreSaveable interface {
	PreSave(db DB) error
}

// PostSaveable is the interface implemented by types which run a post save step.
type PostSaveable interface {
	PostSave() error
}

// SetLogger is the interface implemented by types which have the ability to configure their log entry.
type SetLogger interface {
	SetLog(l *slog.Logger)
}

// Deletable is the interface implemented by types which can delete themselves from the database.
type Deletable interface {
	Delete(db DB) error
}

// PreDeletable is the interface implemented by types which run a pre delete step.
type PreDeletable interface {
	PreDelete() error
}

// PostDeletable is the interface implemented by types which run a post delete step.
type PostDeletable interface {
	PostDelete() error
}

// TransactionFunc is a function to be called within a transaction.
type TransactionFunc func(db DB) error

type TransactionHandler interface {
	Handle(TransactionFunc) error
}

// DBTransactionHandler handles a transaction that will return any error.
type DBTransactionHandler struct {
	db Transactioner
}

// Handle implements the TransactionHandler interface.
func (th *DBTransactionHandler) Handle(f TransactionFunc) error {
	tx, err := th.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	if err := f(tx); err != nil {
		if err2 := tx.Rollback(); err2 != nil {
			return fmt.Errorf("%s: %w", err, err2)
		}
		return fmt.Errorf("action: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

// NewDBTransactionHandler returns a configured instance of DBTransactionHandler
func NewDBTransactionHandler(db Transactioner) *DBTransactionHandler {
	return &DBTransactionHandler{db: db}
}

// LoggableDBTransactionHandler handles a transaction and logs any error.
type LoggableDBTransactionHandler struct {
	db Transactioner
	l  *slog.Logger
}

// NewLoggableDBTransactionHandler returns a configured instance of LoggableDBTransactionHandler
func NewLoggableDBTransactionHandler(db Transactioner, l *slog.Logger) *LoggableDBTransactionHandler {
	return &LoggableDBTransactionHandler{db: db, l: l}
}

// IsKeySet returns true if
// 1. x is an integer and greater than zero.
// 2. x not an integer and is not the zero value.
// Otherwise, returns false
func IsKeySet(x interface{}) bool {
	switch x := x.(type) {
	case int:
		return x > 0
	case int8:
		return x > 0
	case int16:
		return x > 0
	case int32:
		return x > 0
	case int64:
		return x > 0
	case uint:
		return x > 0
	case uint8:
		return x > 0
	case uint16:
		return x > 0
	case uint32:
		return x > 0
	case uint64:
		return x > 0
	}

	return x != reflect.Zero(reflect.TypeOf(x)).Interface()
}
