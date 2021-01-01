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

// Max returns the maximum of the supplied values.
func Max(vals ...int) int {
	if len(vals) == 0 {
		panic("Zero values")
	}
	max := vals[0]
	for _, v := range vals[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Sum returns the sum of the supplied values.
func Sum(vals ...int) int {
	var sum int
	for _, v := range vals {
		sum += v
	}
	return sum
}

// Abs returns the absolute value of v.
func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
