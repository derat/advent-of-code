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

	type state [2]int
	end := state{nrows - 1, ncols - 1}
	got := AStar([]state{{0, 0}},
		func(s state) bool { return s == end },
		func(s state, m map[state]int) {
			for _, off := range []state{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r, c := s[0]+off[0], s[1]+off[1]
				if r >= 0 && r < nrows && c >= 0 && c < ncols && grid[r][c] == '.' {
					m[state{r, c}] = 1
				}
			}
		},
		func(s state) int {
			return Abs(end[0]-s[0]) + Abs(end[1]-s[1]) // Manhattan distance
		})
	if got != want {
		t.Errorf("AStar() returned %v; want %v", got, want)
	}
}

func TestBFS(t *testing.T) {
	type state [2]int
	got, _ := BFS([]state{{0, 0}},
		func(s state, m map[state]struct{}) {
			for _, off := range []state{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				x, y := s[0]+off[0], s[1]+off[1]
				if x >= 0 && x < 3 && y >= 0 && y < 3 {
					m[state{x, y}] = struct{}{}
				}
			}
		}, nil)
	if want := map[state]int{
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
	expSteps := map[string]int{
		"start":       0,
		"start,a":     1,
		"start,b":     1,
		"start,a,end": 2,
		"start,b,end": 2,
	}
	expFrom := map[string]string{
		"start":       "start",
		"start,a":     "start",
		"start,b":     "start",
		"start,a,end": "start,a",
		"start,b,end": "start,b",
	}

	for _, tc := range []struct {
		opts  *BFSOptions[string]
		steps map[string]int
		from  map[string]string
	}{
		{nil, expSteps, expFrom},
		{&BFSOptions[string]{}, expSteps, expFrom},
		{&BFSOptions[string]{NoSteps: true}, nil, expFrom},
		{&BFSOptions[string]{NoFrom: true}, expSteps, nil},
		{&BFSOptions[string]{NoSteps: true, NoFrom: true}, nil, nil},
	} {
		npaths := 0
		steps, from := BFS([]string{"start"}, func(s string, m map[string]struct{}) {
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
