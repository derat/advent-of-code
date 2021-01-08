package main

import (
	"encoding/hex"
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const nitems = 256
	items := make([]int, nitems)
	for i := range items {
		items[i] = i
	}
	var pos, skip int
	for _, length := range lib.InputInts("2017/10") {
		lib.RotateSlice(items, -pos)
		lib.Reverse(items[:length])
		lib.RotateSlice(items, pos)
		pos = (pos + length + skip) % nitems
		skip++
	}
	fmt.Println(items[0] * items[1])

	// Part 2: Do a bunch more random junk. First, interpret input as bytes
	// and do 64 rounds instead of 1.
	const nrounds = 64
	for i := range items {
		items[i] = i
	}
	pos, skip = 0, 0
	lengths := lib.InputLinesBytes("2017/10")[0]
	lengths = append(lengths, 17, 31, 73, 47, 23) // hardcoded in puzzle
	for i := 0; i < nrounds; i++ {
		for _, length := range lengths {
			lib.RotateSlice(items, -pos)
			lib.Reverse(items[:length])
			lib.RotateSlice(items, pos)
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
	fmt.Println(hex.EncodeToString(hash))
}
