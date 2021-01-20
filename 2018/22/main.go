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

	geos := make(map[uint64]int)
	var geo func(int, int) int
	geo = func(r, c int) int {
		k := lib.PackInts(r, c)
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
	end := pack(tr, tc, torch)
	min := lib.AStarVarCost([]uint64{pack(0, 0, torch)},
		func(s uint64) bool { return s == end },
		func(s uint64) map[uint64]int {
			r, c, t := unpack(s)
			m := make(map[uint64]int)

			// If we're holding a tool that can be equipped in an adjacent region,
			// we can move there in 1 minute.
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				r0, c0 := r+off[0], c+off[1]
				if r0 < 0 || c0 < 0 {
					continue
				}
				if typ0 := ero(geo(r0, c0)) % 3; equip(t, typ0) {
					m[pack(r0, c0, t)] = 1
				}
			}

			// Also handle switching tools without moving.
			typ := ero(geo(r, c)) % 3
			for _, t0 := range []tool{climbing, torch, neither} {
				if t0 != t && equip(t0, typ) {
					m[pack(r, c, t0)] = sw
				}
			}
			return m
		},
		func(s uint64) int {
			// The lower bound to get to the target is the Manhattan distance to get
			// there plus the time needed to switch to the torch if it isn't equipped.
			r, c, t := unpack(s)
			cost := lib.Abs(tr-r) + lib.Abs(tc-c)
			if t != torch {
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

func pack(r, c int, t tool) uint64 {
	return lib.PackInts(r, c, int(t))
}

func unpack(p uint64) (int, int, tool) {
	vals := lib.UnpackInts(p, 3)
	return vals[0], vals[1], tool(vals[2])
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
