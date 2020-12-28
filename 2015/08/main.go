package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	hexRegexp := regexp.MustCompile(`^\\x[0-9a-f]{2}`)
	charRegexp := regexp.MustCompile(`^[^\\]`)
	var lit, mem int
	for _, ln := range lib.InputLines("2015/8") {
		var inner string
		lib.Extract(ln, `^"(.*)"$`, &inner)
		toks := lib.Tokenize(inner, `\\`, `\"`, hexRegexp, charRegexp)
		lit += len(ln)
		mem += len(toks)
	}
	fmt.Println(lit - mem)
}
