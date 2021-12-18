package main

import (
	"fmt"
	"regexp"
	"strconv"
	"unicode"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lines := lib.InputLines("2021/18")

	// I initially parsed numbers into a recursive struct, but when I started implementing the
	// reduction process, I realized that it'd probably be way easier to just work with flat strings
	// (particularly for the explode operation). I was worried that part 2 would add a twist that
	// would make this approach no longer work, but it ended up being fine -- 1.3s to check all
	// combinations with my input.

	// Part 1: "Add up all of the snailfish numbers from the homework assignment in the order they
	// appear. What is the magnitude of the final sum?"
	sum := lines[0]
	for _, ln := range lines[1:] {
		sum = add(sum, ln)
	}
	mag, n := magnitude(sum)
	lib.AssertEq(n, len(sum))
	fmt.Println(mag)

	// Part 2: "What is the largest magnitude of any sum of two different snailfish numbers from the
	// homework assignment?"
	var max int
	for i := range lines {
		for j := range lines {
			if i == j {
				continue
			}
			sum := add(lines[i], lines[j])
			mag, n := magnitude(sum)
			lib.AssertEq(n, len(sum))
			max = lib.Max(mag, max)
		}
	}
	fmt.Println(max)
}

// add adds the supplied snailfish numbers and returns the fully-reduced sum.
func add(left, right string) string {
	s := fmt.Sprintf("[%s,%s]", left, right)
	for {
		var reduced bool
		if s, reduced = reduceOnce(s); !reduced {
			break
		}
	}
	return s
}

// reduceOnce reduces the supplied snailfish number once.
// It returns the resulting number and a bool indicating whether reduction was performed.
// If it's true, further reductions may be possible.
func reduceOnce(s string) (string, bool) {
	// "To reduce a snailfish number, you must repeatedly do the first action in this list that
	// applies to the snailfish number:"

	// "If any pair is nested inside four pairs, the leftmost such pair explodes."
	var depth int
	for i, ch := range s {
		switch {
		case ch == '[':
			// "To explode a pair, the pair's left value is added to the first regular number to the
			// left of the exploding pair (if any), and the pair's right value is added to the first
			// regular number to the right of the exploding pair (if any). Exploding pairs will
			// always consist of two regular numbers. Then, the entire exploding pair is replaced
			// with the regular number 0."
			if depth++; depth == 5 {
				var left, right int
				n := lib.Extract(s[i:], `^\[(\d+),(\d+)\]`, &left, &right)
				return explodeAdd(s[0:i], lastRegexp, left) + "0" +
					explodeAdd(s[i+n:], firstRegexp, right), true
			}
		case ch == ']':
			depth--
		}
	}

	// "If any regular number is 10 or greater, the leftmost such regular number splits."
	for i, ch := range s {
		if !unicode.IsDigit(rune(ch)) {
			continue
		}

		var val int
		n := lib.Extract(s[i:], `^(\d+)`, &val)
		if val <= 9 {
			continue
		}

		// "To split a regular number, replace it with a pair; the left element of the pair should
		// be the regular number divided by two and rounded down, while the right element of the
		// pair should be the regular number divided by two and rounded up."
		return s[0:i] + fmt.Sprintf("[%d,%d]", val/2, (val+1)/2) + s[i+n:], true
	}

	return s, false
}

var (
	firstRegexp = regexp.MustCompile(`\d+`)         // beginning of first int
	lastRegexp  = regexp.MustCompile(`\d+[^0-9]*$`) // beginning of last int
)

// explodeAdd uses re to find the beginning of an integer in s and adds val to it.
// If re doesn't match anything, no actions are performed.
func explodeAdd(s string, re *regexp.Regexp, val int) string {
	loc := re.FindStringIndex(s)
	if loc == nil {
		return s
	}
	i := loc[0]
	var old int
	n := lib.Extract(s[i:], `^(\d+)`, &old)
	return s[:i] + strconv.Itoa(old+val) + s[i+n:]
}

// magnitude computes the magnitude of reduced snailfish number s.
// It also returns the number of characters that were consumed.
func magnitude(s string) (val, n int) {
	switch {
	case unicode.IsDigit(rune(s[0])):
		// "The magnitude of a regular number is just that number."
		// We should just have single-digit numbers by this point.
		lib.Assert(len(s) == 1 || !unicode.IsDigit(rune(s[1])))
		return int(s[0] - '0'), 1

	case s[0] == '[':
		// "The magnitude of a pair is 3 times the magnitude of its
		// left element plus 2 times the magnitude of its right element."
		left, ln := magnitude(s[1:])
		lib.AssertEq(s[1+ln], byte(','))
		right, rn := magnitude(s[1+ln+1:])
		lib.AssertEq(s[1+ln+1+rn], byte(']'))
		return 3*left + 2*right, 1 + ln + 1 + rn + 1

	default:
		lib.Panicf("Bad char %q", s[0])
		return 0, 0
	}
}
