// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import (
	"reflect"

	"golang.org/x/exp/constraints"
)

// ReverseBytes reverses b in-place. A pointer to b is also returned.
func ReverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// SliceIndexesWithVal returns indexes into s of elements equal to want.
func SliceIndexesWithVal[T comparable](s []T, want T) []int {
	var idxs []int
	for i, v := range s {
		if v == want {
			idxs = append(idxs, i)
		}
	}
	return idxs
}

// Rotate rotates the elements in [first,last) such that middle becomes the new first element.
// This comes from http://www.cplusplus.com/reference/algorithm/rotate/ .
func Rotate(first, middle, last int, swap func(i, j int)) {
	next := middle
	for first != next {
		swap(first, next)
		first, next = first+1, next+1
		if next == last {
			next = middle
		} else if first == middle {
			middle = next
		}
	}
}

// RotateBy is a wrapper around Rotate that rotates n elements right by amt.
func RotateBy(n, amt int, swap func(i, j int)) {
	if n == 0 || amt == 0 {
		return
	}
	middle := (2*n - (amt % n)) % n
	Rotate(0, middle, n, swap)
}

// RotateSlice is a wrapper around RotateBy that operates on a slice.
func RotateSlice[T any](s []T, amt int) {
	RotateBy(len(s), amt, reflect.Swapper(s))
}

// Move moves the elements in slice v's half-open range [s1,s2) to be at index d.
// Other elements are preserved and shifted as needed.
func Move[T any](v []T, s1, s2, d int) {
	Assertf(s2 >= s1, "Invalid range [%d,%d)", s1, s2)
	sn := s2 - s1 // number of elements being moved
	Assertf(d+sn <= len(v), "Copying [%d,%d) to %d overflows slice of length %v", s1, s2, d, len(v))

	// The elements that need be rotated are the ones in the half-open range
	// [min(s1,d), max(s2,d+sn)).
	start := Min(s1, d)
	end := Max(s2, d+sn)
	amt := d - s1
	RotateSlice(v[start:end], amt)
}

// Reverse reverses the order of the elements in the supplied slice.
func Reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// SetAscInt initializes slice s to ascending signed integer values.
func SetAscInt[T constraints.Integer](s []T, start T) {
	for i, v := 0, start; i < len(s); i, v = i+1, v+1 {
		s[i] = v
	}
}
