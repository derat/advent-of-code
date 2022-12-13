package main

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var packets [][]any
	for _, pg := range lib.InputParagraphs("2022/13") {
		lib.AssertEq(len(pg), 2)
		left, ln := parse(pg[0])
		lib.AssertEq(ln, len(pg[0]))
		right, rn := parse(pg[1])
		lib.AssertEq(rn, len(pg[1]))
		packets = append(packets, left, right)
	}

	// Part 1: get sum of (1-based) indexes of pairs that are properly-ordered.
	var sum int
	for i := 0; i < len(packets)/2; i++ {
		if cmp(packets[i*2], packets[i*2+1]) < 0 {
			sum += i + 1
		}
	}
	fmt.Println(sum)

	// Part 2: add two hardcoded divider packets, sort, and multiply their (1-based) indexes.
	div1 := []any{[]any{2}}
	div2 := []any{[]any{6}}
	packets = append(packets, div1, div2)
	sort.Slice(packets, func(i, j int) bool { return cmp(packets[i], packets[j]) < 0 })
	var i1, i2 int // 1-based
	for i, p := range packets {
		if reflect.DeepEqual(p, div1) {
			i1 = i + 1
		} else if reflect.DeepEqual(p, div2) {
			i2 = i + 1
		}
	}
	lib.Assert(i1 != 0)
	lib.Assert(i2 != 0)
	fmt.Println(i1 * i2)
}

// parse parses a '['-prefixed string like "[1,[2,[3,[4,[5,6,7]]]],8,9]"
// into a corresponding slice of integers and nested slices. The number of
// consumed characters (including opening and closing brackets) is also returned.
func parse(s string) ([]any, int) {
	lib.AssertEq(s[0], '[')
	var ret []any
	var inNum bool
	var num int
	finishNum := func() {
		if inNum {
			ret = append(ret, num)
			inNum = false
		}
	}

	i := 1
	for i < len(s) {
		ch := s[i]
		switch {
		case ch == '[':
			v, n := parse(s[i:])
			ret = append(ret, v)
			i += n
		case ch == ']':
			finishNum()
			return ret, i + 1
		case ch == ',':
			finishNum()
			i++
		case ch >= '0' && ch <= '9':
			d := int(ch - '0')
			if !inNum {
				inNum = true
				num = d
			} else {
				num = 10*num + d
			}
			i++
		default:
			lib.Panicf("Invalid character %q", ch)
		}
	}
	panic("Missing closing bracket")
}

// cmp recursively compares left and right (either int or []any) and returns a negative value if
// left comes first, a positive value if right comes first, and zero if the values are equal.
func cmp(left, right any) int {
	li, lok := left.(int)
	ri, rok := right.(int)
	switch {
	case lok && rok:
		// "If both values are integers, the lower integer should come first. If the left integer is
		// lower than the right integer, the inputs are in the right order. If the left integer is
		// higher than the right integer, the inputs are not in the right order. Otherwise, the
		// inputs are the same integer; continue checking the next part of the input."
		return li - ri
	case lok:
		// "If exactly one value is an integer, convert the integer to a list which contains that
		// integer as its only value, then retry the comparison."
		return cmp([]any{left}, right)
	case rok:
		// See previous case.
		return cmp(left, []any{right})
	default:
		// "If both values are lists, compare the first value of each list, then the second value,
		// and so on."
		ll := left.([]any)
		rl := right.([]any)
		for i := 0; i < len(ll) && i < len(rl); i++ {
			if v := cmp(ll[i], rl[i]); v != 0 {
				return v
			}
		}
		// "If the left list runs out of items first, the inputs are in the right order. If the
		// right list runs out of items first, the inputs are not in the right order. If the lists
		// are the same length and no comparison makes a decision about the order, continue checking
		// the next part of the input."
		return len(ll) - len(rl)
	}
}
