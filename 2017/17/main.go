package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	steps := lib.InputInts("2017/17")[0]
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
}

type node struct {
	v    int
	next *node
}
