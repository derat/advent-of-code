package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	diag := lib.InputLinesBytes("2017/19")

	// Returns true if the supplied position can be entered in the specified direction.
	valid := func(r, c int, vert bool) bool {
		if r < 0 || r >= len(diag) || c < 0 || c >= len(diag[0]) {
			return false
		}
		ch := diag[r][c]
		if ch == ' ' {
			return false
		}
		if vert {
			return ch != '-'
		}
		return ch != '|'
	}

	// Find the starting point at the top of the diagram.
	r, c := 0, -1
	for i, ch := range diag[0] {
		if ch == '|' {
			c = i
			break
		}
	}
	lib.Assert(c >= 0)

	dr, dc := 1, 0 // moving down at first
	var letters string
	var steps int
Loop:
	for {
		// This problem seems pretty loosely specified.
		// I'm not sure whether turns are always denoted with '+' or not.
		if diag[r][c] == '+' {
			if dr != 0 { // moving vertically
				switch {
				case valid(r, c-1, false): // left
					dr, dc = 0, -1
				case valid(r, c+1, false): // right
					dr, dc = 0, 1
				default:
					lib.Panicf("Can't move horizontally at %d, %d", r, c)
				}
			} else { // moving horizontally
				switch {
				case valid(r-1, c, true): // up
					dr, dc = -1, 0
				case valid(r+1, c, true): // down
					dr, dc = 1, 0
				default:
					lib.Panicf("Can't move vertically at %d, %d", r, c)
				}
			}
		}

		r, c = r+dr, c+dc
		steps++
		ch := diag[r][c]
		switch {
		case ch >= 'A' && ch <= 'Z':
			letters += string(ch)
		case ch == ' ':
			break Loop
		}
	}
	fmt.Println(letters)
	fmt.Println(steps)
}
