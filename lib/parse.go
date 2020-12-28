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

var parseCache = make(map[string]*regexp.Regexp)

// Parse executes regular expression re on s and assigns groups to dsts.
func Parse(s, re string, dsts ...interface{}) {
	cre, ok := parseCache[re]
	if !ok {
		cre = regexp.MustCompile(re)
		parseCache[re] = cre
	}

	ms := cre.FindStringSubmatch(s)
	if ms == nil {
		panic(fmt.Sprintf("%q not matched by %q", s, re))
	}
	if len(ms)-1 != len(dsts) {
		panic(fmt.Sprintf("%q has %v group(s), but %v dests were supplied", re, len(ms)-1, len(dsts)))
	}
	for i, dst := range dsts {
		m := ms[i+1]
		if m == "" { // skip optional groups (should maybe set dst to zero?)
			continue
		}
		var err error
		switch t := dst.(type) {
		case *float64:
			*t, err = strconv.ParseFloat(m, 64)
		case *int:
			*t, err = strconv.Atoi(m)
		case *int64:
			*t, err = strconv.ParseInt(m, 10, 64)
		case *string:
			*t = m
		case *uint64:
			*t, err = strconv.ParseUint(m, 10, 64)
		default:
			panic(fmt.Sprintf("Unknown dest type %T for group %v of %q", t, i, re))
		}
		if err != nil {
			panic(fmt.Sprintf("Failed to parse group %q matched by %q: %v", m, re, err))
		}
	}
}
