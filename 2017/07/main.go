package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	progs := make(map[string]*prog)
	for _, ln := range lib.InputLines("2017/7") {
		p := &prog{children: make(map[string]*prog)}
		var rest string
		lib.Extract(ln, `^([a-z]+) \((\d+)\)(.*)`, &p.name, &p.weight, &rest)

		if rest != "" {
			var list string
			lib.Extract(rest, `^ -> (.+)$`, &list)
			for _, s := range strings.Split(list, ", ") {
				p.children[s] = nil
			}
		}

		progs[p.name] = p
	}

	// For part 1, the root is the program that has children but isn't a child of any other program.
	// Fix up child pointers while we're iterating, because why not.
	children := make(map[string]struct{})
	for _, p := range progs {
		for c := range p.children {
			children[c] = struct{}{}
			p.children[c] = progs[c]
		}
	}
	var root *prog
	for n, p := range progs {
		if len(p.children) > 0 {
			if _, ok := children[n]; !ok {
				root = p
				break
			}
		}
	}
	fmt.Println(root.name)

	// For part 2, calculate each program's total weight.
	// Identify the child of the root with a different sum and start looking there.
	root.update()
	bad, diff := root.check()
	bad.find(diff)
}

type prog struct {
	name     string
	weight   int
	children map[string]*prog
	sum      int // weight of program plus all its children
}

// update recursively sets the sum field of p and all of its children.
func (p *prog) update() int {
	p.sum = p.weight
	for _, c := range p.children {
		p.sum += c.update()
	}
	return p.sum
}

// check returns p's child with a different summed weight from the other children.
func (p *prog) check() (bad *prog, diff int) {
	if len(p.children) < 3 {
		return nil, 0
	}

	sums := make([]int, 0, len(p.children))
	for _, c := range p.children {
		sums = append(sums, c.sum)
	}
	sort.Ints(sums)

	good := lib.If(sums[0] == sums[1], sums[0], sums[len(sums)-1])
	for _, c := range p.children {
		if c.sum != good {
			return c, c.sum - good
		}
	}
	lib.Panicf("Bad child not found")
	return nil, 0
}

// find recursively searches for the program under p whose weight needs to be adjusted
// by diff to balance the tree. The correct weight is printed.
func (p *prog) find(diff int) {
	sums := make([]int, 0, len(p.children))
	for _, c := range p.children {
		sums = append(sums, c.sum)
	}
	sort.Ints(sums)

	if len(sums) == 0 || sums[0] == sums[len(sums)-1] {
		fmt.Println(p.weight - diff)
	} else {
		bad, _ := p.check()
		bad.find(diff)
	}
}
