package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var cwd []string                 // path components of current dir (empty for root)
	dirSizes := make(map[string]int) // total sizes keyed by path-separated dir name

	// I was happy that I went with a simple non-recursive approach that happened to
	// make part 2 trivial.
	for _, ln := range lib.InputLines("2022/7") {
		var cd, fn string
		var sz int
		switch {
		case lib.TryExtract(ln, `^\$ cd\s+(.+)$`, &cd):
			switch cd {
			case "/":
				cwd = nil
			case "..":
				lib.Assert(len(cwd) > 0)
				cwd = cwd[:len(cwd)-1]
			default:
				cwd = append(cwd, cd)
			}
		case lib.TryExtract(ln, `^(\d+)\s+(.+)$`, &sz, &fn):
			// This will double-count files if ls is run twice in the same dir,
			// but this doesn't occur in my input.
			for i := 0; i <= len(cwd); i++ {
				dir := "/" + strings.Join(cwd[:i], "/")
				dirSizes[dir] += sz
			}
		}
	}

	// Part 1: Print sum of total sizes of all dirs with total sizes <= 100000.
	var sum int
	for _, sz := range dirSizes {
		if sz <= 100000 {
			sum += sz
		}
	}
	fmt.Println(sum)

	// Part 2: Print total size of smallest dir that, when deleted, will produce free
	// space of at least 30000000.
	const (
		full = 70000000
		req  = 30000000
	)
	used := dirSizes["/"]
	lib.AssertLessEq(used, full)
	free := full - used
	lib.AssertLess(free, req)
	need := req - free
	var best int
	for _, sz := range dirSizes {
		if sz >= need && (best == 0 || sz < best) {
			best = sz
		}
	}
	lib.Assert(best > 0)
	fmt.Println(best)
}
