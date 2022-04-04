package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var target string
	repls := make(map[string][]string)
	for _, ln := range lib.InputLines("2015/19") {
		if strings.Contains(ln, "=>") {
			var src, dst string
			lib.Extract(ln, `^(\w+) => (\w+)$`, &src, &dst)
			repls[src] = append(repls[src], dst)
		} else {
			target = ln
		}
	}

	// Part 1: Count single replacements
	seen := make(map[string]struct{})
	for i := 0; i < len(target); i++ {
		for src, dsts := range repls {
			if !strings.HasPrefix(target[i:], src) {
				continue
			}
			for _, dst := range dsts {
				mol := target[:i] + dst + target[i+len(src):]
				seen[mol] = struct{}{}
			}
		}
	}
	fmt.Println(len(seen))

	// Part 2: Minimum replacements to go from 'e' to target
	//
	// One observation: replacements seem to only go from shorter to longer strings (or equal
	// length, but with more elements), so we shouldn't need to worry about infinite recursion as
	// long as we terminate once we reach the target length. We still need to do something smart
	// around pruning recursion, though, because otherwise we end up expanding Ca => CaCa over and
	// over until we hit the target length.
	//
	// Breadth-first search doesn't seem to help here: there's a huge number of potential strings
	// just by the ninth iteration.
	//
	// The below code works in the opposite direction, starting with the target string and applying
	// reversed replacements to try to get to 'e'. It's still not pruning adequately: I saw it
	// reach an answer once that the website confirmed is the minimum, but it's not fully exploring
	// the problem space, and it's also nondeterministic (since it iterates over a map) so it hasn't
	// even produced any answers on other runs.
	//
	// Update after reading discussion at https://redd.it/3xflz8: It sounds like the "right" way
	// to solve this one was by looking at the specific input instead of writing generalized code.
	// I'm changing this code to return the first solution it finds and calling it a day.

	backRepls := make(map[string]string)
	backRegexps := make(map[string]*regexp.Regexp)
	for src, dsts := range repls {
		for _, dst := range dsts {
			lib.AssertEq(backRepls[dst], "")
			backRepls[dst] = src
			backRegexps[dst] = regexp.MustCompile(dst)
		}
	}

	// Iterating over a map yields (intentionally) non-deterministic behavior.
	// Try longest strings first to avoid this.
	backSorted := lib.MapKeys(backRepls)
	sort.Slice(backSorted, func(i, j int) bool {
		si, sj := backSorted[i], backSorted[j]
		return len(si) > len(sj) || (len(si) == len(sj) && si < sj)
	})

	seen = make(map[string]struct{})
	var recurse func(string, int) int
	recurse = func(mol string, steps int) int {
		if mol == "e" {
			return steps
		}

		steps++
		for _, src := range backSorted {
			dst := backRepls[src]
			for _, idxs := range backRegexps[src].FindAllStringIndex(mol, -1) {
				newMol := mol[:idxs[0]] + dst + mol[idxs[1]:]
				if _, ok := seen[newMol]; !ok {
					seen[mol] = struct{}{}
					if s := recurse(newMol, steps); s != -1 {
						return s
					}
				}
			}
		}
		return -1
	}
	fmt.Println(recurse(target, 0))
}
