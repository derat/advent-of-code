package lib

import (
	"reflect"
	"testing"
)

func TestBFS(t *testing.T) {
	got := BFS(PackInts(0, 0), func(s uint64) []uint64 {
		x, y := UnpackInt2(s)
		var next []uint64
		for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			x0, y0 := x+off[0], y+off[1]
			if x0 >= 0 && x0 < 3 && y0 >= 0 && y0 < 3 {
				next = append(next, PackInts(x0, y0))
			}
		}
		return next
	})
	if want := map[uint64]int{
		PackInts(0, 0): 0,
		PackInts(1, 0): 1,
		PackInts(2, 0): 2,
		PackInts(0, 1): 1,
		PackInts(1, 1): 2,
		PackInts(2, 1): 3,
		PackInts(0, 2): 2,
		PackInts(1, 2): 3,
		PackInts(2, 2): 4,
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("BFS() returned %v; want %v", got, want)
	}
}
