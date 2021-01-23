package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	tracks := lib.InputLinesBytes("2018/13", '-', '|', '\\', '/', '+', ' ', '^', 'v', '<', '>')

	var carts []*cart
	for r, row := range tracks {
		for c, ch := range row {
			switch ch {
			case '^':
				carts = append(carts, &cart{r: r, c: c, dir: lib.Up})
				row[c] = '|'
			case 'v':
				carts = append(carts, &cart{r: r, c: c, dir: lib.Down}) // for want of a ternary operator
				row[c] = '|'
			case '<':
				carts = append(carts, &cart{r: r, c: c, dir: lib.Left})
				row[c] = '-'
			case '>':
				carts = append(carts, &cart{r: r, c: c, dir: lib.Right})
				row[c] = '-'
			}
		}
	}

	// Sigh, I spent way too long on a bug here.
	// I was storing cart indexes rather than pointers here,
	// so things got completely messed up after sorting in the loop.
	locs := make(map[uint64]*cart) // packed r,c of carts
	for _, ca := range carts {
		locs[lib.PackInts(ca.r, ca.c)] = ca
	}

	var tick int
	var ncol int
	for tick = 1; true; tick++ {
		// Sort carts from top to bottom and left to right.
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].r < carts[j].r {
				return true
			} else if carts[i].r > carts[j].r {
				return false
			}
			return carts[i].c < carts[j].c
		})

		for _, ca := range carts {
			// Skip removed carts.
			if ca.rem {
				continue
			}

			// Turn if needed.
			switch tracks[ca.r][ca.c] {
			case '/':
				if ca.dir == lib.Up || ca.dir == lib.Down {
					ca.dir = ca.dir.Right()
				} else {
					ca.dir = ca.dir.Left()
				}
			case '\\':
				if ca.dir == lib.Up || ca.dir == lib.Down {
					ca.dir = ca.dir.Left()
				} else {
					ca.dir = ca.dir.Right()
				}
			case '+':
				switch ca.ints % 3 {
				case 0:
					ca.dir = ca.dir.Left()
				case 1: // go straight
				case 2:
					ca.dir = ca.dir.Right()
				}
				ca.ints++
			case '|':
				lib.Assertf(ca.dir == lib.Up || ca.dir == lib.Down, "Cart not vert at %d,%d", ca.r, ca.c)
			case '-':
				lib.Assertf(ca.dir == lib.Left || ca.dir == lib.Right, "Cart not horiz at %d,%d", ca.r, ca.c)
			}

			// Move the cart.
			delete(locs, lib.PackInts(ca.r, ca.c))
			ca.r += ca.dir.DR()
			ca.c += ca.dir.DC()

			// Check for collision.
			pos := lib.PackInts(ca.r, ca.c)
			if other, ok := locs[pos]; ok {
				// Part 1: Print X,Y of the first collision.
				ncol++
				if ncol == 1 {
					fmt.Printf("%d,%d\n", ca.c, ca.r)
				}

				// Remove both carts.
				ca.rem = true
				other.rem = true
				delete(locs, pos)
			} else {
				locs[pos] = ca
			}
		}

		// Part 2: Print X,Y of last cart at end of first tick where it's the only one left.
		if ncol == len(carts)/2 {
			for _, ca := range carts {
				if !ca.rem {
					fmt.Printf("%d,%d\n", ca.c, ca.r)
					break
				}
			}
			break
		}
	}
}

type cart struct {
	r, c int
	dir  lib.Dir
	ints int  // intersection count
	rem  bool // removed after collision
}
