package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	in := lib.InputInts("2021/17")
	txmin, txmax, tymin, tymax := in[0], in[1], in[2], in[3]

	// Part 1: "Find the initial velocity that causes the probe to reach the highest y position and
	// still eventually be within the target area after any step. What is the highest y position it
	// reaches on this trajectory?"
	//
	// The y position increases by [n, n-1, n-2, ...].
	// During the step that brings the y position back to 0, the y velocity will be -n. So, if we're
	// looking for the initial velocity that brings the probe to the highest y position, we want the
	// next step to bring the probe to the bottom of the target area.
	lib.AssertLess(tymin, 0) // assume target is below starting point
	vymax := -tymin - 1

	// To find the highest y position, we need to sum [n, n-1, n-2, ..., 1], i.e. n+n*(n-1)/2.
	//
	// This approach works for my input (and maybe all official inputs?), but there are some
	// inputs for which it won't work: https://www.reddit.com/r/adventofcode/comments/rid0g3
	// Specifically, it's wrong if there are no inputs such that the probe will eventually settle in
	// the x target range. I've switched to an approach that just uses the data from part 2.
	//fmt.Println(vymax + vymax*(vymax-1)/2)

	// Part 2: "How many distinct initial velocity values cause the probe to be within the target
	// area after any step?"
	//
	// The x position increases by [n, n-1, n-2, ..., n-n].
	// So if we're in the target area when the x velocity has decreased to 0, we'll remain there.
	//
	// Now that we're counting distinct velocities, we need to consider the x and y positions in
	// tandem, which means we care about the number of steps.
	//
	// I think that the way to approach this is to find all initial y velocities that get the probe
	// into the target area and then multiply each of those by the number of initial x velocities
	// that are in the target area after that many steps.
	//
	// These y velocities include ones that reach the target one step after y=0 (either because
	// they went up earlier and are coming down now or because they were initially negative), but
	// also ones that take several steps after y=0 to reach the target. There's probably some faster
	// way to find these, but just iterating over the entire possible range seems fine:
	yvels := make(map[int][]int) // y velocity -> steps during which y position is within target
	var smax int
	for vystart := tymin; vystart <= vymax; vystart++ {
		y := 0
		vy := vystart
		for steps := 1; y >= tymin; steps++ {
			y += vy
			vy--
			if y >= tymin && y <= tymax {
				yvels[vystart] = append(yvels[vystart], steps)
				smax = lib.Max(smax, steps)
			}
		}
	}

	lib.AssertLess(0, txmin)      // assume target is to right of starting point
	xsteps := make(map[int][]int) // steps -> x velocities
	for vxstart := 1; vxstart <= txmax; vxstart++ {
		x := 0
		vx := vxstart
		for steps := 1; x <= txmax && steps <= smax; steps++ {
			x += vx
			vx = lib.Max(vx-1, 0)
			if x >= txmin && x <= txmax {
				xsteps[steps] = append(xsteps[steps], vxstart)
			}
		}
	}

	// This ended up being a bit more subtle than I initially thought. Some trajectories will put
	// the probe in the target area for multiple steps, and we need to make sure that we only count
	// each of those once.
	vymax = 0 // part 1
	vels := make(map[[2]int]struct{})
	for y, steps := range yvels {
		for _, s := range steps {
			for _, x := range xsteps[s] {
				vels[[2]int{x, y}] = struct{}{}
				vymax = lib.Max(vymax, y) // part 1
			}
		}
	}
	fmt.Println(vymax + vymax*(vymax-1)/2) // part 1
	fmt.Println(len(vels))                 // part 2
}
