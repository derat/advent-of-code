// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import (
	"fmt"
	"reflect"

	"golang.org/x/exp/constraints"
)

// Assert panics if v is false.
func Assert(v bool) {
	if !v {
		panic("Assertion failed")
	}
}

// Assertf panics if v is false.
func Assertf(v bool, s string, args ...any) {
	if !v {
		panic(fmt.Sprintf(s, args...))
	}
}

// AssertNil panics if v is non-nil.
func AssertNil(v any) {
	if v != nil {
		panic(v)
	}
}

// Panicf panics with the supplied message.
func Panicf(s string, args ...any) {
	Assertf(false, s, args...)
}

// AssertEq panics if a doesn't equal b.
func AssertEq(a, b any) {
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

func convertInt(i any) (int64, bool) {
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

func convertFloat(i any) (float64, bool) {
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
func AssertLess[T constraints.Ordered](a, b T) {
	if a >= b {
		panic(fmt.Sprintf("%v >= %v", a, b))
	}
}

// AssertLessEq panics if a > b.
func AssertLessEq[T constraints.Ordered](a, b T) {
	if a > b {
		panic(fmt.Sprintf("%v > %v", a, b))
	}
}

// AssertGreater panics if a <= b.
func AssertGreater[T constraints.Ordered](a, b T) {
	if a <= b {
		panic(fmt.Sprintf("%v <= %v", a, b))
	}
}

// AssertGreaterEq panics if a < b.
func AssertGreaterEq[T constraints.Ordered](a, b T) {
	if a < b {
		panic(fmt.Sprintf("%v < %v", a, b))
	}
}

// AssertInRange panics if v is not between min and max (inclusive).
func AssertInRange[T constraints.Ordered](v, min, max T) {
	AssertGreaterEq(v, min)
	AssertLessEq(v, max)
}
