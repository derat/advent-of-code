package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		nrows = 128
		ncols = 128
	)

	key := lib.InputLines("2017/14")[0]
	grid := lib.NewBytes(nrows, ncols, '.')
	for r := range grid {
		for i, b := range knot(fmt.Sprintf("%s-%d", key, r)) {
			for j := 0; j < 8; j++ {
				if b&(1<<(7-j)) != 0 {
					grid[r][i*8+j] = '#'
				}
			}
		}
	}

	// Part 1: Count the number of used squares.
	fmt.Println(lib.CountBytesFull(grid, '#'))

	// Part 2: Count the number of used regions.
	// I was originally thinking about scanning the grid one row at a time, keeping
	// track of active regions and incrementing a counter when they end. It seemed
	// like it'd be tricky to merge regions in patterns like the following, though,
	// so I instead just went with the approach of iterating through all blocks and
	// performing a search whenever we see one that's used that we haven't handled yet.
	var nregions int
	seen := make(map[uint64]struct{})
	for r, row := range grid {
		for c, ch := range row {
			k := lib.PackInts(r, c)
			if ch != '#' || lib.MapHasKey(seen, k) {
				continue
			}

			// Search for all the used blocks in this region.
			todo := [][2]int{[2]int{r, c}}
			for len(todo) > 0 {
				r, c := todo[0][0], todo[0][1]
				todo = todo[1:]
				for _, off := range [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
					r0, c0 := r+off[0], c+off[1]
					if r0 < 0 || r0 >= nrows || c0 < 0 || c0 >= ncols || grid[r0][c0] != '#' {
						continue
					}
					k := lib.PackInts(r0, c0)
					if lib.MapHasKey(seen, k) {
						continue
					}
					seen[k] = struct{}{}
					todo = append(todo, [2]int{r0, c0})
				}
			}
			nregions++
		}
	}
	fmt.Println(nregions)
}

// This is abstracted from 2017/10.
func knot(key string) []byte {
	const (
		nitems  = 256
		nrounds = 64
	)

	items := make([]int, nitems)
	for i := range items {
		items[i] = i
	}

	var pos, skip int
	lengths := append([]byte(key), 17, 31, 73, 47, 23) // from 2017/10
	for round := 0; round < nrounds; round++ {
		for _, length := range lengths {
			for i, j := pos, pos+int(length)-1; i < j; i, j = i+1, j-1 {
				items[i%nitems], items[j%nitems] = items[j%nitems], items[i%nitems]
			}
			pos = (pos + int(length) + skip) % nitems
			skip++
		}
	}

	// Reduce 256-byte sparse hash to 16-byte dense hash by XOR-ing each 16 values.
	hash := make([]byte, 16)
	for i := range hash {
		var b byte
		for j := 0; j < 16; j++ {
			b ^= byte(items[i*16+j])
		}
		hash[i] = b
	}
	return hash
}
