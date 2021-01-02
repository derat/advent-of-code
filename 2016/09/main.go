package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	in := lib.InputLines("2016/9")[0]
	markerRegexp := regexp.MustCompile(`^\((\d+)x(\d+)\)`)
	var i, cnt int
	for i < len(in) {
		s := in[i:]
		if m := markerRegexp.FindString(s); m != "" {
			var length, times int
			lib.Extract(s, markerRegexp.String(), &length, &times)
			cnt += length * times
			i += len(m) + length
		} else {
			cnt++
			i++
		}
	}
	fmt.Println(cnt)
}
