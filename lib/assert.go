package lib

import "fmt"

// Assert panics if v is false.
func Assert(v bool) {
	if !v {
		panic("Assertion failed")
	}
}

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
