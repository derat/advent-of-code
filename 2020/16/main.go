package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2020/16")
	lib.AssertEq(len(pgs), 3)

	rules := make(map[string]rule) // keyed by field name
	for _, ln := range pgs[0] {
		var name string
		var min1, max1, min2, max2 int
		lib.Parse(ln, `^(.+): (\d+)-(\d+) or (\d+)-(\d+)$`,
			&name, &min1, &max1, &min2, &max2)
		rules[name] = rule{[][]int{{min1, max1}, {min2, max2}}}
	}

	lib.AssertEq(len(pgs[1]), 2)
	lib.AssertEq(pgs[1][0], "your ticket:")
	yours := ticket{lib.ExtractInts(pgs[1][1])}

	lib.AssertEq(pgs[2][0], "nearby tickets:")
	var nearby []ticket
	for _, ln := range pgs[2][1:] {
		nearby = append(nearby, ticket{lib.ExtractInts(ln)})
	}

	errRate := 0
	var valid []ticket
	for _, t := range nearby {
		tval := true
		for _, f := range t.fields {
			fval := false
			for _, r := range rules {
				if r.valid(f) {
					fval = true
					break
				}
			}
			if !fval {
				tval = false
				errRate += f
			}
		}
		if tval {
			valid = append(valid, t)
		}
	}
	fmt.Println(errRate)

	// Make a map from field name to possible indexes.
	poss := make(map[string]map[int]struct{}, len(rules))
	for n := range rules {
		fm := make(map[int]struct{})
		for i := range yours.fields {
			fm[i] = struct{}{}
		}
		poss[n] = fm
	}

	// Drop indexes with values that are out of range.
	for _, t := range valid {
		for i, f := range t.fields {
			for n, r := range rules {
				if !r.valid(f) {
					delete(poss[n], i)
				}
			}
		}
	}

	// Determine the mapping from field name to index.
	indexes := make(map[string]int, len(rules))
	for len(poss) != 0 {
		plen := len(poss)
		for n, p := range poss {
			if len(p) == 1 {
				var v int
				for v = range p {
				}
				indexes[n] = v
				for o := range poss {
					delete(poss[o], v)
				}
				delete(poss, n)
				break
			}
		}
		if len(poss) == plen {
			log.Fatal("didn't find an index")
		}
	}

	prod := 1
	for n, i := range indexes {
		if strings.HasPrefix(n, "departure") {
			prod *= yours.fields[i]
		}
	}
	fmt.Println(prod)
}

type rule struct {
	ranges [][]int
}

func (r *rule) valid(n int) bool {
	for _, ra := range r.ranges {
		if n >= ra[0] && n <= ra[1] {
			return true
		}
	}
	return false
}

type ticket struct {
	fields []int
}
