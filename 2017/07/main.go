package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	progs := make(map[string]*prog)
	for _, ln := range lib.InputLines("2017/7") {
		var name, rest string
		p := &prog{children: make(map[string]*prog)}
		lib.Extract(ln, `^([a-z]+) \((\d+)\)(.*)`, &name, &p.weight, &rest)

		if rest != "" {
			var list string
			lib.Extract(rest, `^ -> (.+)$`, &list)
			for _, s := range strings.Split(list, ", ") {
				p.children[s] = nil
			}
		}

		progs[name] = p
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
	for n, p := range progs {
		if len(p.children) > 0 {
			if _, ok := children[n]; !ok {
				fmt.Println(n)
				break
			}
		}
	}
}

type prog struct {
	weight   int
	children map[string]*prog
}
