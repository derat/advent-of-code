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

	root := &node{"COM", nil, nil}
	nodes := map[string]*node{root.name: root}
	next := []string{root.name}

	for len(next) > 0 {
		name := next[0]
		next = next[1:]
		n := nodes[name]
		lib.Assert(n != nil)
		for _, cname := range orbits[name] {
			c := &node{name: cname, parent: n}
			nodes[cname] = c
			n.children = append(n.children, c)
			next = append(next, cname)
		}
	}

	// Part 1: Print total number of direct and indirect orbits.
	var count func(*node, int) int
	count = func(n *node, depth int) int {
		sum := depth // count ourselves
		for _, c := range n.children {
			sum += count(c, depth+1)
		}
		return sum
	}
	fmt.Println(count(root, 0))

	// Part 2: Print number of orbital transfers needed to go from YOU to SAN.
	// (Really seems like the distance from YOU's parent to SAN's parent, though.)
	ancestors := make(map[string]int)
	for n, c := nodes["YOU"].parent, 0; n != nil; n, c = n.parent, c+1 {
		ancestors[n.name] = c
	}
	for n, c := nodes["SAN"].parent, 0; n != nil; n, c = n.parent, c+1 {
		if v, ok := ancestors[n.name]; ok {
			fmt.Println(c + v)
			break
		}
	}
}

type node struct {
	name     string
	children []*node
	parent   *node
}
