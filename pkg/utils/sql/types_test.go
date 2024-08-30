package sql

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testNullBool struct {
	Field1 NullBool
}

type testNullFloat64 struct {
	Field1 NullFloat64
}

type testNullInt64 struct {
	Field1 NullInt64
}

type testNullString struct {
	Field1 NullString
}

type testNullDuration struct {
	Field1 NullDuration
}

type testNullTime struct {
	Field1 NullTime
}

type testDuration struct {
	Field1 Duration
}

type TypesSuite struct {
	suite.Suite
}

func (s *TypesSuite) TestValidNullBoolUnmarshal() {
	ts := &testNullBool{}
	if err := decodeJSON(`{"Field1": true}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.True(ts.Field1.Bool)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullBoolMarshal() {
	ts := &testNullBool{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullBoolMarshal() {
	ts := &testNullBool{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":false}`, j)

}

func (s *TypesSuite) TestValidTrueNullBoolMarshal() {
	ts := &testNullBool{}
	ts.Field1.Valid = true
	ts.Field1.Bool = true
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":true}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullBoolUnmarshal() {
	ts := &testNullBool{Field1: *NewNullBool(true)}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.False(ts.Field1.Bool)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullBoolUnmarshal() {
	ts := &testNullBool{}
	if err := decodeJSON(`{"Field1": "xxx"}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestValidNullFloat64Unmarshal() {
	ts := &testNullFloat64{}
	if err := decodeJSON(`{"Field1": 1.004}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(1.004, ts.Field1.Float64)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullFloat64Marshal() {
	ts := &testNullFloat64{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullFloat64Marshal() {
	ts := &testNullFloat64{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":0}`, j)

}

func (s *TypesSuite) TestValidSetNullFloat64Marshal() {
	ts := &testNullFloat64{}
	ts.Field1.Valid = true
	ts.Field1.Float64 = 6.4
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":6.4}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullFloat64Unmarshal() {
	ts := &testNullFloat64{Field1: *NewNullFloat64(1.004)}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(float64(0), ts.Field1.Float64)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullFloat64Unmarshal() {
	ts := &testNullFloat64{}
	if err := decodeJSON(`{"Field1": "xxx"}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestValidNullInt64Unmarshal() {
	ts := &testNullInt64{}
	if err := decodeJSON(`{"Field1": 10}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(int64(10), ts.Field1.Int64)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullInt64Marshal() {
	ts := &testNullInt64{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullInt64Marshal() {
	ts := &testNullInt64{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":0}`, j)

}

func (s *TypesSuite) TestValidSetNullInt64Marshal() {
	ts := &testNullInt64{}
	ts.Field1.Valid = true
	ts.Field1.Int64 = 64
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":64}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullInt64Unmarshal() {
	ts := &testNullInt64{Field1: *NewNullInt64(10)}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(int64(0), ts.Field1.Int64)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullInt64Unmarshal() {
	ts := &testNullInt64{}
	if err := decodeJSON(`{"Field1": "xxx"}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestValidNullStringUnmarshal() {
	ts := &testNullString{}
	if err := decodeJSON(`{"Field1": "test"}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal("test", ts.Field1.String)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullStringMarshal() {
	ts := &testNullString{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullStringMarshal() {
	ts := &testNullString{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":""}`, j)
}

func (s *TypesSuite) TestValidSetNullStringMarshal() {
	ts := &testNullString{}
	ts.Field1.Valid = true
	ts.Field1.String = "string"
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":"string"}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullStringUnmarshal() {
	ts := &testNullString{Field1: *NewNullString("test")}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal("", ts.Field1.String)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullStringUnmarshal() {
	ts := &testNullString{}
	if err := decodeJSON(`{"Field1": 0}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestValidNullStringRedisArg() {
	ts := &testNullString{
		Field1: *NewNullString("test"),
	}
	s.True(ts.Field1.Valid)
	s.Equal("test", ts.Field1.RedisArg())
}

func (s *TypesSuite) TestInvalidNullStringRedisArg() {
	ts := &testNullString{}
	s.False(ts.Field1.Valid)
	s.Empty(ts.Field1.RedisArg())
}

func (s *TypesSuite) TestValidNullStringRedisScan() {
	ts := &testNullString{}
	s.NoError(ts.Field1.RedisScan("test"))
	s.True(ts.Field1.Valid)
	s.Equal("test", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullStringRedisScan_bytes() {
	ts := &testNullString{}
	s.NoError(ts.Field1.RedisScan([]byte("test")))
	s.True(ts.Field1.Valid)
	s.Equal("test", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullStringRedisScan_emptyString() {
	ts := &testNullString{}
	s.NoError(ts.Field1.RedisScan(""))
	s.True(ts.Field1.Valid)
	s.Equal("", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullStringRedisScan_emptyBytes() {
	ts := &testNullString{}
	s.NoError(ts.Field1.RedisScan([]byte{}))
	s.True(ts.Field1.Valid)
	s.Equal("", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullStringRedisScan_nil() {
	ts := &testNullString{}
	s.NoError(ts.Field1.RedisScan(nil))
	s.False(ts.Field1.Valid)
	s.Equal("", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullStringRedisScan_invalidType() {
	ts := &testNullString{}
	s.EqualError(ts.Field1.RedisScan(123), "unexpected type: int")
	s.False(ts.Field1.Valid)
	s.Equal("", ts.Field1.String)
}

func (s *TypesSuite) TestValidNullDurationUnmarshal() {
	ts := &testNullDuration{}
	if err := decodeJSON(`{"Field1": "15:04:05"}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(Duration(time.Hour*15+time.Minute*4+time.Second*5), ts.Field1.Duration)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullDurationMarshal() {
	ts := &testNullDuration{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullDurationMarshal() {
	ts := &testNullDuration{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":"00:00:00"}`, j)

}

func (s *TypesSuite) TestValidTrueNullDurationMarshal() {
	ts := &testNullDuration{}
	ts.Field1.Valid = true
	ts.Field1.Duration = Duration(time.Hour + time.Minute + time.Second + time.Nanosecond)
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":"01:01:01.1"}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullDurationUnmarshal() {
	ts := &testNullDuration{Field1: *NewNullDuration(time.Hour + time.Minute + time.Second + time.Nanosecond)}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Zero(ts.Field1.Duration)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullDurationUnmarshal() {
	ts := &testNullDuration{}
	if err := decodeJSON(`{"Field1": "xxx"}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestDurationUnmarshal() {
	ts := &testDuration{}
	if err := decodeJSON(`{"Field1": "15:04:05"}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(Duration(time.Hour*15+time.Minute*4+time.Second*5), ts.Field1)
}

func (s *TypesSuite) TestZeroDurationMarshal() {
	ts := &testDuration{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":"00:00:00"}`, j)
}

func (s *TypesSuite) TestDurationMarshal() {
	ts := &testDuration{}
	ts.Field1 = Duration(time.Hour + time.Minute + time.Second + time.Nanosecond)
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":"01:01:01.1"}`, j)
}

func (s *TypesSuite) TestExistingDurationUnmarshal() {
	ts := &testDuration{}
	ts.Field1 = Duration(time.Hour + time.Minute + time.Second + time.Nanosecond)
	if err := decodeJSON(`{"Field1": "15:04:05"}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Equal(Duration(time.Hour*15+time.Minute*4+time.Second*5), ts.Field1)
}

func (s *TypesSuite) TestValidNullTimeUnmarshal() {
	ts := &testNullTime{}
	if err := decodeJSON(`{"Field1": 979311845}`, ts); err != nil {
		s.T().Fatal(err)
	}

	t := time.Date(2001, 1, 12, 15, 4, 5, 0, time.UTC)
	s.Equal(t, ts.Field1.Time)
	s.True(ts.Field1.Valid)
}

func (s *TypesSuite) TestInvalidNullTimeMarshal() {
	ts := &testNullTime{}

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidNullTimeMarshal() {
	ts := &testNullTime{}
	ts.Field1.Valid = true

	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":null}`, j)
}

func (s *TypesSuite) TestValidTrueNullTimeMarshal() {
	ts := &testNullTime{}
	ts.Field1.Valid = true
	ts.Field1.Time = time.Date(2001, 1, 12, 15, 4, 5, 0, time.UTC)
	j, err := encodeJSON(ts)
	if err != nil {
		s.T().Fatal(err)
	}

	s.Equal(`{"Field1":979311845}`, j)
}

func (s *TypesSuite) TestNullExistingValidNullTimeUnmarshal() {
	ts := &testNullTime{Field1: *NewNullTime(time.Date(2001, 1, 12, 15, 4, 5, 0, time.UTC))}
	if err := decodeJSON(`{"Field1": null}`, ts); err != nil {
		s.T().Fatal(err)
	}

	s.Zero(ts.Field1.Time)
	s.False(ts.Field1.Valid)
}

func (s *TypesSuite) TestNullInvalidNullTimeUnmarshal() {
	ts := &testNullTime{}
	if err := decodeJSON(`{"Field1": "xxx"}`, ts); err == nil {
		s.T().Fatal("should have errored")
	}
}

func (s *TypesSuite) TestValidNullTimeRedisArg() {
	now := time.Now().UTC()
	ts := &testNullTime{
		Field1: *NewNullTime(now),
	}
	s.True(ts.Field1.Valid)
	s.Equal(now.UTC().Unix(), ts.Field1.RedisArg())
}

func (s *TypesSuite) TestInvalidNullTimeRedisArg() {
	ts := &testNullTime{}
	s.False(ts.Field1.Valid)
	s.Empty(ts.Field1.RedisArg())
}

func (s *TypesSuite) TestValidNullTimeRedisScan() {
	now := time.Now().UTC()
	ts := &testNullTime{}
	s.NoError(ts.Field1.RedisScan(strconv.FormatInt(now.Unix(), 10)))
	s.True(ts.Field1.Valid)
	s.Equal(now.Unix(), ts.Field1.Time.Unix())
}

func (s *TypesSuite) TestValidNullTimeRedisScan_bytes() {
	now := time.Now().UTC()
	ts := &testNullTime{}
	s.NoError(ts.Field1.RedisScan([]byte(strconv.FormatInt(now.Unix(), 10))))
	s.True(ts.Field1.Valid)
	s.Equal(now.Unix(), ts.Field1.Time.Unix())
}

func (s *TypesSuite) TestValidNullTimeRedisScan_emptyString() {
	ts := &testNullTime{}
	s.NoError(ts.Field1.RedisScan(""))
	s.False(ts.Field1.Valid)
	s.Equal(time.Time{}, ts.Field1.Time)
}

func (s *TypesSuite) TestValidNullTimeRedisScan_emptyBytes() {
	ts := &testNullTime{}
	s.NoError(ts.Field1.RedisScan([]byte{}))
	s.False(ts.Field1.Valid)
	s.Equal(time.Time{}, ts.Field1.Time)
}

func (s *TypesSuite) TestValidNullTimeRedisScan_nil() {
	ts := &testNullTime{}
	s.NoError(ts.Field1.RedisScan(nil))
	s.False(ts.Field1.Valid)
	s.Equal(time.Time{}, ts.Field1.Time)
}

func (s *TypesSuite) TestValidNullTimeRedisScan_invalidType() {
	ts := &testNullTime{}
	s.EqualError(ts.Field1.RedisScan(123), "unexpected type: int")
	s.False(ts.Field1.Valid)
	s.Equal(time.Time{}, ts.Field1.Time)
}

func TestTypes(t *testing.T) {
	suite.Run(t, &TypesSuite{})
}

func decodeJSON(j string, i interface{}) error {
	return json.NewDecoder(strings.NewReader(j)).Decode(i)
}

func encodeJSON(i interface{}) (string, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(i)
	return strings.TrimSpace(buf.String()), err
}

func nullInt64(v int64) NullInt64 {
	return NullInt64{NullInt64: sql.NullInt64{Int64: v, Valid: true}}
}

func TestInt64RedisScan(t *testing.T) {
	ts := &testNullInt64{}

	cases := []struct {
		name     string
		data     interface{}
		err      bool
		expected NullInt64
	}{
		{
			name:     "valid-string",
			data:     "10",
			expected: nullInt64(10),
		},
		{
			name:     "valid-bytes",
			data:     []byte("10"),
			expected: nullInt64(10),
		},
		{
			name: "nil",
		},
		{
			name: "invalid-string",
			data: "invalid",
			err:  true,
		},
		{
			name: "invalid-bytes",
			data: []byte("invalid"),
			err:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ts.Field1.RedisScan(tc.data)
			if tc.err {
				require.Error(t, err)
				require.Equal(t, tc.expected, ts.Field1)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
