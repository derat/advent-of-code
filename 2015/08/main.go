package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	hexRegexp := regexp.MustCompile(`^\\x[0-9a-f]{2}`)
	charRegexp := regexp.MustCompile(`^[^\\]`)
	var lit, mem, enc int
	for _, ln := range lib.InputLines("2015/8") {
		// Part 1
		var inner string
		lib.Extract(ln, `^"(.*)"$`, &inner)
		toks := lib.Tokenize(inner, `\\`, `\"`, hexRegexp, charRegexp)
		lit += len(ln)
		mem += len(toks)

		// Part 2
		enc += 2 // start/end quote
		for _, ch := range ln {
			if ch == '"' || ch == '\\' {
				enc += 2
			} else {
				enc++
			}
		}
	}
	fmt.Println(lit - mem)
	fmt.Println(enc - lit)
}
