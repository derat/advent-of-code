package lib

// Min returns the minimum of the supplied values.
func Min(vals ...int) int {
	Assertf(len(vals) > 0, "No values given")
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
	Assertf(len(vals) > 0, "No values given")
	max := vals[0]
	for _, v := range vals[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Clamp clamps val within [min, max].
func Clamp(val, min, max int) int {
	return Min(Max(val, min), max)
}

// Sum returns the sum of the supplied values.
func Sum(vals ...int) int {
	var sum int
	for _, v := range vals {
		sum += v
	}
	return sum
}

// Product returns the product of the supplied values.
func Product(vals ...int) int {
	Assertf(len(vals) > 0, "No values given")
	prod := vals[1]
	for _, v := range vals[1:] {
		prod *= v
	}
	return prod
}

// Abs returns the absolute value of v.
func Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
