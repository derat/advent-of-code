package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var orig string
	reps := make(map[string][]string)
	for _, ln := range lib.InputLines("2015/19") {
		if strings.Contains(ln, "=>") {
			var src, dst string
			lib.Extract(ln, `^(\w+) => (\w+)$`, &src, &dst)
			reps[src] = append(reps[src], dst)
		} else {
			orig = ln
		}
	}

	mols := make(map[string]struct{})
	for i := 0; i < len(orig); i++ {
		for src, dsts := range reps {
			if strings.HasPrefix(orig[i:], src) {
				for _, dst := range dsts {
					mol := orig[:i] + dst + orig[i+len(src):]
					mols[mol] = struct{}{}
				}
			}
		}
	}
	fmt.Println(len(mols))
}
