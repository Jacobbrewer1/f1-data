package models

import "errors"

// ErrNoAffectedRows is returned if a model update affected no rows
var ErrNoAffectedRows = errors.New("no affected rows")

// ErrDuplicate is returned if a duplicate entry is found
var ErrDuplicate = errors.New("duplicate entry")

// ErrConstraintViolation is returned if a constraint is violated
var ErrConstraintViolation = errors.New("constraint violation")
