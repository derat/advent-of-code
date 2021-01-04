package lib

import "reflect"

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
