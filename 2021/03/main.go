package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Part 1
	var vals int   // number of input lines
	var ones []int // index 0 is MSB

	// Part 2: Store bit patterns in a trie.
	type node struct {
		cnt       int // 0 for root
		zero, one *node
	}
	var root node

	for _, ln := range lib.InputLines("2021/3") {
		// Part 1
		if ones == nil {
			ones = make([]int, len(ln))
		}
		for i, ch := range ln {
			if ch == '1' {
				ones[i]++
			}
		}
		vals++

		// Part 2
		n := &root
		for _, ch := range ln {
			if ch == '1' {
				if n.one == nil {
					n.one = &node{}
				}
				n = n.one
			} else {
				if n.zero == nil {
					n.zero = &node{}
				}
				n = n.zero
			}
			n.cnt++
		}
	}

	// Part 1: Multiply gamma (most common bits) and epsilon (least common bits).
	var gamma, epsilon uint32
	for _, cnt := range ones {
		gamma <<= 1
		epsilon <<= 1
		if cnt >= vals/2 {
			gamma |= 1
		} else {
			epsilon |= 1
		}
	}
	fmt.Println(uint64(gamma) * uint64(epsilon))

	// Part 2: Multiply oxygen (number with most-common bits, 1 wins in ties) and
	// CO2 (number with least-common bits, 0 wins in ties).
	findNum := func(mostCommon bool) uint32 {
		var v uint32
		n := &root
		for range ones {
			onlyOne := n.one != nil && n.zero == nil
			onlyZero := n.zero != nil && n.one == nil
			onesWin := n.one != nil && (n.zero == nil || n.one.cnt > n.zero.cnt)
			zerosWin := n.zero != nil && (n.one == nil || n.zero.cnt > n.one.cnt)
			tie := !onesWin && !zerosWin

			v <<= 1
			if !onlyZero && (onlyOne || (mostCommon && (onesWin || tie)) || (!mostCommon && zerosWin)) {
				v |= 1
				n = n.one
			} else {
				lib.Assert(onlyZero || (mostCommon && zerosWin) || (!mostCommon && (onesWin || tie)))
				n = n.zero
			}
		}
		return v
	}
	oxygen := findNum(true)
	co2 := findNum(false)
	fmt.Println(uint64(oxygen) * uint64(co2))
}
