package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const tlen = 272 // target length (given in puzzle)
	init := lib.InputLinesBytes("2016/16", '1', '0')[0]
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
}

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
