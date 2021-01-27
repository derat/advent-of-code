package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputLinesBytes("2019/20")

	labels := make(map[string][]uint64) // label e.g. "AZ" to packed r,c
	for r := 0; r < len(grid)-1; r++ {
		for c := 0; c < len(grid[0])-1; c++ {
			if ch := grid[r][c]; ch >= 'A' && ch <= 'Z' {
				if chr := grid[r][c+1]; chr >= 'A' && chr <= 'Z' { // horizontal label
					la := string([]byte{ch, chr})
					if c > 0 && grid[r][c-1] == '.' { // point left of label
						labels[la] = append(labels[la], lib.PackInts(r, c-1))
					} else { // point right of label
						lib.Assertf(grid[r][c+2] == '.', "Expected dot at %d, %d for %q", r, c+2, la)
						labels[la] = append(labels[la], lib.PackInts(r, c+2))
					}
				} else if chd := grid[r+1][c]; chd >= 'A' && chd <= 'Z' { // vertical label
					la := string([]byte{ch, chd})
					if r > 0 && grid[r-1][c] == '.' { // point above label
						labels[la] = append(labels[la], lib.PackInts(r-1, c))
					} else { // point below label
						lib.Assertf(grid[r+2][c] == '.', "Expected dot at %d, %d for %q", r+2, c, la)
						labels[la] = append(labels[la], lib.PackInts(r+2, c))
					}
				}
			}
		}
	}

	var start, end uint64                                 // packed r,c
	portals := make(map[uint64]uint64, 2*(len(labels)-2)) // packed r,c to r,c
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
	steps, _ := lib.BFS(start, func(s uint64) []uint64 {
		r, c := lib.UnpackInt2(s)
		ps, hp := portals[s]
		var next []uint64
		for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			r0, c0 := r+off[0], c+off[1]
			if r0 < 0 || c0 < 0 || r0 >= len(grid) || c0 >= len(grid[r]) {
				continue // probably overkill since there are labels around the edges
			}
			if ch := grid[r0][c0]; ch == '.' { // move to empty space
				next = append(next, lib.PackInts(r0, c0))
			} else if hp && ch >= 'A' && ch <= 'Z' { // use portal
				next = append(next, ps)
			}
		}
		return next
	}, &lib.BFSOptions{AnyDests: map[uint64]struct{}{end: struct{}{}}})
	fmt.Println(steps[end])

	// Part 2: Inner labels go to more-deeply nested versions of the maze; outer
	// labels go to less-deeply. AA and ZZ do nothing in nested mazes. Travel from
	// outermost AA to outermost ZZ.
	sr, sc := lib.UnpackInt2(start)
	er, ec := lib.UnpackInt2(end)
	es := pack(er, ec, 0)
	fmt.Println(lib.AStar(
		[]uint64{pack(sr, sc, 0)},
		func(s uint64) bool { return s == es },
		func(s uint64) []uint64 {
			r, c, depth := unpack(s)
			ps, hp := portals[lib.PackInts(r, c)]
			out := r == 2 || c == 2 || r == len(grid)-3 || c == len(grid[0])-3

			var next []uint64
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r0, c0 := r+off[0], c+off[1]
				if r0 < 0 || c0 < 0 || r0 >= len(grid) || c0 >= len(grid[r]) {
					continue // probably overkill since there are labels around the edges
				}
				if ch := grid[r0][c0]; ch == '.' { // move to empty space
					next = append(next, pack(r0, c0, depth))
				} else if hp && ch >= 'A' && ch <= 'Z' { // use portal
					pr, pc := lib.UnpackInt2(ps)
					if out && depth > 0 { // can only take outside portals when nested
						next = append(next, pack(pr, pc, depth-1))
					} else if !out { // take inside portal to deeper nesting
						next = append(next, pack(pr, pc, depth+1))
					}
				}
			}
			return next
		},
		func(s uint64) int {
			// Crappy heuristic: just use the difference in depth from outermost.
			_, _, depth := unpack(s)
			return lib.Abs(depth)
		}))
}

func pack(r, c, depth int) uint64 {
	return lib.PackInts(r, c, depth)
}

func unpack(p uint64) (r, c, depth int) {
	vals := lib.UnpackInts(p, 3)
	return vals[0], vals[1], vals[2]
}
