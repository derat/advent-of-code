package lib

import (
	"reflect"
	"strings"
	"testing"
)

func TestAStar(t *testing.T) {
	grid := NewByteGridString(strings.TrimLeft(`
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

func TestBFS_NoSteps_NoFrom(t *testing.T) {
	expSteps := map[interface{}]int{
		"start":       0,
		"start,a":     1,
		"start,b":     1,
		"start,a,end": 2,
		"start,b,end": 2,
	}
	expFrom := map[interface{}]interface{}{
		"start":       "start",
		"start,a":     "start",
		"start,b":     "start",
		"start,a,end": "start,a",
		"start,b,end": "start,b",
	}

	for _, tc := range []struct {
		opts  *BFSOptions
		steps map[interface{}]int
		from  map[interface{}]interface{}
	}{
		{nil, expSteps, expFrom},
		{&BFSOptions{}, expSteps, expFrom},
		{&BFSOptions{NoSteps: true}, nil, expFrom},
		{&BFSOptions{NoFrom: true}, expSteps, nil},
		{&BFSOptions{NoSteps: true, NoFrom: true}, nil, nil},
	} {
		npaths := 0
		steps, from := BFS([]interface{}{"start"}, func(si interface{}, m map[interface{}]struct{}) {
			s := si.(string)
			if strings.HasSuffix(s, ",end") {
				npaths++
				return
			}
			switch s {
			case "start":
				m["start,a"] = struct{}{}
				m["start,b"] = struct{}{}
			case "start,a":
				m["start,a,end"] = struct{}{}
			case "start,b":
				m["start,b,end"] = struct{}{}
			default:
				t.Fatalf("Unexpected state %q", s)
			}
		}, tc.opts)
		if npaths != 2 {
			t.Errorf("BFS(%+v) found %v path(s); want 2", tc.opts, npaths)
		}
		if !reflect.DeepEqual(steps, tc.steps) {
			t.Errorf("BFS(%+v) returned steps %v; want %v", tc.opts, steps, tc.steps)
		}
		if !reflect.DeepEqual(from, tc.from) {
			t.Errorf("BFS(%+v) returned from %v; want %v", tc.opts, from, tc.from)
		}
	}
}
