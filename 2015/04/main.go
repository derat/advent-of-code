package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strconv"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const zeros = 5               // number of leading zeros
	want := make([]byte, zeros/2) // full zero bytes

	done := false // already printed a solution for part1

	pre := lib.InputLines("2015/4")[0]
	for i := 0; true; i++ {
		s := pre + strconv.Itoa(i)
		hash := md5.Sum([]byte(s))
		if bytes.HasPrefix(hash[:], want) && (zeros%2 == 0 || hash[zeros/2] < 16) {
			if !done {
				fmt.Println(i)
				done = true
			}
			// Part 2: Check the full last byte. (We know zeroes is odd.)
			if hash[zeros/2] == 0 {
				fmt.Println(i)
				break
			}
		}
	}
}
