package lib

import (
	"reflect"
	"sort"
	"testing"
)

func TestFindCombos(t *testing.T) {
	for _, tc := range []struct {
		items   []int
		initial uint64
		target  int
		want    []uint64
	}{
		{[]int{2, 5, 4, 1, 1}, 0b11111, 7, []uint64{0b11, 0b01101, 0b10101, 0b11010}},
		{[]int{2, 5, 4, 1, 1}, 0b01111, 7, []uint64{0b11, 0b01101}},
		{[]int{2, 5, 4, 1, 1}, 0b00000, 7, []uint64{}},
		{[]int{2, 5, 4, 1, 1}, 0b11111, 13, []uint64{0b11111}},
		{[]int{2, 5, 4, 1, 1}, 0b11111, 14, []uint64{}},
	} {
		got := FindCombos(tc.items, tc.initial, tc.target)
		sort.Slice(tc.want, func(i, j int) bool { return tc.want[i] < tc.want[j] })
		sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("FindCombos(%v, %#b, %v) = %#b; want %#b", tc.items, tc.initial, tc.target, got, tc.want)
		}
	}
}
