package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// Part 1
	closers := map[rune]struct {
		opener rune
		points int
	}{
		')': {'(', 3},
		']': {'[', 57},
		'}': {'{', 1197},
		'>': {'<', 25137},
	}

	// Part 2
	openers := map[rune]int{'(': 1, '[': 2, '{': 3, '<': 4}

	var score1 int
	var scores2 []int
Loop:
	for _, ln := range lib.InputLines("2021/10") {
		var open []rune
		for _, ch := range ln {
			switch ch {
			case '(', '[', '{', '<':
				open = append(open, ch)
			case ')', ']', '}', '>':
				// Part 1: Sum scores from syntax errors (first mismatched pairs).
				info := closers[ch]
				if len(open) > 0 && open[len(open)-1] != info.opener {
					score1 += info.points
					continue Loop
				} else {
					open = open[:len(open)-1]
				}
			default:
				log.Panicf("Invalid character %q", ch)
			}
		}

		// Part 2: Compute completion string's score for each line (for each missing
		// closing character, multiply completion score by 5 and add closing char's
		// points).
		lib.AssertLess(0, len(open)) // given in puzzle: all remaining lines are incomplete
		var compScore int
		for i := len(open) - 1; i >= 0; i-- {
			compScore = 5*compScore + openers[open[i]]
		}
		scores2 = append(scores2, compScore)
	}

	fmt.Println(score1)

	// Part 2: Print median of completion scores.
	sort.Ints(scores2)
	fmt.Println(scores2[len(scores2)/2])
}
