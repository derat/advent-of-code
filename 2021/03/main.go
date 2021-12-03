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
	var gamma uint32
	for _, cnt := range ones {
		gamma <<= 1
		if cnt >= vals/2 {
			gamma |= 1
		}
	}
	epsilon := ^gamma & ((1 << len(ones)) - 1) // yuck
	fmt.Println(gamma * epsilon)

	// Part 2: Multiply oxygen (number with most-common bits, 1 wins in ties) and
	// CO2 (number with least-common bits, 0 wins in ties).
	findNum := func(mostCommon bool) uint32 {
		var v uint32
		n := &root
		for range ones {
			v <<= 1
			if n.zero.getCnt() == 0 ||
				(mostCommon && n.one.getCnt() >= n.zero.getCnt()) ||
				(!mostCommon && n.one.getCnt() > 0 && n.one.getCnt() < n.zero.getCnt()) {
				v |= 1
				lib.Assert(n.one.getCnt() > 0)
				n = n.one
			} else {
				lib.Assert(n.zero.getCnt() > 0)
				n = n.zero
			}
		}
		return v
	}
	oxygen := findNum(true)
	co2 := findNum(false)
	fmt.Println(oxygen * co2)
}

type node struct {
	cnt       int // 0 for root
	zero, one *node
}

// getCnt returns n.cnt if n is non-nil or 0 otherwise.
func (n *node) getCnt() int {
	if n == nil {
		return 0
	}
	return n.cnt
}
