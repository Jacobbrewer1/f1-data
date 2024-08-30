package utils

import (
	"reflect"
	"strings"
)

// MultiErrorer is an interface that defines objects that returns multiple errors
type MultiErrorer interface {
	Errors() []error
}

// MultiError is an error type which allows multiple errors to be added and returned as single error.
// This is useful in the case where errors aren't fatal, and you want the combination to be returned at a later date.
type MultiError struct {
	errs ErrorList
}

// NewMultiError returns a usable MultiError
func NewMultiError() *MultiError {
	return &MultiError{
		make([]error, 0),
	}
}

// MultiErrors creates a MultiError and adds all errs to it and returns the resulting error.
func MultiErrors(errs ...error) error {
	merr := NewMultiError()
	for _, err := range errs {
		merr.Add(err)
	}
	return merr.Err()
}

// Add adds the given error if it's not nil and returns true, otherwise it returns false.
func (m *MultiError) Add(err error) bool {
	if err != nil {
		m.errs = append(m.errs, err)
		return true
	}
	return false
}

// Error implements the error interface
func (m *MultiError) Error() string {
	if len(m.errs) == 0 {
		return ""
	}

	out := make([]string, len(m.errs))
	for i, e := range m.errs {
		out[i] = e.Error()
	}

	return strings.Join(out, ". ")
}

// ErrorStrings returns a string array of errors
func (m *MultiError) ErrorStrings() []string {
	if len(m.errs) == 0 {
		return []string{}
	}

	out := make([]string, len(m.errs))
	for i, e := range m.errs {
		out[i] = e.Error()
	}

	return out
}

// Errors returns a flattened list of all errors the MultiError contains
func (m *MultiError) Errors() (errs []error) {
	for _, e := range m.errs {
		if me, ok := e.(MultiErrorer); ok {
			errs = append(errs, me.Errors()...)
		} else {
			errs = append(errs, e)
		}
	}

	return
}

// Err returns either:
// * nil - if no errors where added
// * the added error - if only one error was added
// * itself - otherwise
func (m *MultiError) Err() error {
	switch len(m.errs) {
	case 0:
		return nil
	case 1:
		return m.errs[0]
	}

	return m
}

// Reset clears any previously added errors.
func (m *MultiError) Reset() {
	m.errs = m.errs[:0]
}

func (m *MultiError) Is(err error) bool {
	for _, ie := range m.Errors() {
		if reflect.TypeOf(ie) == reflect.TypeOf(err) {
			return true
		}
	}
	return false
}

// Unwrap allows errors.Is to unwrap a MultiError's list of errors
func (m *MultiError) Unwrap() []error {
	return m.errs
}
