package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		ncols = 25
		nrows = 6
		lsize = ncols * nrows
	)

	input := lib.InputLines("2019/8")[0]
	lib.AssertEq(len(input)%lsize, 0)
	nlayers := len(input) / lsize
	layers := make([][][]byte, nlayers)
	for i := range layers {
		layers[i] = lib.NewBytes(nrows, ncols, 0)
		for r := range layers[i] {
			copy(layers[i][r], []byte(input[i*lsize+r*ncols:]))
		}
	}

	minZeros := math.MaxInt32
	prod := 0
	for _, l := range layers {
		if zeros := lib.CountBytesFull(l, '0'); zeros < minZeros {
			minZeros = zeros
			// It'd be faster to count all three bytes at once, but *shrug*.
			prod = lib.CountBytesFull(l, '1') * lib.CountBytesFull(l, '2')
		}
	}
	fmt.Println(prod)
}
