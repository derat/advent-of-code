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
				carts = append(carts, &cart{r: r, c: c, dir: up})
				row[c] = '|'
			case 'v':
				carts = append(carts, &cart{r: r, c: c, dir: down}) // for want of a ternary operator
				row[c] = '|'
			case '<':
				carts = append(carts, &cart{r: r, c: c, dir: left})
				row[c] = '-'
			case '>':
				carts = append(carts, &cart{r: r, c: c, dir: right})
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
				if ca.dir == up || ca.dir == down {
					ca.dir = ca.dir.right()
				} else {
					ca.dir = ca.dir.left()
				}
			case '\\':
				if ca.dir == up || ca.dir == down {
					ca.dir = ca.dir.left()
				} else {
					ca.dir = ca.dir.right()
				}
			case '+':
				switch ca.ints % 3 {
				case 0:
					ca.dir = ca.dir.left()
				case 1: // go straight
				case 2:
					ca.dir = ca.dir.right()
				}
				ca.ints++
			case '|':
				lib.Assertf(ca.dir == up || ca.dir == down, "Cart not vert at %d,%d", ca.r, ca.c)
			case '-':
				lib.Assertf(ca.dir == left || ca.dir == right, "Cart not horiz at %d,%d", ca.r, ca.c)
			}

			// Move the cart.
			delete(locs, lib.PackInts(ca.r, ca.c))
			ca.r += ca.dir.dr()
			ca.c += ca.dir.dc()

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
	dir  dir
	ints int  // intersection count
	rem  bool // removed after collision
}

// Copied from 2017/22.
type dir int

const (
	up dir = iota
	left
	down
	right
)

func (d dir) left() dir {
	return dir((int(d) + 1) % 4)
}

func (d dir) right() dir {
	return dir((int(d) + 3) % 4)
}

func (d dir) reverse() dir {
	return dir((int(d) + 2) % 4)
}

func (d dir) dr() int {
	switch d {
	case up:
		return -1
	case down:
		return 1
	case left, right:
		return 0
	default:
		panic("Invalid dir")
	}
}

func (d dir) dc() int {
	switch d {
	case up, down:
		return 0
	case left:
		return -1
	case right:
		return 1
	default:
		panic("Invalid dir")
	}
}
