package main

import (
	"fmt"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	tracks := lib.InputLinesBytes("2018/13", '-', '|', '\\', '/', '+', ' ', '^', 'v', '<', '>')

	var carts []cart
	for r, row := range tracks {
		for c, ch := range row {
			switch ch {
			case '^':
				carts = append(carts, cart{r, c, up, 0})
				row[c] = '|'
			case 'v':
				carts = append(carts, cart{r, c, down, 0}) // for want of a ternary operator
				row[c] = '|'
			case '<':
				carts = append(carts, cart{r, c, left, 0})
				row[c] = '-'
			case '>':
				carts = append(carts, cart{r, c, right, 0})
				row[c] = '-'
			}
		}
	}

	locs := make(map[uint64]int) // packed r,c of carts
	for i, ca := range carts {
		locs[lib.PackInts(ca.r, ca.c)] = i
	}

	var tick int
Loop:
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

		for i := range carts {
			ca := &carts[i]

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
				lib.Assertf(ca.dir == up || ca.dir == down, "Cart %d not vert at %d,%d", i, ca.r, ca.c)
			case '-':
				lib.Assertf(ca.dir == left || ca.dir == right, "Card %d not horiz at %d,%d", i, ca.r, ca.c)
			}

			// Move the cart.
			delete(locs, lib.PackInts(ca.r, ca.c))

			ca.r += ca.dir.dr()
			ca.c += ca.dir.dc()

			// Check for collision.
			pos := lib.PackInts(ca.r, ca.c)
			if _, ok := locs[pos]; ok {
				fmt.Printf("%d,%d\n", ca.c, ca.r) // print as X,Y
				break Loop
			}
			locs[pos] = i
		}
	}
}

type cart struct {
	r, c int
	dir  dir
	ints int // intersection count
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
