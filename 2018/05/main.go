package main

import (
	"fmt"
	"unicode"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	polymer := lib.InputLines("2018/5")[0]
	var i int
	for i < len(polymer)-1 {
		a, b := rune(polymer[i]), rune(polymer[i+1])
		low := unicode.IsLower(a)
		if (low && b == unicode.ToUpper(a)) || (!low && b == unicode.ToLower(a)) {
			polymer = polymer[:i] + polymer[i+2:]
			if i > 0 {
				i--
			}
		} else {
			i++
		}
	}
	fmt.Println(len(polymer))
}
