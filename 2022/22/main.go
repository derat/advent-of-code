package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// My solution is horrendous enough that I needed to write tests to find the bugs. :-(
	test()

	pgs := lib.InputParagraphs("2022/22")
	lib.AssertEq(len(pgs), 2)

	// The input lines annoyingly aren't all the same width, so manually create a grid filled with
	// spaces and copy the input into it.
	var mc int
	for _, ln := range pgs[0] {
		mc = lib.Max(mc, len(ln))
	}
	grid := lib.NewByteGrid(len(pgs[0]), mc, ' ')
	for r, ln := range pgs[0] {
		copy(grid[r], []byte(ln))
	}

	var cmin, cmax []int // min and max columns for each row
	for r := range grid {
		min, max := 0, grid.MaxCol()
		for ; min < grid.Cols() && grid[r][min] == ' '; min++ {
		}
		for ; max >= 0 && grid[r][max] == ' '; max-- {
		}
		lib.AssertInRange(min, 0, grid.MaxCol())
		lib.AssertInRange(max, 0, grid.MaxCol())
		cmin = append(cmin, min)
		cmax = append(cmax, max)
	}

	var rmin, rmax []int // min and max rows for each column
	for c := 0; c < grid.Cols(); c++ {
		min, max := 0, grid.MaxRow()
		for ; min < grid.Rows() && grid[min][c] == ' '; min++ {
		}
		for ; max >= 0 && grid[max][c] == ' '; max-- {
		}
		lib.AssertInRange(min, 0, grid.MaxRow())
		lib.AssertInRange(max, 0, grid.MaxRow())
		rmin = append(rmin, min)
		rmax = append(rmax, max)
	}

	move := func(r, c int, dir lib.Dir, steps int) (int, int) {
		for ; steps > 0; steps-- {
			nr, nc := r+dir.DR(), c+dir.DC()
			if dir.DR() != 0 { // vertical
				if nr < rmin[nc] {
					nr = rmax[nc] // wrapped top
				} else if nr > rmax[nc] {
					nr = rmin[nc] // wrapped bottom
				}
			} else { // horizontal
				if nc < cmin[nr] {
					nc = cmax[nr] // wrapped left
				} else if nc > cmax[nr] {
					nc = cmin[nr] // wrapped bottom
				}
			}
			ch := grid[nr][nc]
			switch ch {
			case '#':
				break // stopped by wall
			case '.':
				r, c = nr, nc // empty space
			default:
				lib.Panicf("Moved to %q at row %d, col %d", ch, nr, nc)
			}
		}
		return r, c
	}

	// Part 2: The example specifies a cube that's unfolded like this:
	//
	//    1
	//  234
	//    56
	//
	// My input is like this, though:
	//
	//   12
	//   3
	//  45
	//  6
	//
	// There are various other possible configurations, e.g.
	//
	//  123
	//   4
	//   5
	//   6
	//
	//  12
	//   34
	//    56
	//
	// The prompt says that each face is 50x50 (except in the example input, where they're 4x4).
	// I could move and rotate faces from the original input map, but I'd still need to use the
	// original coordinates to produce the answer.
	//
	// If we're moving off the edge of the map and the current row/column spans 4 times the size,
	// then we just wrap around without changing direction. This doesn't appear in the example or
	// in my input, though. The other cases are trickier!

	// Determine the dimensions of the cube faces. I think that the flattened cube needs to be 3x4
	// or 3x3 faces.
	var size int
	if grid.Rows() != grid.Cols() {
		lib.AssertEq(lib.Max(grid.Rows(), grid.Cols())/4, lib.Min(grid.Rows(), grid.Cols())/3)
		size = lib.Max(grid.Rows(), grid.Cols()) / 4
	} else {
		lib.AssertEq(grid.Rows()%3, 0)
		size = grid.Rows() / 3
	}

	move2 := func(r, c int, dir lib.Dir, steps int) (int, int, lib.Dir) {
		for ; steps > 0; steps-- {
			nr, nc, ndir := step(grid, size, r, c, dir)
			ch := grid.Get(nr, nc, ' ')
			switch ch {
			case '#':
				break // stopped by wall
			case ' ':
				lib.Panicf("Bad move [%v, %v] %c to [%v, %v] %c",
					r, c, dirChar[dir], nr, nc, dirChar[ndir])
			default:
				//grid[r][c] = dirChar[dir]
				//fmt.Println(grid.Dump())
				//fmt.Println()
				r, c, dir = nr, nc, ndir // empty space
			}
		}
		return r, c, dir
	}

	// "You begin the path in the leftmost open tile of the top row of tiles.
	// Initially, you are facing to the right (from the perspective of how the map is drawn)."
	var r, c int
	dir := lib.Right
	for ; c < grid.Rows() && grid[r][c] != '.'; c++ {
	}
	lib.AssertLess(c, grid.Rows())

	r2, c2 := r, c
	dir2 := dir

	lib.AssertEq(len(pgs[1]), 1)
	for _, cmd := range lib.Tokenize(pgs[1][0], regexp.MustCompile(`^\d+|[LR]`)) {
		switch cmd {
		case "L":
			dir = dir.Left()
			dir2 = dir2.Left()
		case "R":
			dir = dir.Right()
			dir2 = dir2.Right()
		default:
			steps := lib.ExtractInt(cmd)
			r, c = move(r, c, dir, steps)
			r2, c2, dir2 = move2(r2, c2, dir2, steps)
		}
	}

	// Part 1: Generate a password from the final position.
	fmt.Println(password(r, c, dir))

	// Part 2: Use cube-like wraparound rules.
	fmt.Println(password(r2, c2, dir2))
}

// vec tracks a point and direction moving around a grid.
// It is unaware of cube-mapping and just moves around the 2-dimensional grid.
// It tracks whether it has ever moved to an invalid position outside of a cube face.
type vec struct {
	g    lib.ByteGrid
	r, c int
	d    lib.Dir
	ok   bool // moved to invalid position during a forward() call
}

// Move forward in the current direction by n spaces.
func (v vec) fore(n int) vec {
	for i := 0; i < n; i++ {
		v.r += v.d.DR()
		v.c += v.d.DC()
		if v.g.Get(v.r, v.c, ' ') == ' ' {
			v.ok = false
		}
	}
	return v
}

// Turn in the specified direction and move n spaces.
func (v vec) left(n int) vec  { return vec{v.g, v.r, v.c, v.d.Left(), v.ok}.fore(n) }
func (v vec) right(n int) vec { return vec{v.g, v.r, v.c, v.d.Right(), v.ok}.fore(n) }
func (v vec) rev(n int) vec   { return vec{v.g, v.r, v.c, v.d.Reverse(), v.ok}.fore(n) }

// step moves one step from [r, c] in the specified direction,
// wrapping around to different faces of the cube when needed.
func step(grid lib.ByteGrid, sz, r, c int, dir lib.Dir) (int, int, lib.Dir) {
	start := vec{grid, r, c, dir, true /* ok */}

	// Simple case: the next face is in front of us.
	if v := start.fore(1); v.ok {
		return v.r, v.c, v.d
	}

	var rs int
	for ; true; rs++ {
		if v := start.right(rs + 1); v.r < 0 || v.r/sz != r/sz || v.c < 0 || v.c/sz != c/sz {
			break
		}
	}
	ls := sz - rs - 1

	// This tedious list of the different paths that we can take to move forward across various
	// flattened cubes is really gross. I haven't completely convinced myself that it's exhaustive,
	// and I'm sure that there's a much easier approach. :-(
	for _, v := range []vec{
		start.rev(4*sz - 1).rev(0), // bottom, back, top

		start.right(rs + 1).left(1 + rs).right(0),                                            // right, top
		start.right(rs + 1 + sz + rs).left(sz).rev(0),                                        // right, back, top
		start.right(rs + 3*sz).left(sz - rs).left(0),                                         // right, back, left, top
		start.right(rs + 1).right(sz).left(3*sz - 1).right(rs).right(0),                      // right, bottom, back, top
		start.right(rs + 1).right(2 * sz).left(sz).right(2*sz - 1).left(sz - rs - 1).left(0), // right, bottom, left, back, top
		start.rev(sz).left(rs + 1 + sz + rs).right(0),                                        // bottom, right, top
		start.rev(sz).left(rs + 1).right(sz).left(2*sz - 1).right(rs).right(0),               // bottom, right, back, top
		start.rev(sz).left(rs + 1).right(3*sz - 1).left(2*sz - rs - 1).left(0),               // bottom, right, back, left, top
		start.rev(2 * sz).left(sz + rs).right(sz + rs).right(0),                              // bottom, back, right, top

		// These are the same as the previous set, but with right/left and ls/rs swapped.
		start.left(ls + 1).right(1 + ls).left(0),                                             // left, top
		start.left(ls + 1 + sz + ls).right(sz).rev(0),                                        // left, back, top
		start.left(ls + 3*sz).right(sz - ls).right(0),                                        // left, back, right, top
		start.left(ls + 1).left(sz).right(3*sz - 1).left(ls).left(0),                         // left, bottom, back, top
		start.left(ls + 1).left(2 * sz).right(sz).left(2*sz - 1).right(sz - ls - 1).right(0), // left, bottom, right, back, top
		start.rev(sz).right(ls + 1 + sz + ls).left(0),                                        // bottom, left, top
		start.rev(sz).right(ls + 1).left(sz).right(2*sz - 1).left(ls).left(0),                // bottom, left, back, top
		start.rev(sz).right(ls + 1).left(3*sz - 1).right(2*sz - ls - 1).right(0),             // bottom, left, back, right, top
		start.rev(2 * sz).right(sz + ls).left(sz + ls).left(0),                               // bottom, back, left, top
	} {
		if v.ok {
			return v.r, v.c, v.d
		}
	}

	lib.Panicf("Failed step from [%d, %d]\n", r, c)
	return r, c, dir // not reached
}

// "Rows start from 1 at the top and count downward; columns start from 1 at the left and count
// rightward. ... The final password is the sum of 1000 times the row, 4 times the column, and
// the facing."
func password(r, c int, dir lib.Dir) int { return 1000*(r+1) + 4*(c+1) + dirValue[dir] }

// "Facing is 0 for right (>), 1 for down (v), 2 for left (<), and 3 for up (^)."
var dirValue = map[lib.Dir]int{
	lib.Right: 0,
	lib.Down:  1,
	lib.Left:  2,
	lib.Up:    3,
}

var dirChar = map[lib.Dir]byte{
	lib.Right: '>',
	lib.Down:  'v',
	lib.Left:  '<',
	lib.Up:    '^',
}

func test() {
	grid := lib.NewByteGridString(strings.Join([]string{
		// 2345678
		"   111222", // 0
		"   111222", // 1
		"   111222", // 2
		"   333   ", // 3
		"   333   ", // 4
		"   333   ", // 5
		"444555   ", // 6
		"444555   ", // 7
		"444555   ", // 8
		"666      ", // 9
		"666      ", // 10
		"666      ", // 11
	}, "\n"))
	const size = 3

	for _, tc := range []struct {
		sr, sc int
		sd     lib.Dir
		dr, dc int
		dd     lib.Dir
	}{
		// up from 1: bottom, back, left, top
		{0, 3, lib.Up, 9, 0, lib.Right},
		{0, 4, lib.Up, 10, 0, lib.Right},
		{0, 5, lib.Up, 11, 0, lib.Right},
		// left from 1: left, back, top
		{0, 3, lib.Left, 8, 0, lib.Right},
		{1, 3, lib.Left, 7, 0, lib.Right},
		{2, 3, lib.Left, 6, 0, lib.Right},
		// up from 2: left, bottom, right, back, top
		{0, 6, lib.Up, 11, 0, lib.Up},
		{0, 7, lib.Up, 11, 1, lib.Up},
		{0, 8, lib.Up, 11, 2, lib.Up},
		// right from 2: bottom, right, top
		{0, 8, lib.Right, 8, 5, lib.Left},
		{1, 8, lib.Right, 7, 5, lib.Left},
		{2, 8, lib.Right, 6, 5, lib.Left},
		// down from 2: right, top
		{2, 6, lib.Down, 3, 5, lib.Left},
		{2, 7, lib.Down, 4, 5, lib.Left},
		{2, 8, lib.Down, 5, 5, lib.Left},
		// left from 3: left, top
		{3, 3, lib.Left, 6, 0, lib.Down},
		{4, 3, lib.Left, 6, 1, lib.Down},
		{5, 3, lib.Left, 6, 2, lib.Down},
		// right from 3: left, top (already tested)
		// up from 4: right, top (already tested)
		// left from 4: bottom, right, top (already tested)
		// right from 5: left, back, top (already tested)
		// down from 5: right, top (already tested)
		// left from 6: right, bottom, back, top
		{9, 0, lib.Left, 0, 3, lib.Down},
		{10, 0, lib.Left, 0, 4, lib.Down},
		{11, 0, lib.Left, 0, 5, lib.Down},
		// right from 6: left, top (already tested)
		// down from 6: bottom, left, back, right, top
		{11, 0, lib.Down, 0, 6, lib.Down},
		{11, 1, lib.Down, 0, 7, lib.Down},
		{11, 2, lib.Down, 0, 8, lib.Down},
		// TODO: This only tests 8 of the 20 possible paths.
		// I just focused on the paths needed for my input.
	} {
		if r, c, d := step(grid, size, tc.sr, tc.sc, tc.sd); r != tc.dr || c != tc.dc || d != tc.dd {
			lib.Panicf("[%v, %v] %c = [%v, %v] %c; want [%v, %v] %c",
				tc.sr, tc.sc, dirChar[tc.sd],
				r, c, dirChar[d],
				tc.dr, tc.dc, dirChar[tc.dd])
		}
	}
}
