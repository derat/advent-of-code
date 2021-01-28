package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLinesBytes("2019/24", '#', '.')

	// Part 1: Print biodiversity score after loop detected.
	grid := lib.CopyBytes(input)
	seen := map[uint64]struct{}{bio(grid): struct{}{}}
	for {
		next := lib.NewBytes(len(grid), len(grid[0]), '.')
		for r, row := range grid {
			for c, ch := range row {
				cnt := lib.CountBytes(grid, r-1, c, r-1, c, '#') +
					lib.CountBytes(grid, r+1, c, r+1, c, '#') +
					lib.CountBytes(grid, r, c-1, r, c-1, '#') +
					lib.CountBytes(grid, r, c+1, r, c+1, '#')
				switch ch {
				case '#':
					if cnt == 1 { // bug dies unless there's exactly one adjacent bug
						next[r][c] = '#'
					}
				case '.': // empty space infested if exactly 1 or 2 adjacent bugs
					if cnt == 1 || cnt == 2 {
						next[r][c] = '#'
					}
				default:
					lib.Panicf("Invalid byte %q", ch)
				}
			}
		}
		grid = next
		score := bio(grid)
		if _, ok := seen[score]; ok {
			fmt.Println(score)
			break
		} else {
			seen[score] = struct{}{}
		}
	}

	// Part 2: Use recursive grids and count bugs present after 200 minutes.

	// This is a hack. To make inner and outer checks easier, we always keep
	// one extra empty grid beyond the last non-empty one. I think that another
	// alternative would be to pass the inner and outer grids into update(),
	// but doing that seemed like it would be a bit fiddly.
	base := &rgrid{
		b:   lib.CopyBytes(input),
		in:  mkgrid(),
		out: mkgrid(),
	}
	base.in.out = base
	base.out.in = base

	for i := 0; i < 200; i++ {
		base = update(base, in|out)
	}
	fmt.Println(base.count(in | out))
}

// bio returns the biodiversity score for part 1.
// This also functions as a unique key for the grid's state.
func bio(grid [][]byte) uint64 {
	var score uint64
	pts := uint64(1)
	for _, row := range grid {
		for _, ch := range row {
			if ch == '#' {
				score += pts
			}
			pts <<= 1
		}
	}
	return score
}

type dir int

const (
	in dir = 1 << iota
	out
)

// update returns a new rgrid containing cur's next state.
// New grids are also created in the specified direction(s), if any.
func update(cur *rgrid, d dir) *rgrid {
	next := mkgrid()
	for r, row := range next.b {
		for c := range row {
			if r == 2 && c == 2 {
				continue
			}
			bug := cur.bug(r, c) == 1
			cnt := cur.neighbors(r, c)
			if (bug && cnt == 1) || (!bug && (cnt == 1 || cnt == 2)) {
				next.b[r][c] = '#'
			}
		}
	}

	if d&in != 0 && cur != nil {
		next.in = update(cur.in, in)
		next.in.out = next
		if next.in.count(in) == 0 {
			next.in.in = nil // prune extra empty grids
		}

	}
	if d&out != 0 && cur != nil {
		next.out = update(cur.out, out)
		next.out.in = next
		if next.out.count(out) == 0 {
			next.out.out = nil // prune extra empty grids
		}
	}
	return next
}

func mkgrid() *rgrid {
	return &rgrid{b: lib.NewBytes(5, 5, '.')}
}

// rgrid represents a recursive 5x5 grid for part 2.
type rgrid struct {
	b       [][]byte
	in, out *rgrid
}

// bug returns 1 if g contains a bug at the specified position and 0 otherwise.
func (g *rgrid) bug(r, c int) int {
	if g == nil {
		return 0
	}
	if r < 0 || c < 0 || r > 4 || c > 4 {
		return 0
	}
	if g.b[r][c] == '#' {
		return 1
	}
	return 0
}

// count returns the total number of bugs in g.
func (g *rgrid) count(d dir) int {
	if g == nil {
		return 0
	}
	cnt := lib.CountBytesFull(g.b, '#')
	if d&in != 0 {
		cnt += g.in.count(in)
	}
	if d&out != 0 {
		cnt += g.out.count(out)
	}
	return cnt
}

// neighbors returns the number of neighbors to r,c containing bugs,
// looking into inner and outer grids if necessary.
func (g *rgrid) neighbors(r, c int) int {
	if g == nil {
		return 0
	}

	cnt := g.bug(r-1, c) + g.bug(r+1, c) + g.bug(r, c-1) + g.bug(r, c+1)

	switch {
	case r == 1 && c == 2: // above center
		cnt += g.in.bug(0, 0) + g.in.bug(0, 1) + g.in.bug(0, 2) + g.in.bug(0, 3) + g.in.bug(0, 4)
	case r == 3 && c == 2: // below center
		cnt += g.in.bug(4, 0) + g.in.bug(4, 1) + g.in.bug(4, 2) + g.in.bug(4, 3) + g.in.bug(4, 4)
	case r == 2 && c == 1: // left of center
		cnt += g.in.bug(0, 0) + g.in.bug(1, 0) + g.in.bug(2, 0) + g.in.bug(3, 0) + g.in.bug(4, 0)
	case r == 2 && c == 3: // right of center
		cnt += g.in.bug(0, 4) + g.in.bug(1, 4) + g.in.bug(2, 4) + g.in.bug(3, 4) + g.in.bug(4, 4)
	}

	if r == 0 { // top edge
		cnt += g.out.bug(1, 2)
	} else if r == 4 { // bottom edge
		cnt += g.out.bug(3, 2)
	}
	if c == 0 { // left edge
		cnt += g.out.bug(2, 1)
	} else if c == 4 { // right edge
		cnt += g.out.bug(2, 3)
	}

	return cnt
}

func (g *rgrid) print(d dir, depth int) {
	if g == nil {
		return
	}

	fmt.Printf("Depth %d:\n", depth)
	fmt.Println(lib.DumpBytes(g.b))
	fmt.Println()

	if d&in != 0 {
		g.in.print(in, depth+1)
	}
	if d&out != 0 {
		g.out.print(out, depth-1)
	}
}
