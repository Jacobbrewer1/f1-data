package sql

import (
	"fmt"
	"regexp"
	"time"
)

type DateTime struct {
	time.Time
}

func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t}
}

// MarshalJSON implements the json.Marshaler interface.
func (d *DateTime) MarshalJSON() ([]byte, error) {
	// Marshal the time.
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *DateTime) UnmarshalJSON(text []byte) error {
	// Remove " from text if present with regex (e.g. "2020-01-01T00:00:00Z" -> 2020-01-01T00:00:00Z)
	reg := regexp.MustCompile(`"(.*)"`)
	text = reg.ReplaceAll(text, []byte("$1"))

	// Parse the time.
	t, err := time.Parse(time.RFC3339, string(text))
	if err != nil {
		return fmt.Errorf("%s is not in the RFC3339 format", text)
	}
	*d = DateTime{t}
	return nil
}

// Scan implements the sql.Scanner interface.
func (d *DateTime) Scan(src any) error {
	switch t := src.(type) {
	case time.Time:
		*d = DateTime{t}
	case string:
		// Parse the time.
		parsedT, err := time.Parse(time.RFC3339, src.(string))
		if err != nil {
			return fmt.Errorf("%s is not in the RFC3339 format", t)
		}
		*d = DateTime{parsedT}
	case []uint8:
		// Parse the time.
		parsedT, err := time.Parse(time.DateTime, string(src.([]uint8)))
		if err != nil {
			return fmt.Errorf("%s is not in the RFC3339 format", t)
		}
		*d = DateTime{parsedT}
	default:
		return fmt.Errorf("unsupported type %T", src)
	}
	return nil
}

// String implements the fmt.Stringer interface.
func (d DateTime) String() string {
	return d.Format(time.RFC3339)
}
