package lib

import (
	"regexp"
	"strconv"
	"strings"
)

var intRegexp = regexp.MustCompile(`-?\d+`)
var uintRegexp = regexp.MustCompile(`\d+`)

// ExtractInts extracts all integers from s.
// Non-digits (besides '-') are ignored.
func ExtractInts(s string) []int {
	v64s := ExtractInt64s(s)
	vals := make([]int, len(v64s))
	for i, v64 := range v64s {
		// TODO: Check for overflow.
		vals[i] = int(v64)
	}
	return vals
}

// ExtractInt64s extracts all integers from s as 64-bit ints.
// Non-digits (besides '-') are ignored.
func ExtractInt64s(s string) []int64 {
	var vals []int64
	for _, s := range intRegexp.FindAllString(s, -1) {
		v, err := strconv.ParseInt(s, 10, 64)
		Assertf(err == nil, "Failed parsing %q as int64: %v", s, err)
		vals = append(vals, v)
	}
	Assertf(len(vals) > 0, "No ints found")
	return vals
}

// ExtractUints extracts all zero or positive integers from s as ints.
// Non-digits (including '-') are ignored.
func ExtractUints(s string) []int {
	var vals []int
	for _, s := range uintRegexp.FindAllString(s, -1) {
		v, err := strconv.Atoi(s)
		Assertf(err == nil, "Failed parsing %q as int64: %v", s, err)
		vals = append(vals, v)
	}
	Assertf(len(vals) > 0, "No ints found")
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
	Assertf(len(vals) > 0, "No digits found")
	return vals
}

var regexpCache = make(map[string]*regexp.Regexp)

func getRegexp(re string) *regexp.Regexp {
	comp, ok := regexpCache[re]
	if !ok {
		comp = regexp.MustCompile(re)
		regexpCache[re] = comp
	}
	return comp
}

// ExtractMaybe executes regular expression re on s and assigns groups to dsts.
// It returns false if re does not match s.
func ExtractMaybe(s, re string, dsts ...interface{}) bool {
	ms := getRegexp(re).FindStringSubmatch(s)
	if ms == nil {
		return false
	}
	Assertf(len(ms)-1 == len(dsts), "%q has %v group(s), but %v dest(s) were supplied", re, len(ms)-1, len(dsts))
	for i, dst := range dsts {
		m := ms[i+1]
		if m == "" { // skip optional groups (should maybe set dst to zero?)
			continue
		}
		var err error
		switch t := dst.(type) {
		case *byte:
			Assertf(len(m) == 1, "Can't store %q in a single byte", m)
			*t = m[0]
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
			Assertf(false, "Unknown dest type %T for group %v of %q", t, i, re)
		}
		Assertf(err == nil, "Failed to parse group %q matched by %q: %v", m, re, err)
	}
	return true
}

// Extract executes regular expression re on s and assigns groups to dsts.
// It panics if re does not match s.
func Extract(s, re string, dsts ...interface{}) {
	Assertf(ExtractMaybe(s, re, dsts...), "%q not matched by %q", s, re)
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
				Assertf(false, "Invalid token %q of type %T", tok, tok)
			}
			if len(v) > 0 {
				out = append(out, v)
				s = s[len(v):]
				found = true
				break
			}
		}
		Assertf(found, "Didn't find token at beginning of %q", s)
	}
	return out
}
