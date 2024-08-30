package pagefilter

import (
	"fmt"
)

// InvalidParamError is returned when the param can never be valid
type InvalidParamError string

func (e InvalidParamError) Error() string {
	return fmt.Sprintf("passing %s is never valid", string(e))
}

// MissingValueError is returned when the param is missing a value
type MissingValueError string

func (e MissingValueError) Error() string {
	return fmt.Sprintf("missing %s filter value", string(e))
}

// InvalidOpError is returned when the operation specified in the param isn't valid
type InvalidOpError string

func (e InvalidOpError) Error() string {
	return fmt.Sprintf("op not implemented: %s", string(e))
}
