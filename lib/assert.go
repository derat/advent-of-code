package lib

import "fmt"

// Assert panics if a doesn't equal b.
func Assert(a, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("%v != %v", a, b))
	}
}
