package lib

import (
	"fmt"
	"regexp"
	"strconv"
)

var nonDigitRegexp = regexp.MustCompile(`\D+`)

// ExtractInts extracts all positive integers from s.
// Non-digits are ignored.
func ExtractInts(s string) []int {
	var vals []int
	for _, m := range nonDigitRegexp.Split(s, -1) {
		if m == "" {
			continue
		}
		v, err := strconv.Atoi(m)
		if err != nil {
			panic(fmt.Sprintf("Failed parsing %q as int: %v", m, err))
		}
		vals = append(vals, v)
	}
	if len(vals) == 0 {
		panic("No ints found")
	}
	return vals
}
