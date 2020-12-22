package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var d1, d2 []int
	var recip *[]int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		switch ln {
		case "":
		case "Player 1:":
			recip = &d1
		case "Player 2:":
			recip = &d2
		default:
			v, err := strconv.Atoi(ln)
			if err != nil {
				log.Fatalf("Bad card %q: %v", ln, err)
			} else if v > 255 { // protect optimization in state()
				log.Fatalf("Card value %v > 255", v)
			}
			*recip = append(*recip, v)
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

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
