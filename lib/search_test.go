package lib

import (
	"reflect"
	"strings"
	"testing"
)

func TestAStar(t *testing.T) {
	grid := ByteLines(strings.TrimLeft(`
..#.....
..#.#...
..#.###.
....#...
#####...
........
...#####
........
`, "\n"), '.', '#')

	const (
		nrows = 8
		ncols = 8
		want  = 30
	)

	start := PackInts(0, 0)
	end := PackInts(nrows-1, ncols-1)
	got := AStar([]uint64{start},
		func(s uint64) bool { return s == end },
		func(s uint64) []uint64 {
			r, c := UnpackInt2(s)
			var next []uint64
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r0, c0 := r+off[0], c+off[1]
				if r0 >= 0 && r0 < nrows && c0 >= 0 && c0 < ncols && grid[r0][c0] == '.' {
					next = append(next, PackInts(r0, c0))
				}
			}
			return next
		},
		func(s uint64) int {
			r, c := UnpackInt2(s)
			return Abs(nrows-1-r) + Abs(ncols-1-c) // Manhattan distance
		})
	if got != want {
		t.Errorf("AStar() returned %v; want %v", got, want)
	}
}

func TestBFS(t *testing.T) {
	got, _ := BFS(PackInts(0, 0), func(s uint64) []uint64 {
		x, y := UnpackInt2(s)
		var next []uint64
		for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			x0, y0 := x+off[0], y+off[1]
			if x0 >= 0 && x0 < 3 && y0 >= 0 && y0 < 3 {
				next = append(next, PackInts(x0, y0))
			}
		}
		return next
	}, nil)
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
