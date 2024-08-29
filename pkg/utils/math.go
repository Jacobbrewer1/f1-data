package utils

import "math"

// Round rounds a float to a specified precision.
// The precision argument is the number of decimal places to round to.
// For example, Round(1.2345, 2) returns 1.23.
// The value argument can be a float32 or float64.
// The return value is the rounded float to the T type.
func Round[T float32 | float64](value T, precision int) T {
	shift := math.Pow(10, float64(precision))
	return T(math.Round(float64(value)*shift) / shift)
}
