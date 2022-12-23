package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputByteGrid("2022/23")
	elves := make(map[point]struct{})
	grid.Iter(func(r, c int) {
		if grid[r][c] == '#' {
			elves[point{c, r}] = struct{}{}
		}
	})
	full := func(p point) bool { return lib.MapHasKey(elves, p) }

	order := []lib.Dir{lib.North, lib.South, lib.West, lib.East} // starting order to consider moves

	const rounds = 10
	for round := 0; round < rounds; round++ {
		// "During the first half of each round, each Elf considers the eight positions adjacent to
		// themself. If no other Elves are in one of those eight positions, the Elf does not do
		// anything during this round. Otherwise, the Elf looks in each of four directions in the
		// following order and proposes moving one step in the first valid direction:
		// - If there is no Elf in the N, NE, or NW adjacent positions, the Elf proposes moving north one step.
		// - If there is no Elf in the S, SE, or SW adjacent positions, the Elf proposes moving south one step.
		// - If there is no Elf in the W, NW, or SW adjacent positions, the Elf proposes moving west one step.
		// - If there is no Elf in the E, NE, or SE adjacent positions, the Elf proposes moving east one step."
		props := make(map[point][]point) // dests to elves proposing moving there
		for elf := range elves {
			dirs := map[lib.Dir]bool{
				lib.North: full(elf.off(-1, -1)) || full(elf.off(0, -1)) || full(elf.off(1, -1)),
				lib.South: full(elf.off(-1, 1)) || full(elf.off(0, 1)) || full(elf.off(1, 1)),
				lib.West:  full(elf.off(-1, -1)) || full(elf.off(-1, 0)) || full(elf.off(-1, 1)),
				lib.East:  full(elf.off(1, -1)) || full(elf.off(1, 0)) || full(elf.off(1, 1)),
			}
			if !dirs[lib.North] && !dirs[lib.South] && !dirs[lib.West] && !dirs[lib.East] {
				continue // no neighbors, so do nothing
			}
			for i := 0; i < len(order); i++ {
				if dir := order[(round+i)%len(order)]; !dirs[dir] {
					dst := elf.off(dir.DC(), dir.DR())
					props[dst] = append(props[dst], elf)
					break
				}
			}
		}

		// "After each Elf has had a chance to propose a move, the second half of the round can
		// begin. Simultaneously, each Elf moves to their proposed destination tile if they were the
		// only Elf to propose moving to that position. If two or more Elves propose moving to the
		// same position, none of those Elves move."
		for dst, es := range props {
			if len(es) == 1 {
				delete(elves, es[0])
				elves[dst] = struct{}{}
			}
		}
	}

	// "Simulate the Elves' process and find the smallest rectangle that contains the Elves after 10
	// rounds. How many empty ground tiles does that rectangle contain?"
	xmin, xmax := math.MaxInt, math.MinInt
	ymin, ymax := math.MaxInt, math.MinInt
	for elf := range elves {
		xmin = lib.Min(xmin, elf.x)
		xmax = lib.Max(xmax, elf.x)
		ymin = lib.Min(ymin, elf.y)
		ymax = lib.Max(ymax, elf.y)
	}
	fmt.Println((xmax-xmin+1)*(ymax-ymin+1) - len(elves))
}

type point struct{ x, y int }

func (p point) off(dx, dy int) point { return point{p.x + dx, p.y + dy} }

func dump(elves map[point]struct{}) string {
	xmin, xmax := math.MaxInt, math.MinInt
	ymin, ymax := math.MaxInt, math.MinInt
	for elf := range elves {
		xmin = lib.Min(xmin, elf.x)
		xmax = lib.Max(xmax, elf.x)
		ymin = lib.Min(ymin, elf.y)
		ymax = lib.Max(ymax, elf.y)
	}
	grid := lib.NewByteGrid(ymax-ymin+1, xmax-xmin+1, '.')
	for elf := range elves {
		grid[elf.y-ymin][elf.x-xmin] = '#'
	}
	return grid.Dump()
}
