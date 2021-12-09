package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	rules := make(map[string]string) // 2->3 and 3->4
	for _, ln := range lib.InputLines("2017/21") {
		var src, dst string
		lib.Extract(ln, `^([.#/]+) => ([.#/]+)$`, &src, &dst)
		for _, s := range transforms(src) {
			rules[s] = dst
		}
	}

	// In each dimension:
	//
	//  - 3 maps to one 4.
	//  - 4 maps to one 6.
	//  - 6 maps to three 3s.
	//
	//    3 ->   4 (3->4)
	//    4 ->   6 (4->6)
	//    6 ->   9 (6->3)
	//    9 ->  12 (3->4)
	//   12 ->  18 (4->6)
	//   18 ->  27 (6->3)
	//   27 ->  36 (3->4)
	//   36 ->  54 (4->6)
	//   54 ->  81 (6->3)
	//   81 -> 108 (3->4)
	//  108 -> 162 (4->6)

	exps := make(map[string][]string) // expansions from one pattern to the others that it creates
	var add func(string)
	add = func(src string) {
		var dsts []string
		switch len(src) {
		case 11: // 3x3 -> one 4x4
			dsts = []string{rules[src]}
		case 19: // 4x4 -> one 6x6
			var subs [][]string
			for _, sub := range split2(src) { // split into four 2x2s
				subs = append(subs, strings.Split(rules[sub], "/")) // convert each to 3x3
			}
			dst := []string{ // join to 6x6
				subs[0][0] + subs[1][0],
				subs[0][1] + subs[1][1],
				subs[0][2] + subs[1][2],
				subs[2][0] + subs[3][0],
				subs[2][1] + subs[3][1],
				subs[2][2] + subs[3][2],
			}
			dsts = []string{strings.Join(dst, "/")}
		case 41: // 6x6 -> nine 3x3
			for _, sub := range split2(src) { // split into nine 2x2s
				dsts = append(dsts, rules[sub]) // convert each to 3x3
			}
		default:
		}

		exps[src] = dsts
		for _, dst := range dsts {
			if _, ok := exps[dst]; !ok {
				add(dst)
			}
		}
	}

	const init = ".#./..#/###" // initial pattern given in puzzle
	add(init)                  // generate all possible 3x3, 4x4, and 6x6 expansions

	var count func(string, int) int
	count = func(s string, n int) int {
		if n == 0 {
			return strings.Count(s, "#")
		}
		var sum int
		for _, d := range exps[s] {
			sum += count(d, n-1)
		}
		return sum
	}
	fmt.Println(count(init, 5))
	fmt.Println(count(init, 18))
}

// split2 splits the supplied string into 2x2 blocks.
func split2(s string) []string {
	var res []string
	rows := strings.Split(s, "/")
	for r := 0; 2*r < len(rows); r++ {
		for c := 0; 2*c < len(rows[r]); c++ {
			sub := make([]string, 2)
			for i := 0; i < 2; i++ {
				sub[i] = rows[2*r+i][2*c : 2*(c+1)]
			}
			res = append(res, strings.Join(sub, "/"))
		}
	}
	return res
}

// transforms returns all transformations that can be produced
// by flipping and rotating the supplied string (e.g. "#./.#").
// This is similar to 2020/20.
func transforms(s string) []string {
	rows := lib.ByteGrid(bytes.Split([]byte(s), []byte("/")))
	lib.AssertEq(len(rows), len(rows[0])) // square

	ts := make([]string, 8)
	for tr := 0; tr < 8; tr++ {
		trows := rows.Copy()
		if tr&1 != 0 {
			trows = trows.FlipVert()
		}
		if tr&2 != 0 {
			trows = trows.FlipHoriz()
		}
		if tr&4 != 0 {
			trows = trows.RotateCW()
		}
		ts[tr] = string(bytes.Join(trows, []byte("/")))
	}

	// Remove duplicates.
	sort.Strings(ts)
	i := 0
	for j := 1; j < len(ts); j++ {
		if ts[i] == ts[j] {
			continue
		}
		i++
		ts[i] = ts[j]
	}
	ts = ts[:i+1]

	return ts
}
