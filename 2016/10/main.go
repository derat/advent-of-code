package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

type instr struct {
	low, high       int
	lowOut, highOut bool // true if going to output instead of bot
}

func main() {
	bots := make(map[int][]int)
	ins := make(map[int]instr)
	outs := make(map[int]int)
	var ready []int // ready bots

	for _, ln := range lib.InputLines("2016/10") {
		switch {
		case strings.Contains(ln, "goes to"):
			var bot, val int
			lib.Extract(ln, `^value (\d+) goes to bot (\d+)$`, &val, &bot)
			bots[bot] = append(bots[bot], val)
			if len(bots[bot]) == 2 {
				ready = append(ready, bot)
			}
		case strings.Contains(ln, "gives"):
			var bot, id1, id2 int
			var typ1, typ2 string
			lib.Extract(ln, `^bot (\d+) gives low to (bot|output) (\d+) and high to (bot|output) (\d+)$`,
				&bot, &typ1, &id1, &typ2, &id2)
			_, ok := ins[bot]
			lib.Assertf(!ok, "Already have instruction for bot %d", bot)
			ins[bot] = instr{id1, id2, typ1 == "output", typ2 == "output"}
		}
	}

	for len(ready) > 0 {
		var newReady []int
		for _, bot := range ready {
			vals := bots[bot]
			lib.AssertEq(len(vals), 2)

			// Part 1: Print number of bot responsible for comparing 17 and 61.
			if lib.Min(vals...) == 17 && lib.Max(vals...) == 61 {
				fmt.Println(bot)
			}

			in, ok := ins[bot]
			lib.Assertf(ok, "No instruction for bot %d", bot)
			if in.lowOut {
				outs[in.low] = lib.Min(vals...)
			} else {
				bots[in.low] = append(bots[in.low], lib.Min(vals...))
				if len(bots[in.low]) == 2 {
					newReady = append(newReady, in.low)
				}
			}
			if in.highOut {
				outs[in.high] = lib.Max(vals...)
			} else {
				bots[in.high] = append(bots[in.high], lib.Max(vals...))
				if len(bots[in.high]) == 2 {
					newReady = append(newReady, in.high)
				}
			}
			delete(bots, bot)
		}
		ready = newReady
	}

	// Part 2
	fmt.Println(outs[0] * outs[1] * outs[2])
}
