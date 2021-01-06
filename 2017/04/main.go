package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var nvalid int
Loop:
	for _, ln := range lib.InputLines("2017/4") {
		seen := make(map[string]struct{})
		for _, s := range strings.Fields(ln) {
			if _, ok := seen[s]; ok {
				continue Loop
			}
			seen[s] = struct{}{}
		}
		nvalid++
	}
	fmt.Println(nvalid)
}
