package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	vals := lib.ExtractDigits(lib.ReadLines()[0])

	// Part 1:
	first, last, one, _ := makeCups(vals)
	for i := 0; i < 100; i++ {
		first, last = move(first, last)
	}

	// Print all the cups after the one labeled 1.
	for c := one.next; c != one; c = c.next {
		fmt.Print(c.val)
	}
	fmt.Println()

	// Part 2:
	first, last, one, max := makeCups(vals)

	// Add additional cups.
	for v := max.val + 1; v <= 1000000; v++ {
		c := &cup{val: v}
		if v == max.val+1 {
			c.less = max
		} else {
			c.less = last
		}
		last.next = c
		last = c
		one.less = c
	}

	// Run 10 million iterations and print the product of the two cups after 1.
	for i := 0; i < 10000000; i++ {
		first, last = move(first, last)
	}
	fmt.Println(int64(one.next.val) * int64(one.next.next.val))
}

type cup struct {
	val        int32
	next, less *cup // next (clockwise) and cup with next-lowest value
}

func makeCups(vals []int) (first, last, one, max *cup) {
	for _, v := range vals {
		c := &cup{val: int32(v)}
		if first == nil {
			first = c
		} else {
			last.next = c
		}
		last = c

		if c.val == 1 {
			one = c
		}
		if max == nil || c.val > max.val {
			max = c
		}
	}

	// Add pointers to cups with decreased value.
	for c := first; c != nil; c = c.next {
		if c == one {
			c.less = max
		} else {
			for o := first; o != nil; o = o.next {
				if c.val-1 == o.val {
					c.less = o
					break
				}
			}
		}
	}

	return
}

func move(first, last *cup) (*cup, *cup) {
	// Remove the second, third, and fourth cups.
	r1 := first.next
	r2 := r1.next
	r3 := r2.next
	first.next = r3.next
	r3.next = nil

	// Find the cup with a value one less than the first cup,
	// counting down further if it's one of the removed cups.
	less := first.less
	for less == r1 || less == r2 || less == r3 {
		less = less.less
	}

	// Splice in the removed cups.
	r3.next = less.next
	less.next = r1
	if less == last {
		last = r3
	}

	// Move the first cup to the end of the list.
	last.next = first
	last = first
	first = first.next

	return first, last
}
