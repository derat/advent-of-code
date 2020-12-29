package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	ln := lib.InputLines("2015/12")[0]

	// I've already written too much code in my life that uses existing
	// libraries to parse JSON and refuse to write more here.
	var sum int
	stringRegexp := regexp.MustCompile(`^"[^"]*"`)
	numRegexp := regexp.MustCompile(`^-?\d+`)
	for _, tok := range lib.Tokenize(ln, "[", "]", "{", "}", ",", ":", stringRegexp, numRegexp) {
		v, err := strconv.Atoi(tok)
		if err == nil {
			sum += v
		}
	}
	fmt.Println(sum)
}
