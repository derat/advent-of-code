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
	locs := make(map[uint64]int) // packed (r, c) to location ID
	for r, row := range rows {
		for c, ch := range row {
			if ch >= '0' && ch <= '9' {
				id := int(ch - '0')
				if id == 0 {
					loc0 = [2]int{r, c}
				}
				p := lib.PackInts(r, c)
				_, seen := locs[p]
				lib.Assertf(!seen, "Duplicate ID %v", id)
				locs[p] = id
			}
		}
	}

	for _, part := range []int{1, 2} {
		steps := lib.AStar([]uint64{pack(loc0[0], loc0[1], 1)},
			func(s uint64) bool {
				r, c, vis := unpack(s)
				// In part 2, we need to end up back at location 0.
				return vis == 1<<len(locs)-1 && (part != 2 || (r == loc0[0] && c == loc0[1]))
			},
			func(s uint64) []uint64 {
				r, c, vis := unpack(s)
				var ns []uint64
				for _, n := range [][2]int{{r - 1, c}, {r + 1, c}, {r, c - 1}, {r, c + 1}} {
					r0, c0, v0 := n[0], n[1], vis
					// Skip moves that go out-of-bounds or hit a wall.
					if r0 < 0 || r0 >= len(rows) || c0 < 0 || c0 >= len(rows[r]) || rows[r0][c0] == '#' {
						continue
					}
					// Check if we've reached a location.
					if id, ok := locs[lib.PackInts(r0, c0)]; ok {
						v0 |= 1 << id
					}
					ns = append(ns, pack(r0, c0, v0))
				}
				return ns
			},
			func(s uint64) int {
				r, c, vis := unpack(s)
				if vis == 1<<len(locs)-1 {
					if part == 1 {
						return 0
					}
					return r + c // need to get back to start
				}
				// Use the distance to the unvisited location farthest from the robot's
				// location as a lower bound for the number of additional steps needed.
				// I initially tried to use the distance to the nearest location plus the
				// distance from there to the farthest location, but that resulted in a
				// too-high answer (presumably it overestimated the minimum distance for
				// some reason -- maybe I had a bug?).
				var fd int
				for p, id := range locs {
					if vis&(1<<id) == 0 {
						lr, lc := lib.UnpackInt2(p)
						ld := lib.Abs(lr-r) + lib.Abs(lc-c)
						fd = lib.Max(fd, ld)
					}
				}
				if part == 1 {
					// TODO: There's some bug here that I don't understand.
					// Now that AStar() randomizes the order of "next" steps
					// internally, I sometimes get a too-high answer for part 1.
					return fd
				}
				return lib.Min(fd, r+c)
			})
		fmt.Println(steps)
	}
}

func pack(r, c, vis int) uint64 {
	return lib.PackInts(r, c, vis)
}

func unpack(s uint64) (r, c, vis int) {
	v := lib.UnpackInts(s, 3)
	return v[0], v[1], v[2]
}
