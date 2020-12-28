package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.ReadParagraphs()
	lib.AssertEq(len(pgs), 2)
	lib.AssertEq(pgs[0][0], "Player 1:")
	d1 := lib.ExtractInts(strings.Join(pgs[0][1:], " "))
	lib.AssertEq(pgs[1][0], "Player 2:")
	d2 := lib.ExtractInts(strings.Join(pgs[1][1:], " "))
	// Technically ought to check that no cards have values above 255.

	// Part 1:
	if p1win, d1, d2 := play(d1, d2, false /* recurse */); p1win {
		fmt.Println(score(d1))
	} else {
		fmt.Println(score(d2))
	}

	// Part 2:
	if p1win, d1, d2 := play(d1, d2, true /* recurse */); p1win {
		fmt.Println(score(d1))
	} else {
		fmt.Println(score(d2))
	}
}

// Returns true if player 1 wins.
func play(d1, d2 []int, recurse bool) (bool, []int, []int) {
	seen := make(map[string]struct{})

	for len(d1) > 0 && len(d2) > 0 {
		// If the decks were seen already, player 1 instantly wins the game.
		st := state(d1, d2)
		if _, ok := seen[st]; ok && recurse {
			return true, d1, d2
		}
		seen[st] = struct{}{}

		// Each player draws their top card.
		v1, v2 := d1[0], d2[0]
		d1, d2 = d1[1:], d2[1:]

		win1 := false
		switch {
		case recurse && len(d1) >= v1 && len(d2) >= v2:
			// Both players have at least as many cards remaining as the values they just drew. Recurse.
			win1, _, _ = play(append([]int{}, d1[:v1]...), append([]int{}, d2[:v2]...), true)
		default:
			// One or both players don't have enough cards left to recurse.
			win1 = v1 > v2
		}
		if win1 {
			d1 = append(d1, v1, v2)
		} else {
			d2 = append(d2, v2, v1)
		}
	}

	return len(d1) > 0, d1, d2
}

// Returns a (non-printable) string uniquely representing the state of the two decks.
func state(d1, d2 []int) string {
	// This could actually go further and pack four cards into each byte.
	b := make([]byte, 0, len(d1)+len(d2)+1)
	for _, v := range d1 {
		b = append(b, byte(v))
	}
	b = append(b, 0)
	for _, v := range d2 {
		b = append(b, byte(v))
	}
	return string(b)
}

func score(d []int) int {
	score := 0
	for i, v := range d {
		score += v * (len(d) - i)
	}
	return score
}
