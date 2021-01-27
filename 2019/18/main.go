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
	steps := lib.AStar([]uint64{pack(sr, sc, 0)},
		func(s uint64) bool {
			_, _, keys := unpack(s)
			return keys&allKeys == allKeys
		},
		func(s uint64) []uint64 {
			var next []uint64
			r, c, keys := unpack(s)
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				nr, nc := r+off[0], c+off[1]
				nkeys := keys
				ch := grid[nr][nc]
				switch {
				case ch == '#': // wall
					continue
				case ch >= 'a' && ch <= 'z': // key
					nkeys |= 1 << (ch - 'a')
				case ch >= 'A' && ch <= 'Z': // door
					if keys&(1<<(ch-'A')) == 0 {
						continue // don't have key
					}
				case ch == '.': // empty space
				case ch == '@': // entrance
				default:
					lib.Panicf("Invalid char %q at %d, %d", ch, nr, nc)
				}
				next = append(next, pack(nr, nc, nkeys))
			}
			return next
		},
		func(s uint64) int {
			// Use the Manhattan distance to the farthest key as a lower bound.
			r, c, keys := unpack(s)
			var max int
			for id, loc := range klocs {
				if keys&id == 0 { // don't have the key yet
					max = lib.Max(lib.Abs(loc[0]-r)+lib.Abs(loc[1]-c), max)
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

	steps = lib.AStarVarCost(
		[]uint64{pack4([4]int{26, 27, 28, 29}, 0)},
		func(s uint64) bool {
			_, keys := unpack4(s)
			return keys == allKeys
		},
		func(s uint64) map[uint64]int {
			next := make(map[uint64]int)
			locs, keys := unpack4(s)
			for ri, loc := range locs {
				for _, d := range dests[loc] {
					if d.doors&^keys == 0 { // have all required keys
						nlocs := locs
						nlocs[ri] = d.loc
						nkeys := keys | d.keys
						p := pack4(nlocs, nkeys)
						// This is very important and led to a lot of pain and swearing:
						// we'll end up with multiple paths to the same destination that
						// are identical in terms of the keys that they acquire but
						// different in terms of number of steps. Use the lowest-cost one.
						if os, ok := next[p]; !ok || d.steps < os {
							next[p] = d.steps
						}
					}
				}
			}
			return next
		},
		func(s uint64) int {
			// As a crappy lower bound, take each robot's current location and add
			// the minimum distance required for it to get another key. A better one
			// would be to find the minimum path to get all needed and reachable keys,
			// but I'm not sure how to get that data without performing a bunch more
			// searches.
			var sum int
			locs, keys := unpack4(s)
			for _, loc := range locs {
				min := math.MaxInt32
				for _, d := range dests[loc] {
					if keys|d.loc != keys {
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

func pack(r, c, keys int) uint64 {
	p := lib.PackInt(0, r, 16, 0)
	p = lib.PackInt(p, c, 16, 16)
	p = lib.PackInt(p, keys, 32, 32)
	return p
}

func unpack(p uint64) (r, c, keys int) {
	r = lib.UnpackInt(p, 16, 0)
	c = lib.UnpackInt(p, 16, 16)
	keys = lib.UnpackInt(p, 32, 32)
	return r, c, keys
}

func pack4(locs [4]int, keys int) uint64 {
	var p uint64
	p = lib.PackInt(p, locs[0], 8, 0)
	p = lib.PackInt(p, locs[1], 8, 8)
	p = lib.PackInt(p, locs[2], 8, 16)
	p = lib.PackInt(p, locs[3], 8, 24)
	p = lib.PackInt(p, keys, 32, 32)
	return p
}

func unpack4(p uint64) (locs [4]int, keys int) {
	locs[0] = lib.UnpackInt(p, 8, 0)
	locs[1] = lib.UnpackInt(p, 8, 8)
	locs[2] = lib.UnpackInt(p, 8, 16)
	locs[3] = lib.UnpackInt(p, 8, 24)
	keys = lib.UnpackInt(p, 32, 32)
	return
}

// explore performs a BFS on grid starting at r,c and returns all
// minimum unique (in terms of doors and keys passed) paths to keys.
func explore(grid [][]byte, r, c int) []dest {
	ids := make(map[byte]int) // key/door (e.g. 'A' or 'z') to ID
	getID := func(ch byte) int {
		id, ok := ids[ch]
		if !ok {
			id = len(ids)
			ids[ch] = id
			lib.AssertLessEq(len(ids), 48) // limit in pack()
		}
		return id
	}

	pack := func(r, c, feats int) uint64 {
		var p uint64
		p = lib.PackInt(p, r, 8, 0)
		p = lib.PackInt(p, c, 8, 8)
		p = lib.PackInt(p, feats, 48, 16)
		return p
	}

	unpack := func(p uint64) (r, c, feats int) {
		r = lib.UnpackInt(p, 8, 0)
		c = lib.UnpackInt(p, 8, 8)
		feats = lib.UnpackInt(p, 48, 16)
		return
	}

	states, _ := lib.BFS(pack(r, c, 0),
		func(s uint64) []uint64 {
			r, c, feats := unpack(s)
			var next []uint64
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				nr, nc := r+off[0], c+off[1]
				nfeats := feats
				if ch := grid[nr][nc]; ch == '#' {
					continue
				} else if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
					nfeats |= 1 << getID(ch)
				}
				next = append(next, pack(nr, nc, nfeats))
			}
			return next
		}, nil)

	// Only return states ending on other keys.
	var dests []dest
	for s, steps := range states {
		r0, c0, feats := unpack(s)
		if r0 == r && c0 == c {
			continue // skip paths leading back to us
		}
		ch := grid[r0][c0]
		if ch < 'a' || ch > 'z' {
			continue // skip non-keys
		}
		var keys, doors int
		for ch, id := range ids {
			if feats&(1<<id) != 0 {
				switch {
				case ch >= 'a' && ch <= 'z':
					keys |= 1 << int(ch-'a')
				case ch >= 'A' && ch <= 'Z':
					doors |= 1 << int(ch-'A')
				default:
					lib.Panicf("Invalid char %v", ch)
				}
			}
		}
		dests = append(dests, dest{int(ch - 'a'), steps, keys, doors})
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
