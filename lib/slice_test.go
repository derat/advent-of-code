// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import (
	"reflect"
	"testing"
)

func TestRotateSlice(t *testing.T) {
	for _, tc := range []struct {
		in   string
		amt  int
		want string
	}{
		{"abcd", 0, "abcd"},
		{"abcd", 1, "dabc"},
		{"abcd", 2, "cdab"},
		{"abcd", 3, "bcda"},
		{"abcd", 4, "abcd"},
		{"abcd", 5, "dabc"},
		{"abcd", -1, "bcda"},
		{"abcd", -2, "cdab"},
		{"abcd", -3, "dabc"},
		{"abcd", -4, "abcd"},
		{"abcd", -5, "bcda"},
		{"", 0, ""},
		{"", 1, ""},
		{"", -1, ""},
		{"a", 0, "a"},
		{"a", 1, "a"},
		{"a", -1, "a"},
	} {
		b := []byte(tc.in)
		if RotateSlice(b, tc.amt); string(b) != tc.want {
			t.Errorf("RotateBy(%q, %d) = %q; want %q", tc.in, tc.amt, b, tc.want)
		}
	}
}

func TestMove(t *testing.T) {
	for _, tc := range []struct {
		in        string
		s1, s2, d int
		want      string
	}{
		// Move beginning.
		{"abcdef", 0, 0, 0, "abcdef"},
		{"abcdef", 0, 1, 0, "abcdef"},
		{"abcdef", 0, 2, 1, "cabdef"},
		{"abcdef", 0, 2, 2, "cdabef"},
		{"abcdef", 0, 2, 3, "cdeabf"},
		{"abcdef", 0, 2, 4, "cdefab"},
		{"abcdef", 0, 5, 1, "fabcde"},

		// Move end.
		{"abcdef", 6, 6, 0, "abcdef"},
		{"abcdef", 5, 5, 0, "abcdef"},
		{"abcdef", 5, 6, 5, "abcdef"},
		{"abcdef", 5, 6, 4, "abcdfe"},
		{"abcdef", 4, 6, 3, "abcefd"},
		{"abcdef", 1, 6, 0, "bcdefa"},

		// Move middle.
		{"abcdef", 2, 2, 0, "abcdef"},
		{"abcdef", 2, 4, 2, "abcdef"},
		{"abcdef", 2, 4, 1, "acdbef"},
		{"abcdef", 2, 4, 3, "abecdf"},
		{"abcdef", 2, 4, 0, "cdabef"},
		{"abcdef", 2, 4, 4, "abefcd"},
	} {
		b := []byte(tc.in)
		if Move(b, tc.s1, tc.s2, tc.d); string(b) != tc.want {
			t.Errorf("Move(%q, %d, %d, %d) = %q; want %q", tc.in, tc.s1, tc.s2, tc.d, b, tc.want)
		}
	}
}

func TestReverse(t *testing.T) {
	for _, tc := range []struct {
		orig, want []int
	}{
		{[]int{}, []int{}},
		{[]int{1}, []int{1}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	} {
		got := append([]int{}, tc.orig...)
		Reverse(got)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("Reverse(%v) = %v; want %v", tc.orig, got, tc.want)
		}
	}
}
