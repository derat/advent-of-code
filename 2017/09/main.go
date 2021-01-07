package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var depth int // brackets
	var score int // bracket score for part 1
	var garbage, escape bool
	for i, ch := range lib.InputLines("2017/9")[0] {
		switch {
		case escape:
			escape = false // skip char after '!'
		case ch == '!':
			lib.Assertf(garbage, "'!' outside garbage at %d", i)
			escape = true
		case ch == '>':
			lib.Assertf(garbage, "'>' outside garbage at %d", i)
			garbage = false
		case garbage: // consume chars within garbage
		case ch == '<':
			garbage = true
		case ch == '{':
			depth++
			score += depth
		case ch == '}':
			depth--
		}
	}
	lib.AssertEq(depth, 0)
	fmt.Println(score)
}
