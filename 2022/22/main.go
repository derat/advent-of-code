package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
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

	warp := make(map[vec]vec)

	// TODO: I tried a general approach but couldn't get it working.
	// I ended up just hardcoding the joined edges needed for my input. :-(

	/*
		warp := make(map[vec]vec)
		exists := func(p point) bool { return grid.Get(p.r, p.c, ' ') != ' ' }
		for r := 0; r < grid.Rows(); r++ {
			for c := 0; c < grid.Cols(); c++ {
				p := point{r, c}
				if !exists(p) {
					continue
				}
			NeighborLoop:
				for _, dir := range []lib.Dir{lib.Up, lib.Down, lib.Left, lib.Right} {
					if exists(p.step(dir)) || lib.MapHasKey(warp, vec{p, dir}) {
						continue
					}

					for _, tdir := range []lib.Dir{dir.Left(), dir.Right()} {
						for next := p.step(tdir); next.dist(p) <= 2*size && exists(next); next = next.step(tdir) {
							if exists(next.step(dir)) {
								dist := next.dist(p)
								for i := 0; i < dist; i++ {
									if i < size {
										next = next.step(dir)
									} else {
										next = next.step(tdir)
									}
								}
								ndir := tdir
								if dist > size {
									ndir = dir.Reverse()
								}
								warp[vec{p, dir}] = vec{next, ndir}
								warp[vec{next, ndir.Reverse()}] = vec{p, dir.Reverse()}
								continue NeighborLoop
							}
						}
					}
				}
			}
		}
	*/

	// TODO: This is disgusting and only works for my input.
	join := func(a point, as, ad lib.Dir, b point, bs, bd lib.Dir) {
		for i := 0; i < size; i++ {
			warp[vec{a, ad}] = vec{b, bd.Reverse()}
			warp[vec{b, bd}] = vec{a, ad.Reverse()}
			a = a.step(as)
			b = b.step(bs)
		}
	}

	join(point{0, size}, lib.Right, lib.Up, point{3 * size, 0}, lib.Down, lib.Left)
	join(point{0, size}, lib.Down, lib.Left, point{3*size - 1, 0}, lib.Up, lib.Left)
	join(point{2 * size, 0}, lib.Right, lib.Up, point{size, size}, lib.Down, lib.Left)
	join(point{0, 2 * size}, lib.Right, lib.Up, point{4*size - 1, 0}, lib.Right, lib.Down)
	join(point{3 * size, size - 1}, lib.Down, lib.Right, point{3*size - 1, size}, lib.Right, lib.Down)
	join(point{2 * size, 2*size - 1}, lib.Down, lib.Right, point{size - 1, 3*size - 1}, lib.Up, lib.Right)
	join(point{size - 1, 2 * size}, lib.Right, lib.Down, point{size, 2*size - 1}, lib.Down, lib.Right)

	/* For example input:
	join(point{size, 3*size - 1}, lib.Down, lib.Right, point{2 * size, 4*size - 1}, lib.Left, lib.Up)
	join(point{3*size - 1, 2 * size}, lib.Right, lib.Down, point{2*size - 1, size - 1}, lib.Left, lib.Down)
	join(point{size, size}, lib.Right, lib.Up, point{0, 2 * size}, lib.Down, lib.Left)
	*/

	move2 := func(r, c int, dir lib.Dir, steps int) (int, int, lib.Dir) {
		for ; steps > 0; steps-- {

			var nr, nc int
			var ndir lib.Dir
			if dst, ok := warp[vec{point{r, c}, dir}]; ok {
				nr, nc, ndir = dst.p.r, dst.p.c, dst.d
			} else {
				nr, nc, ndir = r+dir.DR(), c+dir.DC(), dir
			}

			ch := grid.Get(nr, nc, ' ')
			switch ch {
			case '#':
				break // stopped by wall
			case '.', '<', '>', '^', 'v':
				//grid[r][c] = dirChar[dir]
				//fmt.Println(grid.Dump())
				//fmt.Println()
				r, c, dir = nr, nc, ndir // empty space
			default:
				lib.Panicf("Moved to %q at row %d, col %d", ch, nr, nc)
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

type point struct{ r, c int }

func (p point) stepn(dir lib.Dir, n int) point { return point{p.r + n*dir.DR(), p.c + n*dir.DC()} }
func (p point) step(dir lib.Dir) point         { return p.stepn(dir, 1) }
func (p point) dist(o point) int               { return lib.Abs(p.r-o.r) + lib.Abs(p.c-o.c) }

type vec struct {
	p point
	d lib.Dir
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
