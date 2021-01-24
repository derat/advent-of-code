package main

import (
	"fmt"
	"math"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const (
	animateSearch = false
	animateFill   = false
	delay         = 50 * time.Millisecond
)

func main() {
	vm := lib.NewIntcode(lib.InputInt64s("2019/15"))

	var r, c int                                           // droid location
	var mr, mc int                                         // location we tried to move to
	var or, oc int                                         // oxygen location
	states := map[uint64]status{lib.PackInts(r, c): empty} // known locations
	dests := make(map[uint64]struct{})                     // unknown locations

	// Returns neighboring squares to p with non-wall or unknown states.
	neighbors := func(p uint64) []uint64 {
		r, c := lib.UnpackIntSigned2(p)
		ns := make([]uint64, 0, 4)
		for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			r, c := r+off[0], c+off[1]
			p := lib.PackInts(r, c)
			if st, ok := states[p]; !ok || st != wall {
				ns = append(ns, p)
			}
		}
		return ns
	}

	// Initially explore neighbors of the starting square.
	for _, p := range neighbors(lib.PackInts(r, c)) {
		dests[p] = struct{}{}
	}

	// Silly code for visualizing the search.
	var drows int // last drawn rows
	drawSearch := func() {
		if !animateSearch {
			return
		}
		if drows > 0 {
			fmt.Printf("\033[%dA", drows) // clear previous grid
		}
		grid, _, _ := dump(states, r, c, dests)
		drows = len(grid)
		fmt.Println(lib.DumpBytes(grid))
		time.Sleep(delay)
	}
	drawSearch()

	vm.InFunc = func() int64 {
		// Perform BFS from r, c to find the nearest square from dests.
		start := lib.PackInts(r, c)
		_, from := lib.BFS(start, neighbors, &lib.BFSOptions{AnyDests: dests})
		for s := range from {
			if _, ok := dests[s]; !ok {
				continue
			}
			// Walk backward to find the square after our current one.
			for ; ; s = from[s] {
				if f := from[s]; f == start {
					mr, mc = lib.UnpackIntSigned2(s)
					switch {
					case mr == r-1:
						return 1 // north
					case mr == r+1:
						return 2 // south
					case mc == c-1:
						return 3 // west
					case mc == c+1:
						return 4 // east
					}
				}
			}
		}
		panic("No dest")
	}

	vm.OutFunc = func(v int64) {
		mp := lib.PackInts(mr, mc)
		delete(dests, mp)
		st := status(v)
		states[mp] = st

		switch st {
		case empty, oxygen:
			// Move was successful. Add unknown neighboring locations to dests.
			r, c = mr, mc
			for _, p := range neighbors(lib.PackInts(r, c)) {
				if _, ok := states[p]; !ok {
					dests[p] = struct{}{}
				}
			}
			if st == oxygen {
				or, oc = r, c
			}
		case wall:
		default:
			lib.Panicf("Invalid output %v", v)
		}

		drawSearch()

		// Exit if we've fully explored the maze.
		if len(dests) == 0 {
			vm.Halt()
		}
	}

	// Actually explore the maze (needed for both parts).
	lib.Assert(vm.Run())

	// Part 1: Print length of shortest path from starting location to oxygen.
	os := lib.PackInts(or, oc)
	fmt.Println(lib.AStar([]uint64{lib.PackInts(0, 0)},
		func(s uint64) bool { return s == os },
		neighbors,
		func(s uint64) int {
			r, c := lib.UnpackIntSigned2(s)
			return lib.Abs(or-r) + lib.Abs(oc-c)
		}))

	// Part 2: Print number of steps for oxygen to spread to all open locations.
	steps, _ := lib.BFS(os, neighbors, nil)
	maxSteps := lib.Max(lib.MapIntVals(steps)...)
	fmt.Println(maxSteps)

	// Silly code for animating oxygen filling the maze.
	if animateFill {
		grid, rmin, cmin := dump(states, 0, 0, nil)
		fmt.Println(lib.DumpBytes(grid))

		for i := 1; i <= maxSteps; i++ {
			for r, row := range grid {
				for c := range row {
					if v, ok := steps[lib.PackInts(r+rmin, c+cmin)]; ok && v == i {
						row[c] = 'O'
					}
				}
			}
			fmt.Printf("\033[%dA", len(grid)) // clear previous grid
			fmt.Println(lib.DumpBytes(grid))
			time.Sleep(delay)
		}
	}
}

type status int

const (
	wall   status = 0
	empty         = 1
	oxygen        = 2
)

// dump dumps the supplied tile states, droid position, and destination locations
// to a printable grid. It also returns the min row and column number from states
// and dests (so additional data can be added to the supplied grid).
func dump(states map[uint64]status, dr, dc int, dests map[uint64]struct{}) ([][]byte, int, int) {
	rmin, rmax := math.MaxInt32, math.MinInt32
	cmin, cmax := math.MaxInt32, math.MinInt32
	for p := range states {
		r, c := lib.UnpackIntSigned2(p)
		rmin, rmax = lib.Min(rmin, r), lib.Max(rmax, r)
		cmin, cmax = lib.Min(cmin, c), lib.Max(cmax, c)
	}
	for p := range dests {
		r, c := lib.UnpackIntSigned2(p)
		rmin, rmax = lib.Min(rmin, r), lib.Max(rmax, r)
		cmin, cmax = lib.Min(cmin, c), lib.Max(cmax, c)
	}

	grid := lib.NewBytes(rmax-rmin+1, cmax-cmin+1, ' ')
	for p, s := range states {
		r, c := lib.UnpackIntSigned2(p)
		var ch byte
		switch {
		case r == dr && c == dc:
			ch = 'D'
		case s == wall:
			ch = '#'
		case s == empty:
			ch = '.'
		case s == oxygen:
			ch = '@'
		}
		grid[r-rmin][c-cmin] = ch
	}
	for p := range dests {
		r, c := lib.UnpackIntSigned2(p)
		grid[r-rmin][c-cmin] = '?'
	}

	return grid, rmin, cmin
}
