// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestAssertEq(t *testing.T) {
	for _, tc := range []struct {
		a, b any
		eq   bool
	}{
		{true, true, true},
		{true, false, false},
		{3, 3, true},
		{3, 2, false},
		{'a', 'a', true},
		{int64(3), int64(3), true},
		{int(3), int64(3), true},
		{byte('a'), 'a', true},
		{byte('a'), rune('a'), true},
		{3.5, 3.5, true},
		{3.5, 3.4, false},
		{float32(3.5), float64(3.5), true},
		{"abc", "abc", true},
		{"abc", "def", false},
		{[]int{}, []int{}, true},
		{[]int{2, 3}, []int{2, 3}, true},
		{[]int{2, 3}, []int{3, 2}, false},
		{[]int{}, []int(nil), false},
		{map[int]int{}, map[int]int{}, true},
		{map[int]int{1: 2}, map[int]int{1: 2}, true},
		{map[int]int{1: 2}, map[int]int{2: 1}, false},
		{map[int]int{}, map[int]int(nil), false},
	} {
		t.Run(fmt.Sprintf("%v (%v), %v (%v)", tc.a, reflect.TypeOf(tc.a), tc.b, reflect.TypeOf(tc.b)),
			func(t *testing.T) {
				defer func() {
					if r := recover(); r != nil && tc.eq {
						t.Errorf("AssertEq panicked: %s", r)
					} else if r == nil && !tc.eq {
						t.Error("AssertEq unexpectedly passed")
					}
				}()
				AssertEq(tc.a, tc.b)
			})
	}
}
