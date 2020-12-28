package lib

// Min returns the minimum of the supplied values.
func Min(vals ...int) int {
	if len(vals) == 0 {
		panic("Zero values")
	}
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}
