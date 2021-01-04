package lib

import "fmt"

// Assert panics if v is false.
func Assert(v bool) {
	if !v {
		panic("Assertion failed")
	}
}

// Assertf panics if v is false.
func Assertf(v bool, s string, args ...interface{}) {
	if !v {
		panic(fmt.Sprintf(s, args...))
	}
}

// Panicf panics with the supplied message.
func Panicf(s string, args ...interface{}) {
	Assertf(false, s, args...)
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
