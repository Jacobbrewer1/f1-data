package pagefilter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testFilter struct{}

func (testFilter) Join() (string, []interface{}) {
	return "\n\njoin\n\n", []interface{}{1, 2, 3}
}
func (testFilter) Where() (string, []interface{}) {
	return "where", []interface{}{4, 5, 6}
}

func TestMultiFilter(t *testing.T) {
	mf := NewMultiFilter()
	mf.Add(testFilter{})
	mf.Add(testFilter{})
	mf.Add(testFilter{})

	jSQL, jArgs := mf.Join()
	require.Equal(t, "join\njoin\njoin", jSQL)
	require.Equal(t, []interface{}{1, 2, 3, 1, 2, 3, 1, 2, 3}, jArgs)

	wSQL, wArgs := mf.Where()
	require.Equal(t, "where\nwhere\nwhere", wSQL)
	require.Equal(t, []interface{}{4, 5, 6, 4, 5, 6, 4, 5, 6}, wArgs)
}
