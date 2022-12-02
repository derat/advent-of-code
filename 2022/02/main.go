package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Part 1: played shapes.
	themChars := map[byte]shape{
		'A': rock,
		'B': paper,
		'C': scissors,
	}
	meChars := map[byte]shape{
		'X': rock,
		'Y': paper,
		'Z': scissors,
	}

	// Part 2: desired outcomes.
	outChars := map[byte]outcome{
		'X': lose,
		'Y': draw,
		'Z': win,
	}

	var total, total2 int
	for _, ln := range lib.InputLines("2022/2") {
		var tc, mc byte
		lib.Extract(ln, `^([ABC]) ([XYZ])$`, &tc, &mc)

		them := themChars[tc]
		me := meChars[mc]
		total += result(me, them)

		out := outChars[mc]
		me2 := choose(them, out)
		total2 += result(me2, them)
	}
	fmt.Println(total)
	fmt.Println(total2)
}

type shape int // values are shape scores

const (
	rock     shape = 1
	paper    shape = 2
	scissors shape = 3
)

var beats = map[shape]shape{
	rock:     scissors,
	paper:    rock,
	scissors: paper,
}
var loses = lib.InvertMap(beats)

// result returns the score for a round in which the supplied shapes are played.
func result(me, them shape) int {
	switch {
	case beats[me] == them:
		return 6 + int(me) // win
	case loses[me] == them:
		return 0 + int(me) // loss
	default:
		return 3 + int(me) // draw
	}
}

type outcome int

const (
	lose outcome = iota
	draw
	win
)

// choose returns the shape that should be played to get the supplied outcome.
func choose(them shape, out outcome) shape {
	switch out {
	case lose:
		return beats[them]
	case draw:
		return them
	case win:
		return loses[them]
	}
	panic("not reached")
}
