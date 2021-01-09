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
	fmt.Println(lib.CountBytesFull(grid, '#'))
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
