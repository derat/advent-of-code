package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2018/22")
	var depth, tr, tc int
	lib.Extract(lines[0], `^depth: (\d+)$`, &depth)
	lib.Extract(lines[1], `^target: (\d+),(\d+)$`, &tc, &tr)

	ero := func(geo int) int { return (geo + depth) % 20183 }

	geos := make(map[[2]int]int)
	var geo func(int, int) int
	geo = func(r, c int) int {
		k := [2]int{r, c}
		if g, ok := geos[k]; ok {
			return g
		}
		var g int
		switch {
		case r == 0 && c == 0, r == tr && c == tc:
			g = 0
		case r == 0:
			g = c * 16807
		case c == 0:
			g = r * 48271
		default:
			g = ero(geo(r, c-1)) * ero(geo(r-1, c))
		}
		geos[k] = g
		return g
	}

	// Part 1: Print total risk level for 0,0 to target coords.
	var total int
	for r := 0; r <= tr; r++ {
		for c := 0; c <= tc; c++ {
			total += ero(geo(r, c)) % 3
		}
	}
	fmt.Println(total)

	// Part 2: Minimum time to reach target.
	// Moving to another square takes 1 minute; switching tools takes 7 minutes.
	// - Rocky (ero % 3 == 0) needs climbing gear or torch.
	// - Wet (ero % 3 == 1) needs climbing gear or neither.
	// - Narrow (ero % 3 == 2) needs torch or neither.
	// Start at 0,0 with torch and get to tr,tc with climbing gear.
	const sw = 7 // time to switch tools
	end := state{tr, tc, torch}
	min := lib.AStar(
		[]interface{}{state{0, 0, torch}},
		func(si interface{}) bool { return si.(state) == end },
		func(si interface{}, m map[interface{}]int) {
			s := si.(state)

			// If we're holding a tool that can be equipped in an adjacent region,
			// we can move there in 1 minute.
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				n := state{s.r + off[0], s.c + off[1], s.t}
				if n.r < 0 || n.c < 0 {
					continue
				}
				if typ0 := ero(geo(n.r, n.c)) % 3; equip(n.t, typ0) {
					m[n] = 1
				}
			}

			// Also handle switching tools without moving.
			typ := ero(geo(s.r, s.c)) % 3
			for _, t := range []tool{climbing, torch, neither} {
				if t != s.t && equip(t, typ) {
					m[state{s.r, s.c, t}] = sw
				}
			}
		},
		func(si interface{}) int {
			// The lower bound to get to the target is the Manhattan distance to get
			// there plus the time needed to switch to the torch if it isn't equipped.
			s := si.(state)
			cost := lib.Abs(tr-s.r) + lib.Abs(tc-s.c)
			if s.t != torch {
				cost += sw
			}
			return cost
		})
	fmt.Println(min)
}

type tool int

const (
	climbing tool = iota
	torch
	neither
)

type state struct {
	r, c int
	t    tool
}

// equip returns true if tl can be used in a region of the supplied type.
func equip(tl tool, typ int) bool {
	switch typ {
	case 0: // rocky
		return tl == climbing || tl == torch
	case 1: // wet
		return tl == climbing || tl == neither
	case 2: // narrow
		return tl == torch || tl == neither
	}
	panic("Not reached")
}
