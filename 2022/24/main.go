package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2022/24", '#', '.', '>', '<', '^', 'v')
	var enter, exit point
	initial := make(map[point]lib.Dir)
	grid.Iter(func(r, c int) {
		p := point{r, c}
		ch := grid[r][c]
		switch {
		case ch == '.' && r == 0:
			enter = p
		case ch == '.' && r == grid.MaxRow():
			exit = p
		case ch == '^':
			initial[p] = lib.Up
		case ch == 'v':
			initial[p] = lib.Down
		case ch == '>':
			initial[p] = lib.Right
		case ch == '<':
			initial[p] = lib.Left
		}
	})

	// Ignore the border.
	rmin, rmax := 1, grid.MaxRow()-1
	cmin, cmax := 1, grid.MaxCol()-1
	rows := rmax - rmin + 1
	cols := cmax - cmin + 1

	// Precompute the positions containing blizzards at each turn until the board loops.
	cycle := lib.LCM(rows, cols)
	blizzards := make([]map[point]struct{}, cycle)
	for i := 0; i < cycle; i++ {
		m := make(map[point]struct{})
		for p, dir := range initial {
			// I had a bug here where up- and left-moving blizzards would end up with negative
			// positions due to negative modulus being the worst thing ever.
			r := rmin + (p.r-rmin+(dir.DR()*(i%rows)))%rows
			c := cmin + (p.c-cmin+(dir.DC()*(i%cols)))%cols
			if r < rmin {
				r += rows
			}
			if c < cmin {
				c += cols
			}
			lib.AssertInRange(r, rmin, rmax)
			lib.AssertInRange(c, cmin, cmax)
			m[point{r, c}] = struct{}{}
		}
		blizzards[i] = m
	}

	// Part 1: "What is the fewest number of minutes required to avoid the blizzards and reach the goal?"
	fmt.Println(lib.AStar(
		[]state{{enter, 0}},
		func(s state) bool { return s.p == exit },
		func(s state, next map[state]int) {
			// I had a hard-to-find bug here where I was looking at the blizzards from the
			// current turn rather than from the next turn (which is what's relevant when
			// determining which next states are possible).
			nt := (s.turn + 1) % cycle
			bl := blizzards[nt]

			// We can move to the adjacent squares if they're in-bounds and won't be hit
			// by blizzards.
			for _, dir := range []lib.Dir{lib.Up, lib.Down, lib.Left, lib.Right} {
				np := point{s.p.r + dir.DR(), s.p.c + dir.DC()}
				valid := np == enter || np == exit ||
					(np.r >= rmin && np.r <= rmax && np.c >= cmin && np.c <= cmax)
				if valid && !lib.MapHasKey(bl, np) {
					next[state{np, nt}] = 1
				}
			}

			// We can also wait at the current position if a blizzard won't arrive there.
			if !lib.MapHasKey(bl, s.p) {
				next[state{s.p, nt}] = 1
			}
		},
		func(s state) int { return lib.Abs(exit.r-s.p.r) + lib.Abs(exit.c-s.p.c) },
	))
}

type point struct{ r, c int }

type state struct {
	p    point
	turn int // index into blizzards
}
