package main

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	polymer := lib.InputLines("2018/5")[0]

	// Part 1: Print length after performing all reactions.
	fmt.Println(collapse(polymer))

	// Part 2: Print shortest possible length by removing one type.
	min := math.MaxInt32
	for ch := rune('a'); ch <= rune('z'); ch++ {
		// This is inefficient, but I'm not sure if there's a simpler way
		// to remove multiple runes from a string without iterating.
		mod := strings.ReplaceAll(polymer, string(ch), "")
		mod = strings.ReplaceAll(mod, string(unicode.ToUpper(ch)), "")
		if n := collapse(mod); n < min {
			min = n
		}
	}
	fmt.Println(min)
}

func collapse(polymer string) int {
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
	return len(polymer)
}
