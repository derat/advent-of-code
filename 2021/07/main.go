package main

import (
	"fmt"
	"math"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	initial := lib.InputInts("2021/7")
	sort.Ints(initial)
	max := initial[len(initial)-1]

	// Part 1: Print cost of position with minimum cost.
	fmt.Println(slowCost(initial, initial[len(initial)/2]))

	// Part 2: Each step costs 1 more unit of fuel than the last.
	// Compute the cost of moving all crabs from the left of each position
	// to that position (index in the array is position).
	left := make([]int, max+1)
	costs := make([]int, 0, max+1) // last cost of crabs that we've picked up so far
	var j int                      // index into initial
	for i := range left {
		for ; j < len(initial) && initial[j] < i; j++ {
			costs = append(costs, 0)
		}
		if i > 0 {
			left[i] += left[i-1]
		}
		for k := range costs {
			costs[k]++
			left[i] += costs[k]
		}
	}

	// Do the same thing from the right end.
	right := make([]int, max+1)
	costs = costs[:0]
	j = len(initial) - 1
	for i := len(right) - 1; i >= 0; i-- {
		for ; j >= 0 && initial[j] > i; j-- {
			costs = append(costs, 0)
		}
		if i < len(right)-1 {
			right[i] += right[i+1]
		}
		for k := range costs {
			costs[k]++
			right[i] += costs[k]
		}
	}

	// Sum the costs of moving the crabs from the left and right of each position.
	min := math.MaxInt32
	for i := range left {
		min = lib.Min(min, left[i]+right[i])
	}
	fmt.Println(min)

	// Sigh, much simpler approach: the cost of moving from p to t is the summation
	// from 1 to the absolute value of distance d, which is just d*(d+1)/2.
	min2 := math.MaxInt32
	for t := 0; t <= max; t++ {
		var cost int
		for _, p := range initial {
			d := lib.Abs(p - t)
			cost += d * (d + 1) / 2
		}
		min2 = lib.Min(min2, cost)
	}
	lib.AssertEq(min2, min)
}

func slowCost(initial []int, target int) int {
	var cost int
	for _, p := range initial {
		cost += lib.Abs(target - p)
	}
	return cost
}
