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
	layers := make([]lib.ByteGrid, nlayers)
	for i := range layers {
		layers[i] = lib.NewByteGrid(nrows, ncols, 0)
		for r := range layers[i] {
			copy(layers[i][r], []byte(input[i*lsize+r*ncols:]))
		}
	}

	// Part 1: Print number of 1s times number of 2s from layer with fewest 0s.
	minZeros := math.MaxInt32
	prod := 0
	for _, l := range layers {
		if zeros := l.Count('0'); zeros < minZeros {
			minZeros = zeros
			// It'd be faster to count all three bytes at once, but *shrug*.
			prod = l.Count('1') * l.Count('2')
		}
	}
	fmt.Println(prod)

	// Part 2: Stack layers: first in front; 0 is black, 1 is white, 2 is transparent.
	img := lib.NewByteGrid(nrows, ncols, ' ')
	for r, row := range img {
		for c := range row {
		Loop:
			for i, l := range layers {
				ch := l[r][c]
				switch ch {
				case '0': // black
					row[c] = ' '
					break Loop
				case '1': // white
					row[c] = '#'
					break Loop
				case '2': // transparent
				default:
					lib.Panicf("Invalid char %q at %v,%v in layer %v", ch, r, c, i)
				}
			}
		}
	}
	//fmt.Println(img.Dump())
	// Hack: lib.OCR expects a blank column between letters, but
	// my input has the second letter flush against the third letter.
	left := lib.OCR(img.CopyRect(0, 0, nrows-1, 11), ' ')
	right := lib.OCR(img.CopyRect(0, 12, nrows-1, ncols-1), ' ')
	fmt.Println(left + right)
}
