package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/7")

	// Part 1: Find maximum output signal.
	fmt.Println(run(input, []int64{0, 1, 2, 3, 4}, false))

	// Part 2: Feed e's output into a's input.
	fmt.Println(run(input, []int64{5, 6, 7, 8, 9}, true))
}

func run(input, vals []int64, feedback bool) int64 {
	pch := make(chan []int64)
	go lib.Perms(vals, pch)

	// The problem just calls for running each amplifier in parallel
	// and manually copying values, but this is Go, so why not use
	// goroutines and channels?
	var max int64
	for phases := range pch {
		// Create the amplifiers and wire them up.
		amps := make([]*lib.Intcode, len(phases))
		for i := range amps {
			amps[i] = lib.NewIntcode(input)
			if i > 0 {
				amps[i].In = amps[i-1].Out
			}
		}
		if feedback {
			amps[0].In = amps[len(amps)-1].Out
		}

		// Start the amplifiers and feed them their phase signals.
		for i, a := range amps {
			a.Start()
			a.In <- phases[i]
		}

		// Send the input signal to the first and read the output from the last.
		amps[0].In <- 0
		lib.Assert(amps[0].Wait())
		if out := <-amps[len(amps)-1].Out; out > max {
			max = out
		}
	}
	return max
}
