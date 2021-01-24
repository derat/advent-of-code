package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/15")

	vm := lib.NewIntcode(input)

	var r, c int                                           // current location
	var mr, mc int                                         // location we tried to move to
	var or, oc int                                         // oxygen location
	states := map[uint64]status{lib.PackInts(r, c): empty} // known locations

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

	dests := make(map[uint64]struct{}) // unknown locations
	for _, p := range neighbors(lib.PackInts(r, c)) {
		dests[p] = struct{}{}
	}

	vm.InFunc = func() int64 {
		// Perform BFS from r, c. As soon as we find a square in dests, start
		// moving toward it.
		start := lib.PackInts(r, c)
		queue := []uint64{start}
		from := map[uint64]uint64{start: start}
		var dest uint64 // chosen location from dests
	Search:
		for len(queue) > 0 {
			p := queue[0]
			queue = queue[1:]
			for _, n := range neighbors(p) {
				if _, ok := from[n]; ok {
					continue // already been there
				}
				if st, ok := states[n]; ok && st == wall {
					continue // can't move there
				}
				from[n] = p
				if _, ok := dests[n]; ok {
					dest = n
					break Search
				}
				queue = append(queue, n)
			}
		}
		// Walk backward to find the square after our current one.
		for {
			if f := from[dest]; f == start {
				mr, mc = lib.UnpackIntSigned2(dest)
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
			} else {
				dest = from[dest]
			}
		}
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

		// Exit if we've fully explored the maze.
		if len(dests) == 0 {
			vm.Halt()
		}
	}

	lib.Assert(vm.Run())

	/*
		rmin, rmax := math.MaxInt32, math.MinInt32
		cmin, cmax := math.MaxInt32, math.MinInt32
		for p := range states {
			r, c := lib.UnpackIntSigned2(p)
			rmin, rmax = lib.Min(rmin, r), lib.Max(rmax, r)
			cmin, cmax = lib.Min(cmin, c), lib.Max(cmax, c)
		}
		grid := lib.NewBytes(rmax-rmin+1, cmax-cmin+1, ' ')
		for p, s := range states {
			r, c := lib.UnpackIntSigned2(p)
			var ch byte
			switch {
			case r == 0 && c == 0:
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
		fmt.Println(lib.DumpBytes(grid))
	*/

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
	steps := lib.BFS(os, neighbors, nil, -1)
	fmt.Println(lib.Max(lib.MapIntVals(steps)...))
}

type status int

const (
	wall   status = 0
	empty         = 1
	oxygen        = 2
)
