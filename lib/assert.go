package lib

import "fmt"

// AssertEq panics if a doesn't equal b.
func AssertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("%v != %v", a, b))
	}
}
