package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	orbits := make(map[string][]string) // vals are in orbit around keys
	for _, ln := range lib.InputLines("2019/6") {
		var a, b string
		lib.Extract(ln, `^([\w]+)\)([\w]+)$`, &a, &b)
		orbits[a] = append(orbits[a], b)
	}

	root := &node{"COM", nil}
	nodes := map[string]*node{root.name: root}
	next := []string{root.name}

	for len(next) > 0 {
		name := next[0]
		next = next[1:]
		n := nodes[name]
		lib.Assert(n != nil)
		for _, cname := range orbits[name] {
			c := &node{name: cname}
			nodes[cname] = c
			n.children = append(n.children, c)
			next = append(next, cname)
		}
	}

	var count func(*node, int) int
	count = func(n *node, depth int) int {
		sum := depth // count ourselves
		for _, c := range n.children {
			sum += count(c, depth+1)
		}
		return sum
	}
	fmt.Println(count(root, 0))
}

type node struct {
	name     string
	children []*node
}
