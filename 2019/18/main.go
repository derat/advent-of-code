package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	grid := lib.InputLinesBytes("2019/18")

	var sr, sc int                // starting position
	keyLocs := make([][2]int, 26) // index is id, vals are r,c
	for r, row := range grid {
		for c, ch := range row {
			switch {
			case ch == '@': // entrance
				sr, sc = r, c
			case ch >= 'a' && ch <= 'z': // key
				keyLocs[ch-'a'] = [2]int{r, c}
			}
		}
	}

	allKeys := (1 << 26) - 1 // bitfield representing all keys

	cost := lib.AStar([]uint64{pack(sr, sc, 0)},
		func(s uint64) bool {
			_, _, keys := unpack(s)
			return keys&allKeys == allKeys
		},
		func(s uint64) []uint64 {
			var next []uint64
			r, c, keys := unpack(s)
			for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
				nr, nc := r+off[0], c+off[1]
				nkeys := keys
				ch := grid[nr][nc]
				switch {
				case ch == '#': // wall
					continue
				case ch >= 'a' && ch <= 'z': // key
					nkeys |= 1 << (ch - 'a')
				case ch >= 'A' && ch <= 'Z': // door
					if keys&(1<<(ch-'A')) == 0 {
						continue // don't have key
					}
				case ch == '.': // empty space
				case ch == '@': // entrance
				default:
					lib.Panicf("Invalid char %q at %d, %d", ch, nr, nc)
				}
				next = append(next, pack(nr, nc, nkeys))
			}
			return next
		},
		func(s uint64) int {
			// Use the Manhattan distance to the farthest key as a lower bound.
			r, c, keys := unpack(s)
			var max int
			for id, loc := range keyLocs {
				if keys&id == 0 { // don't have the key yet
					max = lib.Max(lib.Abs(loc[0]-r)+lib.Abs(loc[1]-c), max)
				}
			}
			return max
		})
	fmt.Println(cost)
}

func pack(r, c, keys int) uint64 {
	// Use 32 bits to track keys and 16 each for the row and column.
	p := lib.PackInt(0, r, 16, 0)
	p = lib.PackInt(p, c, 16, 16)
	p = lib.PackInt(p, keys, 32, 32)
	return p
}

func unpack(p uint64) (r, c, keys int) {
	r = lib.UnpackInt(p, 16, 0)
	c = lib.UnpackInt(p, 16, 16)
	keys = lib.UnpackInt(p, 32, 32)
	return r, c, keys
}
