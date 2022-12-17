package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var seq []int // dx of jets of gas (input)
	for _, ch := range strings.TrimSpace(lib.Input("2022/17")) {
		switch ch {
		case '<':
			seq = append(seq, -1)
		case '>':
			seq = append(seq, 1)
		default:
			lib.Panicf("Invalid input char %q", ch)
		}
	}

	const (
		width  = 7
		rocks  = 2022          // part 1
		rocks2 = 1000000000000 // part 2
	)

	chamber := make(region)
	heights := make([]int, width)
	var maxHeight int

	// test returns true if the bottom-left corner of shape can be placed at pos
	// without overlapping with the walls, floor, or a stopped rock.
	test := func(shape region, pos point) bool {
		if pos.x < 0 || pos.y <= 0 {
			return false
		}
		for p := range shape {
			p.x += pos.x
			p.y += pos.y
			if p.x >= width || lib.MapHasKey(chamber, p) {
				return false
			}
		}
		return true
	}

	// isFull returns true if the specified row is completely filled by rocks.
	isFull := func(y int) bool {
		for x := 0; x < width; x++ {
			if !lib.MapHasKey(chamber, point{x, y}) {
				return false
			}
		}
		return true
	}

	var step int
	var startRock, startHeight, loopLen, loopHeight, endRock int
	for rock := 0; true; rock++ {
		shape := shapes[rock%len(shapes)]

		// "Each rock appears so that its left edge is two units away from the left wall and its
		// bottom edge is three units above the highest rock in the room (or the floor, if there
		// isn't one)."
		pos := point{2, maxHeight + 4}

		// "After a rock appears, it alternates between being pushed by a jet of hot gas one unit
		// (in the direction indicated by the next symbol in the jet pattern) and then falling one
		// unit down. If any movement would cause any part of the rock to move into the walls,
		// floor, or a stopped rock, the movement instead does not occur. If a downward movement
		// would have caused a falling rock to move into the floor or an already-fallen rock, the
		// falling rock stops where it is (having landed on something) and a new rock immediately
		// begins falling."
		for {
			p0 := point{pos.x + seq[step%len(seq)], pos.y}
			step++
			if test(shape, p0) {
				// Shape is moved by gas.
				pos = p0
			}

			p0 = point{pos.x, pos.y - 1}
			if test(shape, p0) {
				// Shape moves down one space.
				pos = p0
			} else {
				// Shape is stopped.
				for p := range shape {
					p.x += pos.x
					p.y += pos.y
					chamber[p] = r
					heights[p.x] = lib.Max(heights[p.x], p.y)
				}
				maxHeight = lib.Max(heights...)
				break
			}
		}

		if rock == rocks-1 {
			fmt.Println(maxHeight) // part 1
		}

		// TODO: This produces the correct answer with my input but it fails to find a loop in the
		// example input in a reasonable amount of time. It actually doesn't even find a single full
		// line, so maybe I was just lucky that this approach works. :-/
		next := rock + 1
		switch {
		case startRock == 0: // looking for start of loop
			if isFull(maxHeight) {
				startRock = next
				startHeight = maxHeight
			}
		case loopLen == 0: // looking for end of loop
			if next%len(shapes) == startRock%len(shapes) && isFull(maxHeight) {
				loopLen = next - startRock
				loopHeight = maxHeight - startHeight
				endRock = next + int((rocks2-int64(startRock))%int64(loopLen))
			}
		default: // looking for added post-loop height
			if next == endRock {
				loops := (rocks2-int64(startRock))/int64(loopLen) - 1
				fmt.Println(int64(maxHeight) + loops*int64(loopHeight))
				return
			}
		}
	}
}

type rock struct{}
type point struct{ x, y int }
type region map[point]rock

var r = rock{}
var shapes = []region{
	// ####
	{point{0, 0}: r, point{1, 0}: r, point{2, 0}: r, point{3, 0}: r},
	// .#.
	// ###
	// .#.
	{point{1, 2}: r, point{0, 1}: r, point{1, 1}: r, point{2, 1}: r, point{1, 0}: r},
	// ..#
	// ..#
	// ###
	{point{2, 2}: r, point{2, 1}: r, point{0, 0}: r, point{1, 0}: r, point{2, 0}: r},
	// #
	// #
	// #
	// #
	{point{0, 3}: r, point{0, 2}: r, point{0, 1}: r, point{0, 0}: r},
	// ##
	// ##
	{point{0, 1}: r, point{1, 1}: r, point{0, 0}: r, point{1, 0}: r},
}
