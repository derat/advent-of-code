package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	rows := lib.InputLinesBytes("2016/24", '#', '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	// The state needs to include the robot's current row and column and the locations that we've
	// reached at least once.
	var loc0 [2]int
	locs := make(map[[2]int]int) // r,c to location ID
	for r, row := range rows {
		for c, ch := range row {
			if ch >= '0' && ch <= '9' {
				id := int(ch - '0')
				if id == 0 {
					loc0 = [2]int{r, c}
				}
				p := [2]int{r, c}
				_, seen := locs[p]
				lib.Assertf(!seen, "Duplicate ID %v", id)
				locs[p] = id
			}
		}
	}

	for _, part := range []int{1, 2} {
		steps := lib.AStar([]interface{}{state{loc0[0], loc0[1], 1}},
			func(si interface{}) bool {
				s := si.(state)
				// In part 2, we need to end up back at location 0.
				return s.vis == 1<<len(locs)-1 && (part != 2 || (s.r == loc0[0] && s.c == loc0[1]))
			},
			func(si interface{}, m map[interface{}]int) {
				s := si.(state)
				for _, p := range [][2]int{{s.r - 1, s.c}, {s.r + 1, s.c}, {s.r, s.c - 1}, {s.r, s.c + 1}} {
					n := state{p[0], p[1], s.vis}
					// Skip moves that go out-of-bounds or hit a wall.
					if n.r < 0 || n.r >= len(rows) || n.c < 0 || n.c >= len(rows[n.r]) || rows[n.r][n.c] == '#' {
						continue
					}
					// Check if we've reached a location.
					if id, ok := locs[[2]int{n.r, n.c}]; ok {
						n.vis |= 1 << id
					}
					m[n] = 1
				}
			},
			func(si interface{}) int {
				s := si.(state)
				if s.vis == 1<<len(locs)-1 {
					if part == 1 {
						return 0
					}
					return s.r + s.c // need to get back to start
				}
				// Use the distance to the unvisited location farthest from the robot's
				// location as a lower bound for the number of additional steps needed.
				// I initially tried to use the distance to the nearest location plus the
				// distance from there to the farthest location, but that resulted in a
				// too-high answer (presumably it overestimated the minimum distance for
				// some reason -- maybe I had a bug?).
				var fd int
				for p, id := range locs {
					if s.vis&(1<<id) == 0 {
						ld := lib.Abs(p[0]-s.r) + lib.Abs(p[1]-s.c)
						fd = lib.Max(fd, ld)
					}
				}
				if part == 1 {
					// TODO: There's some bug here that I don't understand.
					// Now that AStar() randomizes the order of "next" steps
					// internally, I sometimes get a too-high answer for part 1.
					return fd
				}
				return lib.Min(fd, s.r+s.c)
			})
		fmt.Println(steps)
	}
}

type state struct{ r, c, vis int }
