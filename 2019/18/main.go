package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputLinesBytes("2019/18")

	var sr, sc int              // starting position
	klocs := make([][2]int, 26) // index is id, vals are r,c
	for r, row := range grid {
		for c, ch := range row {
			switch {
			case ch == '@': // entrance
				sr, sc = r, c
			case ch >= 'a' && ch <= 'z': // key
				klocs[ch-'a'] = [2]int{r, c}
			}
		}
	}

	allKeys := (1 << 26) - 1 // bitfield representing all keys

	// Part 1: Print minimum number of steps to pick up all keys.
	steps := lib.AStar(
		[]state{{sr, sc, 0}},
		func(s state) bool { return s.keys&allKeys == allKeys },
		func(s state, m map[state]int) {
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				n := s
				n.r += off[0]
				n.c += off[1]
				ch := grid[n.r][n.c]
				switch {
				case ch == '#': // wall
					continue
				case ch >= 'a' && ch <= 'z': // key
					n.keys |= 1 << (ch - 'a')
				case ch >= 'A' && ch <= 'Z': // door
					if s.keys&(1<<(ch-'A')) == 0 {
						continue // don't have key
					}
				case ch == '.': // empty space
				case ch == '@': // entrance
				default:
					lib.Panicf("Invalid char %q at %d, %d", ch, n.r, n.c)
				}
				m[n] = 1
			}
		},
		func(s state) int {
			// Use the Manhattan distance to the farthest key as a lower bound.
			var max int
			for id, loc := range klocs {
				if s.keys&id == 0 { // don't have the key yet
					max = lib.Max(lib.Abs(loc[0]-s.r)+lib.Abs(loc[1]-s.c), max)
				}
			}
			return max
		})
	fmt.Println(steps)

	// Part 2: Split the map into four and print the minimum number of steps to pick
	// up all the keys using four robots.
	grid[sr][sc] = '#'
	grid[sr-1][sc] = '#'
	grid[sr+1][sc] = '#'
	grid[sr][sc-1] = '#'
	grid[sr][sc+1] = '#'

	// Robot starting locations.
	rlocs := [4][2]int{
		{sr - 1, sc - 1},
		{sr - 1, sc + 1},
		{sr + 1, sc - 1},
		{sr + 1, sc + 1},
	}

	// My input has the property that the top-right quadrant consisted only of keys
	// while the bottom-right quadrant has no immediately-accessible keys. All three
	// of the non-right quadrants have dependencies on each other, though.
	//
	// I first just adapted my A* code to support interface{} states (so I could encode
	// four r,c pairs and 26 bits of keys as a []byte and then cast it to a string).
	// This was way too slow, though.
	//
	// Next, I tried to take advantage of the properties of my input. It's clearly
	// optimal to grab all the keys in the top-right quadrant first, since none of them
	// are blocked by doors. I was at a loss about what to do next, though. I would
	// need to switch off between the robots in the other three quadrants, but I didn't
	// know how to determine the order in which they should pick up keys when there were
	// multiple possibilities. Performing a search in parallel using all three robots
	// was still too slow.
	//
	// After taking a walk and thinking about it some more, I realized that it's never
	// advantageous to move a robot unless it goes all the way to pick up another key.
	// So, I could just perform a BFS from each key to find the shortest paths to all
	// other reachable keys, keeping track of the keys and doors passed on the way.
	// Then I could use that data in the original A* search, just tracking the current
	// key (or starting) location for each robot and the keys that are held. I struggled
	// for a while with a bug in my next-state function where shorter paths were
	// overwritten with longer ones, leading to inconsistent and incorrect results.

	dests := make(map[int][]dest, 30) // keys are key (0-25) or starting loc (26-29)
	for i, loc := range klocs {
		dests[i] = explore(grid, loc[0], loc[1])
	}
	for i, loc := range rlocs {
		dests[26+i] = explore(grid, loc[0], loc[1])
	}

	steps = lib.AStar(
		[]state2{{[4]int{26, 27, 28, 29}, 0}},
		func(s state2) bool { return s.keys == allKeys },
		func(s state2, m map[state2]int) {
			for ri, loc := range s.locs {
				for _, d := range dests[loc] {
					if d.doors&^s.keys == 0 { // have all required keys
						n := s
						n.locs[ri] = d.loc
						n.keys |= d.keys
						// This is very important and led to a lot of pain and swearing:
						// we'll end up with multiple paths to the same destination that
						// are identical in terms of the keys that they acquire but
						// different in terms of number of steps. Use the lowest-cost one.
						if os, ok := m[n]; !ok || d.steps < os {
							m[n] = d.steps
						}
					}
				}
			}
		},
		func(s state2) int {
			// As a crappy lower bound, take each robot's current location and add
			// the minimum distance required for it to get another key. A better one
			// would be to find the minimum path to get all needed and reachable keys,
			// but I'm not sure how to get that data without performing a bunch more
			// searches.
			var sum int
			for _, loc := range s.locs {
				min := math.MaxInt32
				for _, d := range dests[loc] {
					if s.keys|d.loc != s.keys {
						min = lib.Min(min, d.steps)
					}
				}
				if min != math.MaxInt32 {
					sum += min
				}
			}
			return sum
		})
	fmt.Println(steps)
}

type state struct {
	r, c, keys int
}

type state2 struct {
	locs [4]int // index into dests in main()
	keys int    // bitfield; position 0 is 'a'
}

// explore performs a BFS on grid starting at r,c and returns all
// minimum unique (in terms of doors and keys passed) paths to keys.
func explore(grid [][]byte, r, c int) []dest {
	type state struct{ r, c, keys, doors int }
	states, _ := lib.BFS(
		[]state{{r, c, 0, 0}},
		func(s state, m map[state]struct{}) {
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				n := s
				n.r += off[0]
				n.c += off[1]
				if ch := grid[n.r][n.c]; ch == '#' {
					continue
				} else if ch >= 'a' && ch <= 'z' {
					n.keys |= 1 << (ch - 'a')
				} else if ch >= 'A' && ch <= 'Z' {
					n.doors |= 1 << (ch - 'A')
				}
				m[n] = struct{}{}
			}
		}, nil)

	// Only return states ending on other keys.
	var dests []dest
	for s, steps := range states {
		if s.r == r && s.c == c {
			continue // skip paths leading back to us
		}
		ch := grid[s.r][s.c]
		if ch >= 'a' && ch <= 'z' {
			dests = append(dests, dest{int(ch - 'a'), steps, s.keys, s.doors})
		}
	}
	return dests
}

// dest describes a path from a point to a key.
type dest struct {
	loc         int // dest key in [0,25]
	steps       int // number of steps to get to dest
	keys, doors int // keys and doors passed (including dest)
}

// letters is a debug function that prints the letters represented by a bitfield.
func letters(vals int, upper bool) string {
	base := 'a'
	if upper {
		base = 'A'
	}
	var str string
	for i := 0; i < 26; i++ {
		if vals&(1<<i) != 0 {
			str += string(rune(i) + base)
		}
	}
	return str
}
