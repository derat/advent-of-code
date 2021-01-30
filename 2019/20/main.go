package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputLinesBytes("2019/20")

	labels := make(map[string][][2]int) // label e.g. "AZ" to r,c
	for r := 0; r < len(grid)-1; r++ {
		for c := 0; c < len(grid[0])-1; c++ {
			if ch := grid[r][c]; ch >= 'A' && ch <= 'Z' {
				if chr := grid[r][c+1]; chr >= 'A' && chr <= 'Z' { // horizontal label
					la := string([]byte{ch, chr})
					if c > 0 && grid[r][c-1] == '.' { // point left of label
						labels[la] = append(labels[la], [2]int{r, c - 1})
					} else { // point right of label
						lib.Assertf(grid[r][c+2] == '.', "Expected dot at %d, %d for %q", r, c+2, la)
						labels[la] = append(labels[la], [2]int{r, c + 2})
					}
				} else if chd := grid[r+1][c]; chd >= 'A' && chd <= 'Z' { // vertical label
					la := string([]byte{ch, chd})
					if r > 0 && grid[r-1][c] == '.' { // point above label
						labels[la] = append(labels[la], [2]int{r - 1, c})
					} else { // point below label
						lib.Assertf(grid[r+2][c] == '.', "Expected dot at %d, %d for %q", r+2, c, la)
						labels[la] = append(labels[la], [2]int{r + 2, c})
					}
				}
			}
		}
	}

	var start, end [2]int                                 // r,c
	portals := make(map[[2]int][2]int, 2*(len(labels)-2)) // r,c to r,c
	for la, ps := range labels {
		switch la {
		case "AA":
			lib.Assertf(len(ps) == 1, "Want 1 point for %q; got %v", la, ps)
			start = ps[0]
		case "ZZ":
			lib.Assertf(len(ps) == 1, "Want 1 point for %q; got %v", la, ps)
			end = ps[0]
		default:
			lib.Assertf(len(ps) == 2, "Want 2 points for %q; got %v", la, ps)
			portals[ps[0]] = ps[1]
			portals[ps[1]] = ps[0]
		}
	}

	// Part 1: Minimum number of steps to go from AA to ZZ.
	// I'm just using BFS instead of A* here since the number of states is small
	// and it'd be tricky to write a proper heuristic function.
	steps, _ := lib.BFS([]interface{}{start}, func(si interface{}, m map[interface{}]struct{}) {
		s := si.([2]int)
		ps, hp := portals[s]
		for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			r, c := s[0]+off[0], s[1]+off[1]
			if r < 0 || c < 0 || r >= len(grid) || c >= len(grid[r]) {
				continue // probably overkill since there are labels around the edges
			}
			if ch := grid[r][c]; ch == '.' { // move to empty space
				m[[2]int{r, c}] = struct{}{}
			} else if hp && ch >= 'A' && ch <= 'Z' { // use portal
				m[ps] = struct{}{}
			}
		}
	}, &lib.BFSOptions{AnyDests: map[interface{}]struct{}{end: struct{}{}}})
	fmt.Println(steps[end])

	// Part 2: Inner labels go to more-deeply nested versions of the maze; outer
	// labels go to less-deeply. AA and ZZ do nothing in nested mazes. Travel from
	// outermost AA to outermost ZZ.
	es := state{end[0], end[1], 0}
	fmt.Println(lib.AStar(
		[]interface{}{state{start[0], start[1], 0}},
		func(si interface{}) bool { return si.(state) == es },
		func(si interface{}, m map[interface{}]int) {
			s := si.(state)
			ps, hp := portals[[2]int{s.r, s.c}]
			out := s.r == 2 || s.c == 2 || s.r == len(grid)-3 || s.c == len(grid[0])-3

			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r, c := s.r+off[0], s.c+off[1]
				if r < 0 || c < 0 || r >= len(grid) || c >= len(grid[r]) {
					continue // probably overkill since there are labels around the edges
				}
				if ch := grid[r][c]; ch == '.' { // move to empty space
					m[state{r, c, s.depth}] = 1
				} else if hp && ch >= 'A' && ch <= 'Z' { // use portal
					if out && s.depth > 0 { // can only take outside portals when nested
						m[state{ps[0], ps[1], s.depth - 1}] = 1
					} else if !out { // take inside portal to deeper nesting
						m[state{ps[0], ps[1], s.depth + 1}] = 1
					}
				}
			}
		},
		func(si interface{}) int {
			// Crappy heuristic: just use the difference in depth from outermost.
			return lib.Abs(si.(state).depth)
		}))
}

type state struct{ r, c, depth int }
