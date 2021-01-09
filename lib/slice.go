package lib

import (
	"reflect"
)

// ReverseBytes reverses b in-place. A pointer to b is also returned.
func ReverseBytes(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

// SliceIndexesWithVal returns indexes into s of elements equal to v.
func SliceIndexesWithVal(s, v interface{}) []int {
	sv := reflect.ValueOf(s)
	var idxs []int
	for i := 0; i < sv.Len(); i++ {
		if sv.Index(i).Interface() == v {
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

// RotateBy is a wrapper around Rotate that rotates n elements by amt.
func RotateBy(n, amt int, swap func(i, j int)) {
	if n == 0 || amt == 0 {
		return
	}
	middle := (2*n - (amt % n)) % n
	Rotate(0, middle, n, swap)
}

// RotateSlice is a wrapper around RotateBy that operates on a slice.
func RotateSlice(v interface{}, amt int) {
	RotateBy(reflect.ValueOf(v).Len(), amt, reflect.Swapper(v))
}

// Move moves the elements in slice v's half-open range [s1,s2) to be at index d.
// Other elements are preserved and shifted as needed.
func Move(v interface{}, s1, s2, d int) {
	Assertf(s2 >= s1, "Invalid range [%d,%d)", s1, s2)
	sn := s2 - s1 // number of elements being moved
	r := reflect.ValueOf(v)
	Assertf(d+sn <= r.Len(), "Copying [%d,%d) to %d overflows slice of length %v", s1, s2, d, r.Len())

	// The elements that need be rotated are the ones in the half-open range
	// [min(s1,d), max(s2,d+sn)).
	start := Min(s1, d)
	end := Max(s2, d+sn)
	amt := d - s1
	RotateSlice(r.Slice(start, end).Interface(), amt)
}

// Reverse reverses the order of the elements in the supplied slice.
func Reverse(s interface{}) {
	sv := reflect.ValueOf(s)
	swap := reflect.Swapper(s)
	for i, j := 0, sv.Len()-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// SetAscInt initializes slice s to ascending signed integer values.
func SetAscInt(s interface{}, start int64) {
	sv := reflect.ValueOf(s)
	for i, v := 0, start; i < sv.Len(); i, v = i+1, v+1 {
		sv.Index(i).SetInt(v)
	}
}

// SetAscUint initializes slice s to ascending unsigned integer values.
func SetAscUint(s interface{}, start uint64) {
	sv := reflect.ValueOf(s)
	for i, v := 0, start; i < sv.Len(); i, v = i+1, v+1 {
		sv.Index(i).SetUint(v)
	}
}
