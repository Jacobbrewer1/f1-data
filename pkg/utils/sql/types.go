package sql

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var null = []byte{0x6e, 0x75, 0x6c, 0x6c}

// NullBool represents a nullable bool type which supports json Marshaler, sql Scanner, and sql driver Valuer interfaces.
type NullBool struct {
	sql.NullBool
}

// MarshalJSON implements the json.Marshaler interface for a NullBool.
func (r NullBool) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Bool)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullBool.
func (r *NullBool) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.Bool = false
		r.Valid = false

		return nil
	}

	if err := json.Unmarshal(data, &r.Bool); err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// NewNullBool returns a valid new NullBool for a given boolean value.
func NewNullBool(b bool) *NullBool {
	return &NullBool{NullBool: sql.NullBool{Bool: b, Valid: true}}
}

// NullFloat64 represents a nullable float64 type which supports json Marshaler, sql Scanner, and sql driver Valuer interfaces.
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON implements the json.Marshaler interface for a NullFloat64.
func (r NullFloat64) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Float64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullFloat64.
func (r *NullFloat64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.Float64 = 0
		r.Valid = false

		return nil
	}

	if err := json.Unmarshal(data, &r.Float64); err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// NewNullFloat64 returns a valid new NullFloat64 for a given float64 value.
func NewNullFloat64(f float64) *NullFloat64 {
	return &NullFloat64{NullFloat64: sql.NullFloat64{Float64: f, Valid: true}}
}

// NullInt64 represents a nullable int64 type which supports json Marshaler, sql Scanner, and sql driver Valuer interfaces.
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON implements the json.Marshaler interface for a NullInt64.
func (r NullInt64) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Int64)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullInt64.
func (r *NullInt64) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.Int64 = 0
		r.Valid = false

		return nil
	}

	if err := json.Unmarshal(data, &r.Int64); err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// RedisArg implements redis.Argument.
//
// The caller should explicitly check the Valid field when putting NullString
// values into Redis. If this is not done it can lead to confusion as RedisScan
// on this type will treat an empty non-nil string as valid.
func (r NullInt64) RedisArg() interface{} {
	return r.Int64
}

// RedisScan implements redis.Scanner.
func (r *NullInt64) RedisScan(src interface{}) error {
	if src == nil {
		r.Int64, r.Valid = 0, false
		return nil
	}

	var err error
	switch src := src.(type) {
	case []byte:
		r.Int64, err = strconv.ParseInt(string(src), 10, 64)
	case string:
		r.Int64, err = strconv.ParseInt(src, 10, 64)
	default:
		return fmt.Errorf("unexpected type: %T", src)
	}

	if err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// NewNullInt64 returns a valid new NullInt64 for a given int64 value.
func NewNullInt64(i int64) *NullInt64 {
	return &NullInt64{NullInt64: sql.NullInt64{Int64: i, Valid: true}}
}

// NullString represents a nullable string type which supports json Marshaler, sql Scanner, and sql driver Valuer interfaces.
type NullString struct {
	sql.NullString
}

// MarshalJSON implements the json.Marshaler interface for a NullString.
func (r NullString) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.String)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullString.
func (r *NullString) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.String = ""
		r.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &r.String); err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// RedisArg implements redis.Argument.
//
// The caller should explicitly check the Valid field when putting NullString
// values into Redis. If this is not done it can lead to confusion as RedisScan
// on this type will treat an empty non-nil string as valid.
func (r NullString) RedisArg() interface{} {
	return r.String
}

// RedisScan implements redis.Scanner.
func (r *NullString) RedisScan(src interface{}) error {
	if src == nil {
		r.String, r.Valid = "", false
		return nil
	}

	switch src := src.(type) {
	case []byte:
		r.String = string(src)
	case string:
		r.String = src
	default:
		return fmt.Errorf("unexpected type: %T", src)
	}
	r.Valid = true

	return nil
}

// NewNullString returns a valid new NullString for a given string value.
func NewNullString(s string) *NullString {
	return &NullString{NullString: sql.NullString{String: s, Valid: true}}
}

// Enum represents non-nullable enum values in sql
type Enum string

func NewEnum(s string) Enum {
	return Enum(s)
}

// NullEnum represents nullable enum values in sql
type NullEnum NullString

// NewNullEnum return a valid new NullEnum for a given value
func NewNullEnum(s string) *NullEnum {
	return &NullEnum{NullString: sql.NullString{String: s, Valid: true}}
}

// Duration represents a mysql TIME column
type Duration time.Duration

func parseTime(str string, loc *time.Location) (time.Time, error) {
	parts := strings.Split(str, ":")
	for len(parts) != 3 {
		parts = append([]string{"00"}, parts...)
	}
	str = strings.Join(parts, ":")

	return time.ParseInLocation("15:04:05.999999999", str, loc)
}

// Scan implements the Scanner interface
func (d *Duration) Scan(value interface{}) (err error) {
	var t time.Time
	switch v := value.(type) {
	case time.Time:
		t = v
	case []byte:
		t, err = parseTime(string(v), time.UTC)
	case string:
		t, err = parseTime(v, time.UTC)
	default:
		return fmt.Errorf("can't convert %T to time.Duration", value)
	}

	if err != nil {
		return err
	}

	h, m, s := t.Clock()
	ns := t.Nanosecond()
	*d = Duration((time.Hour * time.Duration(h)) + (time.Minute * time.Duration(m)) + (time.Second * time.Duration(s)) + time.Duration(ns))

	return nil
}

// Value implements the driver Valuer interface
func (d Duration) Value() (driver.Value, error) {
	var t time.Time
	t = t.Add(time.Duration(d))
	h, m, s := t.Clock()
	str := fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	if t.Nanosecond() > 0 {
		str = fmt.Sprintf("%v.%d", str, t.Nanosecond())
	}
	return str, nil
}

// MarshalJSON implements the json.Marshaler interface
func (d Duration) MarshalJSON() ([]byte, error) {
	s, _ := d.Value()
	return json.Marshal(fmt.Sprint(s))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (d *Duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	return d.Scan(str)
}

// NullDuration represents a nullable string type which supports json Marshaler, sql Scanner, and sql driver Valuer interfaces.
type NullDuration struct {
	Duration
	Valid bool
}

// MarshalJSON implements the json.Marshaler interface for a NullDuration.
func (r NullDuration) MarshalJSON() ([]byte, error) {
	if r.Valid {
		return json.Marshal(r.Duration)
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullDuration.
func (r *NullDuration) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.Duration = 0
		r.Valid = false
		return nil
	}

	if err := json.Unmarshal(data, &r.Duration); err != nil {
		return err
	}
	r.Valid = true

	return nil
}

// NewNullDuration returns a valid new NullDuration for a given time.Duration value.
func NewNullDuration(t time.Duration) *NullDuration {
	return &NullDuration{Duration: Duration(t), Valid: true}
}

// NullTime represents a nullable string type which supports json.Marshaler,
// sql.Scanner, and sql/driver.Valuer interfaces.
type NullTime struct {
	sql.NullTime
}

// MarshalJSON implements json.Marshaller.
func (r NullTime) MarshalJSON() ([]byte, error) {
	if r.Valid && !r.Time.IsZero() {
		return []byte(strconv.FormatInt(r.Time.UTC().Unix(), 10)), nil
	}

	return json.Marshal(nil)
}

// UnmarshalJSON implements the json.Unmarshaler interface for a NullTime.
func (r *NullTime) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, null) {
		r.Time = time.Time{}
		r.Valid = false
		return nil
	}

	return r.parse(string(data))
}

// RedisArg implements redis.Argument.
//
// The time will be serialized as a UTC unix timestamp.
func (r NullTime) RedisArg() interface{} {
	if r.Valid {
		return r.Time.UTC().Unix()
	}
	return nil
}

// RedisScan implements redis.Scanner.
//
// The src value is expected to be a UTC timestamp formatted as a unix timestamp.
func (r *NullTime) RedisScan(src interface{}) error {
	if src == nil {
		return nil
	}

	switch src := src.(type) {
	case []byte:
		if err := r.parse(string(src)); err != nil {
			return err
		}
	case string:
		if err := r.parse(src); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unexpected type: %T", src)
	}

	return nil
}

func (r *NullTime) parse(src string) error {
	if src == "" {
		return nil
	}
	val, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		return err
	}
	r.Time = time.Unix(val, 0).UTC()
	r.Valid = true
	return nil
}

// NewNullTime returns a valid new NullTime for a given time.Time value.
func NewNullTime(t time.Time) *NullTime {
	return &NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}
