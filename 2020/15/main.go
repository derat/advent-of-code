package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	const (
		end  = 2020
		end2 = 30000000
	)

	turn := 1
	seen := make(map[int]int) // turn in which number last spoken
	last := -1                // last number spoken
	for _, n := range lib.InputInts("2020/15") {
		last = n
		seen[n] = turn
		turn++
	}

	next := 0 // next number to speak (i.e. age of last number)
	for {
		last = next
		if lastTurn, ok := seen[last]; ok {
			next = turn - lastTurn
		} else {
			next = 0
		}
		seen[last] = turn

		if turn == end {
			fmt.Println(last)
		}
		if turn == end2 {
			fmt.Println(last)
			break
		}
		turn++
	}
}
