package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	in := lib.InputLines("2016/9")[0]

	cnt := expand(in, false)
	fmt.Println(cnt)

	cnt2 := expand(in, true)
	fmt.Println(cnt2)
}

var markerRegexp = regexp.MustCompile(`^\((\d+)x(\d+)\)`)

func expand(s string, recurse bool) (expanded int64) {
	for i := 0; i < len(s); {
		s := s[i:]
		if m := markerRegexp.FindString(s); m != "" {
			var length, times int
			lib.Extract(s, markerRegexp.String(), &length, &times)
			if recurse {
				expanded += int64(times) * expand(s[len(m):len(m)+length], true)
			} else {
				expanded += int64(times * length)
			}
			i += len(m) + length
		} else {
			expanded++
			i++
		}
	}
	return expanded
}
