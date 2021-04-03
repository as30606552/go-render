package mathutils

import "math"

// Returns the maximum value among all parameters.
func Max(values ...float64) float64 {
	if len(values) == 0 {
		return math.Inf(-1)
	}
	var max = values[0]
	for i := 1; i < len(values); i++ {
		max = math.Max(max, values[i])
	}
	return max
}

// Returns the minimum value among all parameters.
func Min(values ...float64) float64 {
	if len(values) == 0 {
		return math.Inf(+1)
	}
	var min = values[0]
	for i := 1; i < len(values); i++ {
		min = math.Min(min, values[i])
	}
	return min
}
