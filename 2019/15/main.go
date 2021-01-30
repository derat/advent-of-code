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

	var dp [2]int                          // droid location
	var mp [2]int                          // location we tried to move to
	var op [2]int                          // oxygen location
	states := map[[2]int]status{dp: empty} // known locations
	dests := make(map[[2]int]struct{})     // unknown locations

	// Returns neighboring squares to p with non-wall or unknown states.
	neighbors := func(p [2]int) [][2]int {
		ns := make([][2]int, 0, 4)
		for _, o := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			n := [2]int{p[0] + o[0], p[1] + o[1]}
			if st, ok := states[n]; !ok || st != wall {
				ns = append(ns, n)
			}
		}
		return ns
	}

	// Initially explore neighbors of the starting square.
	for _, p := range neighbors(dp) {
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
		grid, _, _ := dump(states, dp, dests)
		drows = len(grid)
		fmt.Println(lib.DumpBytes(grid))
		time.Sleep(delay)
	}
	drawSearch()

	vm.InFunc = func() int64 {
		// Perform BFS from r, c to find the nearest square from dests.
		_, from := lib.BFS([]interface{}{dp}, func(si interface{}, m map[interface{}]struct{}) {
			for _, n := range neighbors(si.([2]int)) {
				m[n] = struct{}{}
			}
		}, &lib.BFSOptions{AnyDests: lib.GenericSet(dests)})
		for si := range from {
			if _, ok := dests[si.([2]int)]; !ok {
				continue
			}
			// Walk backward to find the square after our current one.
			for ; ; si = from[si] {
				if f := from[si].([2]int); f == dp {
					mp = si.([2]int)
					switch {
					case mp[0] == dp[0]-1:
						return 1 // north
					case mp[0] == dp[0]+1:
						return 2 // south
					case mp[1] == dp[1]-1:
						return 3 // west
					case mp[1] == dp[1]+1:
						return 4 // east
					}
				}
			}
		}
		panic("No dest")
	}

	vm.OutFunc = func(v int64) {
		delete(dests, mp)
		st := status(v)
		states[mp] = st

		switch st {
		case empty, oxygen:
			// Move was successful. Add unknown neighboring locations to dests.
			dp = mp
			for _, p := range neighbors(mp) {
				if _, ok := states[p]; !ok {
					dests[p] = struct{}{}
				}
			}
			if st == oxygen {
				op = mp
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
	fmt.Println(lib.AStar([]interface{}{[2]int{0, 0}},
		func(si interface{}) bool { return si.([2]int) == op },
		func(si interface{}, m map[interface{}]int) {
			for _, n := range neighbors(si.([2]int)) {
				m[n] = 1
			}
		},
		func(si interface{}) int {
			s := si.([2]int)
			return lib.Abs(op[0]-s[0]) + lib.Abs(op[1]-s[1])
		}))

	// Part 2: Print number of steps for oxygen to spread to all open locations.
	steps, _ := lib.BFS([]interface{}{op}, func(si interface{}, m map[interface{}]struct{}) {
		for _, n := range neighbors(si.([2]int)) {
			m[n] = struct{}{}
		}
	}, nil)
	maxSteps := lib.Max(lib.MapIntVals(steps)...)
	fmt.Println(maxSteps)

	// Silly code for animating oxygen filling the maze.
	if animateFill {
		grid, rmin, cmin := dump(states, [2]int{0, 0}, nil)
		fmt.Println(lib.DumpBytes(grid))

		for i := 1; i <= maxSteps; i++ {
			for r, row := range grid {
				for c := range row {
					if v, ok := steps[[2]int{r + rmin, c + cmin}]; ok && v == i {
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
func dump(states map[[2]int]status, dp [2]int, dests map[[2]int]struct{}) ([][]byte, int, int) {
	rmin, rmax := math.MaxInt32, math.MinInt32
	cmin, cmax := math.MaxInt32, math.MinInt32
	for p := range states {
		rmin, rmax = lib.Min(rmin, p[0]), lib.Max(rmax, p[0])
		cmin, cmax = lib.Min(cmin, p[1]), lib.Max(cmax, p[1])
	}
	for p := range dests {
		rmin, rmax = lib.Min(rmin, p[0]), lib.Max(rmax, p[0])
		cmin, cmax = lib.Min(cmin, p[1]), lib.Max(cmax, p[1])
	}

	grid := lib.NewBytes(rmax-rmin+1, cmax-cmin+1, ' ')
	for p, s := range states {
		var ch byte
		switch {
		case p == dp:
			ch = 'D'
		case s == wall:
			ch = '#'
		case s == empty:
			ch = '.'
		case s == oxygen:
			ch = '@'
		}
		grid[p[0]-rmin][p[1]-cmin] = ch
	}
	for p := range dests {
		grid[p[0]-rmin][p[1]-cmin] = '?'
	}

	return grid, rmin, cmin
}
