package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/derat/advent-of-code/lib"
)

const (
	animate = false
	delay   = 100 * time.Millisecond
)

func main() {
	input := lib.InputLinesBytes("2018/15", '#', '.', 'G', 'E')

	// Part 1: Print number of *full* rounds times summed HP of living units.
	ca := newCave(input, 200, 3, 3)
	if animate {
		ca.dump(false)
	}
	rounds := 0
	for {
		if !ca.round() {
			break
		}
		rounds++
		if animate {
			ca.dump(true)
			time.Sleep(delay)
		}
	}
	if animate {
		ca.dump(true)
	}
	fmt.Println(score(rounds, ca))

	// Part 2: Print score using minimum elf AP to not suffer any casualties.
	// This is slow but I am so done with this problem.
Loop:
	for ap := 4; true; ap++ {
		ca := newCave(input, 200, ap, 3)
		elves := ca.left[elf]
		rounds := 0
		for {
			done := !ca.round()
			if ca.left[elf] < elves {
				continue Loop // elves died; continue
			}
			if done {
				fmt.Println(score(rounds, ca))
				break Loop
			}
			rounds++
		}
	}
}

func score(rounds int, ca *cave) int {
	var sum int
	for _, u := range ca.units {
		if u.alive() {
			sum += u.hp
		}
	}
	return rounds * sum
}

type cave struct {
	grid  [][]byte         // '.' is empty, '#' is wall
	units []*unit          // both living and dead
	locs  map[uint64]*unit // living units keyed by packed R,C
	left  map[species]int  // living units per species
}

func newCave(grid [][]byte, hp, eap, gap int) *cave {
	ca := &cave{
		grid: lib.CopyBytes(grid),
		locs: make(map[uint64]*unit),
		left: make(map[species]int),
	}

	for r, row := range ca.grid {
		for c, ch := range row {
			if ch == 'E' || ch == 'G' {
				spec := elf
				ap := eap
				if ch == 'G' {
					spec = goblin
					ap = gap
				}
				u := &unit{spec, r, c, hp, ap}
				ca.units = append(ca.units, u)
				ca.locs[lib.PackInts(r, c)] = u
				ca.left[spec]++
				row[c] = '.'
			}
		}
	}

	return ca
}

// done returns true if combat is done because a species has been wiped out.
func (ca *cave) done() bool {
	return ca.left[elf] == 0 || ca.left[goblin] == 0
}

// round performs a round of combat, with each living unit taking a turn.
// Returns true if a full round is completed and false if the round ended
// early due to one side being wiped out.
func (ca *cave) round() bool {
	// Iterate over units in reading order.
	sort.Slice(ca.units, func(i, j int) bool { return ca.units[i].cmp(ca.units[j], false) })
	for _, u := range ca.units {
		if u.alive() {
			if !ca.act(u) {
				return false
			}
		}
	}
	return true
}

// act performs a turn on behalf of u.
// Returns true if the turn was completed normally and false if it
// ended early due to there not being any enemies left.
func (ca *cave) act(u *unit) bool {
	lib.Assertf(u.alive(), "Dead unit can't take turn")

	// Identify open squares within range of living enemy units.
	dests := make(map[uint64][]*unit)
	var enemies int
	for _, o := range ca.units {
		if o.alive() && o.spec != u.spec {
			enemies++
			for _, k := range ca.neighbors(o.r, o.c, u) {
				dests[k] = append(dests[k], o)
			}
		}
	}
	// "If no targets remain, combat ends."
	if enemies == 0 {
		return false
	}
	// "If the unit is not already in range of a target, and there are no open
	// squares which are in range of a target, the unit ends its turn."
	if len(dests) == 0 {
		return true
	}

	// If we're not already in range of the enemy, move.
	if !lib.MapHasKey(dests, lib.PackInts(u.r, u.c)) {
		ca.move(u, lib.MapUint64Keys(dests))
	}

	// If the unit isn't in range after moving, it ends its turn.
	targets, ok := dests[lib.PackInts(u.r, u.c)]
	if !ok {
		return true
	}

	// "Otherwise, the adjacent target with the fewest hit points is selected; in a tie, the
	// adjacent target with the fewest hit points which is first in reading order is selected."
	sort.Slice(targets, func(i, j int) bool { return targets[i].cmp(targets[j], true) })
	t := targets[0]
	t.hp -= u.ap
	if !t.alive() {
		delete(ca.locs, lib.PackInts(t.r, t.c))
		ca.left[t.spec]--
	}

	return true
}

// move moves u to the preferred packed location from dests.
func (ca *cave) move(u *unit, dests []uint64) {
	// Find the neighboring square and dest with the shortest path between them.
	var dest, next uint64
	min := math.MaxInt32
	neighbors := ca.neighbors(u.r, u.c, nil)
	for _, d := range dests {
		// Find minimum number of steps from dest to each neighboring square.
		costs, _ := lib.BFS(d, func(s uint64) []uint64 {
			r, c := lib.UnpackInt2(s)
			return ca.neighbors(r, c, u)
		}, &lib.BFSOptions{AllDests: neighbors, MaxSteps: min})
		for _, n := range neighbors {
			cost, ok := costs[n]
			if !ok {
				continue
			}
			// Favor shortest path, then better dest, then better neighbor.
			if cost < min ||
				(cost == min && (cmp(d, dest) || (d == dest && cmp(n, next)))) {
				min = cost
				dest = d
				next = n
			}
		}
	}

	// If none of the dests were reachable, we can't move.
	if min == math.MaxInt32 {
		return
	}

	delete(ca.locs, lib.PackInts(u.r, u.c))
	u.r, u.c = lib.UnpackInt2(next)
	ca.locs[lib.PackInts(u.r, u.c)] = u
}

// open returns true if the specified position is in bounds and doesn't
// contain a wall or a unit.
func (ca *cave) open(r, c int) bool {
	if r < 0 || r >= len(ca.grid) || c < 0 || c >= len(ca.grid[r]) {
		return false
	}
	return ca.grid[r][c] == '.' && !lib.MapHasKey(ca.locs, lib.PackInts(r, c))
}

// neighbors returns open neighboring squares to r, c.
// If u is non-nil, neighboring squares are also returned if they are occupied by u.
func (ca *cave) neighbors(r, c int, u *unit) []uint64 {
	var next []uint64
	for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		r0, c0 := r+off[0], c+off[1]
		if ca.open(r0, c0) || (u != nil && u.r == r0 && u.c == c0) {
			next = append(next, lib.PackInts(r0, c0))
		}
	}
	return next
}

// dump writes the current state to stdout.
// If replace is true, overwrites the previously-dumped state.
func (ca *cave) dump(replace bool) {
	if replace {
		fmt.Printf("\033[%dA", len(ca.grid)) // move cursor up
	}
	for r, row := range ca.grid {
		var us []string
		for c, ch := range row {
			p := lib.PackInts(r, c)
			if u, ok := ca.locs[p]; ok {
				if u.spec == elf {
					fmt.Print("E")
					us = append(us, fmt.Sprintf("E(%d)", u.hp))
				} else {
					fmt.Print("G")
					us = append(us, fmt.Sprintf("G(%d)", u.hp))
				}
			} else {
				fmt.Print(string(rune(ch)))
			}
		}
		if len(us) > 0 {
			fmt.Print("   ", strings.Join(us, ", "))
		}
		if replace {
			fmt.Print("\033[K") // clear from cursor right
		}
		fmt.Println()
	}
}

type species int

const (
	elf species = iota
	goblin
)

type unit struct {
	spec   species
	r, c   int
	hp, ap int
}

func (u *unit) alive() bool {
	return u.hp > 0
}

// cmp returns true if u precedes o.
// If hp is true, the unit with the lower HP is favored.
// Otherwise, or if the HP are the same, compares coordinates in reading order.
func (u *unit) cmp(o *unit, hp bool) bool {
	if hp {
		if u.hp < o.hp {
			return true
		} else if u.hp > o.hp {
			return false
		}
	}
	if u.r < o.r {
		return true
	} else if u.r > o.r {
		return false
	}
	return u.c < o.c
}

// dist returns the Manhattan distance between two packed R,C coordinates.
func dist(a, b uint64) int {
	ar, ac := lib.UnpackInt2(a)
	br, bc := lib.UnpackInt2(b)
	return lib.Abs(ar-br) + lib.Abs(ac-bc)
}

// cmp returns true if packed R,C coordinate a precedes b in reading order.
func cmp(a, b uint64) bool {
	ar, ac := lib.UnpackInt2(a)
	br, bc := lib.UnpackInt2(b)
	if ar < br {
		return true
	} else if ar > br {
		return false
	}
	return ac < bc
}
