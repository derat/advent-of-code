package lib

import "fmt"

// AssertEq panics if a doesn't equal b.
func AssertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("%v != %v", a, b))
	}
}

// AssertLess panics if a >= b.
func AssertLess(a, b int) {
	if a >= b {
		panic(fmt.Sprintf("%v >= %v", a, b))
	}
}

// AssertLessEq panics if a > b.
func AssertLessEq(a, b int) {
	if a > b {
		panic(fmt.Sprintf("%v > %v", a, b))
	}
}
