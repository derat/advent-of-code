package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	black := make(map[uint64]struct{})
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}

		var x, y int32
		for i := 0; i < len(ln); i++ {
			switch ln[i] {
			case 'e':
				x++
			case 'w':
				x--
			case 'n':
				switch ln[i+1] {
				case 'e':
					x += ex(y)
				case 'w':
					x += wx(y)
				default:
					log.Fatalf("Bad line %q", ln)
				}
				y++
				i++
			case 's':
				switch ln[i+1] {
				case 'e':
					x += ex(y)
				case 'w':
					x += wx(y)
				default:
					log.Fatalf("Bad line %q", ln)
				}
				y--
				i++
			default:
				log.Fatalf("Bad line %q", ln)
			}
		}

		key := makeKey(x, y)
		if _, ok := black[key]; ok {
			delete(black, key) // flip to white
		} else {
			black[key] = struct{}{} // flip to black
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	fmt.Println(len(black))

	const days = 100
	for i := 0; i < days; i++ {
		nextBlack := make(map[uint64]struct{}) // black tiles for next iteration
		doneWhite := make(map[uint64]struct{}) // white tiles already updated

		// Checks the white tile at (x, y).
		updateWhite := func(x, y int32) {
			key := makeKey(x, y)
			if _, ok := black[key]; ok {
				return // not white
			}
			if _, ok := doneWhite[key]; ok {
				return // already handled
			}
			// White with 2 black adj is flipped.
			if countAdj(black, x, y) == 2 {
				nextBlack[key] = struct{}{}
			}
			doneWhite[key] = struct{}{}
		}

		for key := range black {
			x, y := unpackKey(key)
			// Black with 0 or >2 black adj is flipped, so only preserve if 1 or 2.
			if cnt := countAdj(black, x, y); cnt == 1 || cnt == 2 {
				nextBlack[key] = struct{}{}
			}

			// Check neighbors.
			updateWhite(x+1, y)       // e
			updateWhite(x-1, y)       // w
			updateWhite(x+ex(y), y+1) // ne
			updateWhite(x+wx(y), y+1) // nw
			updateWhite(x+ex(y), y-1) // se
			updateWhite(x+wx(y), y-1) // sw
		}

		black = nextBlack
	}

	fmt.Println(len(black))
}

func makeKey(x, y int32) uint64 {
	return uint64(x)<<32 | (uint64(y) & 0xffffffff)
}

func unpackKey(k uint64) (x, y int32) {
	x = int32((k >> 32) & 0xffffffff)
	y = int32(k & 0xffffffff)
	return
}

// X offset for NE or SE move starting from y.
func ex(y int32) int32 {
	if y < 0 {
		y = -y
	}
	return int32(y % 2)
}

// X offset for NW or SW move starting from y.
func wx(y int32) int32 {
	if y < 0 {
		y = -y
	}
	return -int32((y + 1) % 2)
}

// Returns number of black tiles adjacent to (x, y).
func countAdj(m map[uint64]struct{}, x, y int32) int {
	var cnt int
	check := func(x, y int32) {
		if _, ok := m[makeKey(x, y)]; ok {
			cnt++
		}
	}
	check(x+1, y)       // e
	check(x-1, y)       // w
	check(x+ex(y), y+1) // ne
	check(x+wx(y), y+1) // nw
	check(x+ex(y), y-1) // se
	check(x+wx(y), y-1) // sw
	return cnt
}
