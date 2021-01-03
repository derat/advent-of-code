package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const tlen = 272 // target length (given in puzzle)
	init := lib.InputLinesBytes("2016/16", '1', '0')[0]

	// Part 1: Naively generate the full data.
	b := append([]byte{}, init...)
	for len(b) < tlen {
		b0 := make([]byte, 2*len(b)+1)
		copy(b0, b)
		b0[len(b)] = '0'
		for i, ch := range b {
			if ch == '0' {
				ch = '1'
			} else {
				ch = '0'
			}
			b0[len(b0)-i-1] = ch
		}
		b = b0
	}
	fmt.Println(string(checksum(b[:tlen])))

	// Part 2: Disk size is much bigger!
	const tlen2 = 35651584
	fmt.Println(string(checksum2(init, tlen2)))
}

// checksum is a straightforward recursive implementation of the checksumming algorithm
// described in the problem.
func checksum(b []byte) []byte {
	csum := make([]byte, len(b)/2)
	for i := 0; i < len(csum); i++ {
		if b[2*i] == b[2*i+1] {
			csum[i] = '1'
		} else {
			csum[i] = '0'
		}
	}
	// "If the length of the checksum is even, repeat the process until
	// you end up with a checksum with an odd length."
	if len(csum)%2 == 0 {
		return checksum(csum)
	}
	return csum
}

// checksum2 returns the checksum produced after expanding init to tlen.
func checksum2(init []byte, tlen int) []byte {
	// Reversing and flipping state 'a' yields 'b'.
	// Note that reversing and flipping 'b' also yields 'a'.
	//
	// After repeated steps:
	//  a (initial)
	//  a0b
	//  a0b0a1b
	//  a0b0a1b0a0b1a1b
	//  a0b0a1b0a0b1a1b0a0b0a1b1a0b1a1b
	//
	// So we end up with 'a', a separator, 'b', a separator, 'a', a separator,
	// and so on. If we just look at the separators, they follow the same pattern:
	//
	//  0
	//  001
	//  0010011
	//  001001100011011
	//
	// The separator sequence is fixed -- 0 is always the initial state.

	// Determine how many times we'll need to reduce the full data to get
	// an odd-length checksum.
	red := 0
	for i := tlen; i%2 != 1; i /= 2 {
		red++
	}

	// Determine the length of each checksum block.
	blen := lib.Pow(2, red)

	// Each char of the checksum will be produced from blen chars of input.
	// For any n chars of input (where n is a power of 2), the resulting single
	// checksum char will be 1 if there are even numbers of 0s and 1s in n and
	// 0 if there are odd numbers.

	// Iterating through all of the chars and counting the ones seems wasteful,
	// but it's actually still pretty fast.
	//
	// I suspect that it'd be possible to optimize further by determining the number
	// of full 'a' and 'b' sequences in a given block and getting the parity of ones
	// that they contain, and then also counting the ones from any partial 'a' or 'b' at
	// the beginning or end of a block, plus the ones provided by separators in the
	// block. That seems like it'd provide a lot of opportunities for fiddly off-by-one
	// errors, though, so I'm just sticking with what I have here.
	csum := make([]byte, tlen/blen)
	for i := 0; i < len(csum); i++ {
		ones := 0
		start := i * blen
		for j := 0; j < blen; j++ {
			if ch := get(init, start+j); ch == '1' {
				ones++
			}
		}
		if ones%2 == 0 {
			csum[i] = '1'
		} else {
			csum[i] = '0'
		}
	}
	return csum
}

// get returns the char at position i given an initial state of a.
func get(a []byte, i int) byte {
	// The full sequence consists of a, followed by a separator, followed by b, followed
	// by another separator, followed by a, and so on.
	n := len(a) + 1
	div := i / n

	// If the byte is within a or b (its reversed/inverted counterpart), return it directly.
	if rem := i % n; rem < n-1 {
		if div%2 == 0 { // in a
			return a[rem]
		} else { // in b
			if ch := a[len(a)-rem-1]; ch == '1' {
				return '0'
			} else {
				return '1'
			}
		}
	}

	// Otherwise, we're in a separator. The separator follows the same expansion pattern
	// with a starting byte of 0, so we call ourselves recursively.
	return get([]byte{'0'}, div)
}
