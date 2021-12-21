package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

const (
	nspaces = 10
	nsides  = 100
	win     = 1000
	win2    = 21
)

func main() {
	input := lib.InputLines("2021/21")
	lib.AssertEq(len(input), 2)
	var start1, start2 int
	lib.Extract(input[0], `^Player 1 starting position: (\d+)$`, &start1)
	lib.Extract(input[1], `^Player 2 starting position: (\d+)$`, &start2)

	// Part 1: "Play a practice game using the deterministic 100-sided die. The moment either player
	// wins, what do you get if you multiply the score of the losing player by the number of times
	// the die was rolled during the game?"
	p1 := counter{1, nspaces, start1, 0}
	p2 := counter{1, nspaces, start2, 0}
	die := counter{1, nsides, 0, 0}
	var s1, s2 int
	for {
		if s1 += p1.inc(die.inc(1) + die.inc(1) + die.inc(1)); s1 >= win {
			break
		}
		if s2 += p2.inc(die.inc(1) + die.inc(1) + die.inc(1)); s2 >= win {
			break
		}
	}
	fmt.Println(lib.Min(s1, s2) * die.incs)

	// Part 2: "Using your given starting positions, determine every possible outcome. Find the
	// player that wins in more universes; in how many universes does that player win?"
	//
	// The state of a game can be represented by the following:
	//  - Turn             0-1, 1 bit
	//  - Player 1 space: 1-10, 4 bits
	//  - Player 2 space: 1-10, 4 bits
	//  - Player 1 score: 0-20, 5 bits
	//  - Player 2 score: 0-20, 5 bits
	//
	// Since I always err on the side of bit-packing the state into an uint64 prematurely,
	// I'm just going to throw everything into a struct to begin with.
	type state struct {
		turn  int8
		pos   [2]int8
		score [2]int8
	}

	// Find all possible outcomes of three rolls of the die.
	rolls := make([]int8, 0, 3*3*3)
	for i := 1; i <= 3; i++ {
		for j := 1; j <= 3; j++ {
			for k := 1; k <= 3; k++ {
				rolls = append(rolls, int8(i+j+k))
			}
		}
	}

	advance := func(start, amt int8) int8 { return 1 + (start+amt-1)%nspaces }

	cache := make(map[state][2]int) // state to wins
	var play func(state) [2]int
	play = func(s state) (wins [2]int) {
		if w, ok := cache[s]; ok {
			return w
		}

		for _, r := range rolls {
			npos := advance(s.pos[s.turn], r)
			nscore := s.score[s.turn] + npos
			if nscore >= win2 {
				wins[s.turn]++
				continue
			}

			ns := s
			ns.pos[s.turn] = npos
			ns.score[s.turn] = nscore
			ns.turn = (ns.turn + 1) % 2
			w := play(ns)
			wins[0] += w[0]
			wins[1] += w[1]
		}

		cache[s] = wins
		return wins
	}

	scores := play(state{turn: 0, pos: [2]int8{int8(start1), int8(start2)}})
	fmt.Println(lib.Max(scores[0], scores[1]))
}

// This type is sort of dumb, but the board and die share similar logic, so why not.
// (I wrote this for part 1, but it ended up not being useful for part 2.)
type counter struct{ min, max, v, incs int }

func (c *counter) inc(amt int) int {
	c.incs++
	c.v = c.min + ((c.v + amt - c.min) % (c.max - c.min + 1))
	return c.v
}
