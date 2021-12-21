package lib

import (
	"fmt"
	"reflect"
)

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
	var eq bool
	if av, ok := convertInt(a); ok {
		if bv, ok := convertInt(b); ok {
			eq = av == bv
		} else {
			panic(fmt.Sprintf("%v not convertible to int64", reflect.TypeOf(b)))
		}
	} else if av, ok := convertFloat(a); ok {
		if bv, ok := convertFloat(b); ok {
			eq = av == bv
		} else {
			panic(fmt.Sprintf("%v not convertible to float64", reflect.TypeOf(b)))
		}
	} else if k := reflect.TypeOf(a).Kind(); k == reflect.Slice || k == reflect.Map {
		eq = reflect.DeepEqual(a, b)
	} else {
		eq = a == b
	}

	if !eq {
		panic(fmt.Sprintf("%v != %v", a, b))
	}
}

func convertInt(i interface{}) (int64, bool) {
	switch v := i.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	default:
		return 0, false
	}
}

func convertFloat(i interface{}) (float64, bool) {
	switch v := i.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	default:
		return 0, false
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
