package lib

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var nonDigitRegexp = regexp.MustCompile(`\D+`)

// ExtractInts extracts all positive integers from s.
// Non-digits are ignored.
func ExtractInts(s string) []int {
	v64s := ExtractInt64s(s)
	vals := make([]int, len(v64s))
	for i, v64 := range v64s {
		// TODO: Check for overflow.
		vals[i] = int(v64)
	}
	return vals
}

// ExtractInt64s extracts all positive integers from s as 64-bit ints.
// Non-digits are ignored.
func ExtractInt64s(s string) []int64 {
	var vals []int64
	for _, m := range nonDigitRegexp.Split(s, -1) {
		if m == "" {
			continue
		}
		v, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("Failed parsing %q as int64: %v", m, err))
		}
		vals = append(vals, v)
	}
	if len(vals) == 0 {
		panic("No ints found")
	}
	return vals
}

// ExtractDigits extracts individual digits from s and returns them as ints.
func ExtractDigits(s string) []int {
	var vals []int
	for _, ch := range []byte(s) {
		if ch >= '0' && ch <= '9' {
			vals = append(vals, int(ch-'0'))
		}
	}
	if len(vals) == 0 {
		panic("No digits found")
	}
	return vals
}

var extractCache = make(map[string]*regexp.Regexp)

// Extract executes regular expression re on s and assigns groups to dsts.
func Extract(s, re string, dsts ...interface{}) {
	cre, ok := extractCache[re]
	if !ok {
		cre = regexp.MustCompile(re)
		extractCache[re] = cre
	}

	ms := cre.FindStringSubmatch(s)
	if ms == nil {
		panic(fmt.Sprintf("%q not matched by %q", s, re))
	}
	if len(ms)-1 != len(dsts) {
		panic(fmt.Sprintf("%q has %v group(s), but %v dest(s) were supplied", re, len(ms)-1, len(dsts)))
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

// Tokenize splits s into tokens from the supplied args (either string or *regexp.Regexp).
// Whitespace is ignored. Regexps should be '^'-prefixed for better performance.
func Tokenize(s string, tokens ...interface{}) []string {
	var out []string
	for s = strings.TrimSpace(s); len(s) > 0; s = strings.TrimSpace(s) {
		found := false
		for _, tok := range tokens {
			var v string
			switch t := tok.(type) {
			case string:
				if strings.HasPrefix(s, t) {
					v = t
				}
			case *regexp.Regexp:
				if loc := t.FindStringIndex(s); loc != nil && loc[0] == 0 {
					v = s[loc[0]:loc[1]]
				}
			default:
				panic(fmt.Sprintf("Invalid token %q of type %T", tok, tok))
			}
			if len(v) > 0 {
				out = append(out, v)
				s = s[len(v):]
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Didn't find token at beginning of %q", s))
		}
	}
	return out
}
