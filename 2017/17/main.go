package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	steps := lib.InputInts("2017/17")[0]

	// Part 1
	type node struct {
		v    int
		next *node
	}
	cur := &node{v: 0}
	cur.next = cur
	for i := 1; i <= 2017; i++ {
		for j := 0; j < steps; j++ {
			cur = cur.next
		}
		n := &node{i, cur.next}
		cur.next = n
		cur = n
	}
	fmt.Println(cur.next.v)

	// Part 2: 50 million iterations.
	// The key insight here is that we only are about the value after 0,
	// so we only need to figure out when we'd be inserting there.
	var pos int  // current position
	var next int // value following 0
	cnt := 1     // number of values
	for v := 1; v <= 50_000_000; v++ {
		if pos = (pos + steps) % cnt; pos == 0 {
			next = v
		}
		cnt++                 // insert new value
		pos = (pos + 1) % cnt // new value becomes current
	}
	fmt.Println(next)
}
