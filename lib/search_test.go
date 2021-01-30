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

	end := [2]int{nrows - 1, ncols - 1}
	got := AStar([]interface{}{[2]int{0, 0}},
		func(si interface{}) bool { return si.([2]int) == end },
		func(si interface{}, m map[interface{}]int) {
			s := si.([2]int)
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r, c := s[0]+off[0], s[1]+off[1]
				if r >= 0 && r < nrows && c >= 0 && c < ncols && grid[r][c] == '.' {
					m[[2]int{r, c}] = 1
				}
			}
		},
		func(si interface{}) int {
			s := si.([2]int)
			return Abs(end[0]-s[0]) + Abs(end[1]-s[1]) // Manhattan distance
		})
	if got != want {
		t.Errorf("AStar() returned %v; want %v", got, want)
	}
}

func TestBFS(t *testing.T) {
	got, _ := BFS([]interface{}{[2]int{0, 0}},
		func(si interface{}, m map[interface{}]struct{}) {
			s := si.([2]int)
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				x, y := s[0]+off[0], s[1]+off[1]
				if x >= 0 && x < 3 && y >= 0 && y < 3 {
					m[[2]int{x, y}] = struct{}{}
				}
			}
		}, nil)
	if want := map[interface{}]int{
		[2]int{0, 0}: 0,
		[2]int{1, 0}: 1,
		[2]int{2, 0}: 2,
		[2]int{0, 1}: 1,
		[2]int{1, 1}: 2,
		[2]int{2, 1}: 3,
		[2]int{0, 2}: 2,
		[2]int{1, 2}: 3,
		[2]int{2, 2}: 4,
	}; !reflect.DeepEqual(got, want) {
		t.Errorf("BFS() returned %v; want %v", got, want)
	}
}
