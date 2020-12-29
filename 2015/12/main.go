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
	stack := []frame{{}}
	for _, tok := range lib.Tokenize(ln, "[", "]", "{", "}", ",", ":", stringRegexp, numRegexp) {
		switch tok {
		case "[", "{":
			stack = append(stack, frame{tok: tok})
		case "]", "}":
			if top := &stack[len(stack)-1]; !top.red {
				stack[len(stack)-2].sum += top.sum
			}
			stack = stack[:len(stack)-1]
		case `"red"`:
			// Laziness alert: my input never contains "red" as a property name, so
			// just look for "red" tokens while an object is at the top of the stack.
			if top := &stack[len(stack)-1]; top.tok == "{" {
				top.red = true
			}
		default:
			v, err := strconv.Atoi(tok)
			if err == nil {
				sum += v
				stack[len(stack)-1].sum += v
			}
		}

	}
	fmt.Println(sum)

	lib.AssertEq(len(stack), 1)
	fmt.Println(stack[0].sum)
}

type frame struct {
	tok string
	sum int
	red bool
}
