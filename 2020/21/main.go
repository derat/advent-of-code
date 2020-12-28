package main

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

type allerg string
type allergMap map[allerg]struct{}

type ingred string
type ingredMap map[ingred]struct{}

var set = struct{}{}

type food struct {
	ingreds ingredMap
	allergs allergMap
}

func main() {
	var foods []food

	for _, ln := range lib.InputLines("2020/21") {
		fd := food{make(ingredMap), make(allergMap)}
		var left, right string
		lib.Extract(ln, `^(.+) \(contains (.+)\)$`, &left, &right)

		for _, f := range strings.Fields(left) {
			fd.ingreds[ingred(f)] = set
		}
		if len(fd.ingreds) == 0 {
			log.Fatalf("No ingredients in %q", ln)
		}

		for _, f := range strings.Split(right, ", ") {
			fd.allergs[allerg(f)] = set
		}
		if len(fd.allergs) == 0 {
			log.Fatalf("No allergens in %q", ln)
		}

		foods = append(foods, fd)
	}

	rules := make(map[allerg][]ingredMap) // values are ingredient lists
	known := make(map[ingred]allerg)
	for _, fd := range foods {
		for al := range fd.allergs {
			im := make(ingredMap, len(fd.ingreds))
			for in := range fd.ingreds {
				im[in] = set
			}
			rules[al] = append(rules[al], im)
		}
	}

	for len(rules) > 0 {
		for al, lists := range rules {
			// Eliminate any ingredients that aren't present in all lists.
			inCnt := make(map[ingred]int)
			for _, ins := range lists {
				for in := range ins {
					inCnt[in]++
				}
			}
			for in, cnt := range inCnt {
				if cnt != len(lists) {
					for _, ins := range lists {
						delete(ins, in)
					}
				}
			}

			// Check if there's a list with a single possible ingredient.
			for _, ins := range lists {
				if len(ins) == 1 {
					var in ingred // getting a single key is painful :-(
					for i := range ins {
						in = i
						break
					}

					//fmt.Println(in, "contains", al)
					delete(rules, al)
					known[in] = al

					// Delete the ingredient from all other allergens' lists.
					for _, lists := range rules {
						for _, ins := range lists {
							delete(ins, in)
						}
					}
					break
				}
			}
		}
	}

	// Part 1:
	var cnt int
	for _, fd := range foods {
		for in := range fd.ingreds {
			if _, ok := known[in]; !ok {
				cnt++
			}
		}
	}
	fmt.Println(cnt)

	// Part 2:
	var danger []string
	for in := range known {
		danger = append(danger, string(in)) // sigh, go
	}
	sort.Slice(danger, func(i, j int) bool {
		return known[ingred(danger[i])] < known[ingred(danger[j])]
	})
	fmt.Println(strings.Join([]string(danger), ","))
}
