package utils

import (
	"encoding/json"
	"errors"
)

// ErrorList is a JSON that can serialise a list of errors
type ErrorList []error

// MarshalJSON implements the json.Marshaler interface
func (el *ErrorList) MarshalJSON() ([]byte, error) {
	if el == nil {
		return json.Marshal([]string{})
	}

	out := make([]string, len(*el))
	for i, e := range *el {
		out[i] = e.Error()
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (el *ErrorList) UnmarshalJSON(data []byte) error {
	if el == nil {
		return nil
	}

	in := make([]string, 0)
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}

	*el = make([]error, len(in))
	for i, s := range in {
		(*el)[i] = errors.New(s)
	}

	return nil
}
