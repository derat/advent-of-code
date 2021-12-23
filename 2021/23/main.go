package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLines("2021/23")
	var init [8]byte
	lib.AssertEq(len(input), 5)
	lib.AssertEq(input[0], "#############")
	lib.AssertEq(input[1], "#...........#")
	lib.Extract(input[2], "^###([A-D])#([A-D])#([A-D])#([A-D])###$",
		&init[0], &init[1], &init[2], &init[3])
	lib.Extract(input[3], "^  #([A-D])#([A-D])#([A-D])#([A-D])#$",
		&init[4], &init[5], &init[6], &init[7])
	lib.AssertEq(input[4], "  #########")

	tr := func(ch byte) typ { return typ(ch - 'A' + 1) }
	start := pack([19]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		tr(init[0]), tr(init[1]), tr(init[2]), tr(init[3]),
		tr(init[4]), tr(init[5]), tr(init[6]), tr(init[7])})
	end := pack([19]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 1, 2, 3, 4})

	cost := lib.AStar([]interface{}{start},
		func(si interface{}) bool { return si.(uint64) == end },
		func(si interface{}, next map[interface{}]int) {
			s := unpack(si.(uint64))
			for start, t := range s {
				if t == 0 {
					continue
				}
				for end := range s {
					if legal(start, end, s) {
						ns := s
						ns[start] = 0
						ns[end] = t
						cost := tinfos[t].cost * len(getPath(start, end))
						next[pack(ns)] = cost
					}
				}
			}
		},
		func(si interface{}) int {
			// Return the sum of the minimum cost to move each amphipod to its room.
			var sum int
			for start, t := range unpack(si.(uint64)) {
				if t == 0 {
					continue
				}
				ti := tinfos[t]
				if start == ti.inner || start == ti.outer {
					continue
				}
				sum += ti.cost * len(getPath(start, ti.outer))
			}
			return sum
		})
	fmt.Println(cost)
}

// The shape of the burrow appears to be fixed, e.g.:
//
//  #############
//  #...........#
//  ###B#C#B#D###
//    #A#D#C#A#
//    #########
//
// (Note that the last two lines don't have trailing spaces.)
//
// There are 19 spaces, and each can have 5 states (empty, A, B, C, D),
// so it looks like I should be able to pack the spaces into an int64
// as 3-bit numbers. I know from a past year's puzzle that I should
// treat all the amphipods of a given type as interchangeable to count
// down on the total number of possible states, but this approach should
// do that already.
//
//  0123456789A
//    B C D E
//    F G H I

type typ int

const empty typ = 0

func pack(b [19]typ) (s uint64) {
	for i, t := range b {
		s = lib.PackInt(s, int(t), 3, 3*i)
	}
	return s
}

func unpack(s uint64) (b [19]typ) {
	for i := range b {
		b[i] = typ(lib.UnpackInt(s, 3, 3*i))
	}
	return b
}

func dump(b [19]typ) string {
	ch := func(t typ) string {
		if t == empty {
			return "."
		}
		return string('A' + rune(t-1))
	}
	s := "#############\n"
	s += "#"
	for _, t := range b[0:11] {
		s += ch(t)
	}
	s += "#\n"
	s += "###" + ch(b[11]) + "#" + ch(b[12]) + "#" + ch(b[13]) + "#" + ch(b[14]) + "###\n"
	s += "  #" + ch(b[15]) + "#" + ch(b[16]) + "#" + ch(b[17]) + "#" + ch(b[18]) + "#\n"
	s += "  #########\n"
	return s
}

type tinfo struct {
	t            typ
	door         int // index of hallway just outside room
	outer, inner int // index of outer and inner room position
	cost         int // cost of moving a single step
}

var tinfos = []tinfo{
	tinfo{0, -1, -1, -1, 0},   // make it easy to index with type
	tinfo{1, 2, 11, 15, 1},    // "Amber amphipods require 1 energy per step"
	tinfo{2, 4, 12, 16, 10},   // "Bronze amphipods require 10 energy"
	tinfo{3, 6, 13, 17, 100},  // "Copper amphipods require 100"
	tinfo{4, 8, 14, 18, 1000}, // "Desert ones require 1000"
}

var paths = map[[2]int][]int{} // map from start/end to steps (including end)

// getPath returns the shortest path from start to end.
func getPath(start, end int) []int {
	p, ok := paths[[2]int{start, end}]
	lib.Assertf(ok, "No path from %d to %d", start, end)
	return p
}

// concat flattens ps into a new slice.
func concat(ps ...[]int) []int {
	var path []int
	for _, p := range ps {
		path = append(path, p...)
	}
	return path
}

func init() {
	// First, generate paths along the hallway (positions 0 to 10, inclusive).
	for start := 0; start <= 10; start++ {
		for end := 0; end <= 10; end++ {
			var p []int
			if start <= end {
				for i := start + 1; i <= end; i++ {
					p = append(p, i)
				}
			} else if start > end {
				for i := start - 1; i >= end; i-- {
					p = append(p, i)
				}
			} else {
				continue
			}
			paths[[2]int{start, end}] = p
		}
	}

	// Next, generate paths between rooms and their doorways.
	for _, t := range tinfos[1:] {
		paths[[2]int{t.door, t.outer}] = []int{t.outer}
		paths[[2]int{t.door, t.inner}] = []int{t.outer, t.inner}
		paths[[2]int{t.outer, t.door}] = []int{t.door}
		paths[[2]int{t.outer, t.inner}] = []int{t.inner}
		paths[[2]int{t.inner, t.outer}] = []int{t.outer}
		paths[[2]int{t.inner, t.door}] = []int{t.outer, t.door}

		for i := 0; i <= 10; i++ {
			if i == t.door {
				continue
			}
			paths[[2]int{t.outer, i}] = concat(getPath(t.outer, t.door), getPath(t.door, i))
			paths[[2]int{i, t.outer}] = concat(getPath(i, t.door), getPath(t.door, t.outer))
			paths[[2]int{t.inner, i}] = concat(getPath(t.inner, t.door), getPath(t.door, i))
			paths[[2]int{i, t.inner}] = concat(getPath(i, t.door), getPath(t.door, t.inner))
		}
	}

	// Finally, generate paths between different rooms.
	for i, st := range tinfos[1:] {
		for j, et := range tinfos[1:] {
			if i == j {
				continue
			}
			paths[[2]int{st.outer, et.outer}] = concat(getPath(st.outer, st.door), getPath(st.door, et.outer))
			paths[[2]int{st.outer, et.inner}] = concat(getPath(st.outer, st.door), getPath(st.door, et.inner))
			paths[[2]int{st.inner, et.outer}] = concat(getPath(st.inner, st.door), getPath(st.door, et.outer))
			paths[[2]int{st.inner, et.inner}] = concat(getPath(st.inner, st.door), getPath(st.door, et.inner))
		}
	}
}

func getRoomInfo(p int) (tinfo, bool) {
	for _, ti := range tinfos[1:] {
		if p == ti.outer || p == ti.inner {
			return ti, true
		}
	}
	return tinfo{}, false
}

func legal(start, end int, b [19]typ) bool {
	t := b[start]
	if t == empty || start == end || b[end] != 0 {
		return false
	}

	// "Amphipods will never stop on the space immediately outside any room.
	// They can move into that space so long as they immediately continue moving."
	for _, ti := range tinfos[1:] {
		if end == ti.door {
			return false
		}
	}

	_, sr := getRoomInfo(start)
	ei, er := getRoomInfo(end)

	// "Once an amphipod stops moving in the hallway, it will stay in that spot
	// until it can move into a room."
	if !sr && !er {
		return false
	}

	// "Amphipods will never move from the hallway into a room unless that room is
	// their destination room and that room contains no amphipods which do not also
	// have that room as their own destination."
	if er {
		if ei.t != t {
			return false
		}
		// Sigh. I wasted a bunch of time staring at my cost-estimation code when I actually
		// had a bug here: I was checking et.outer and et.inner instead of b[et.outer] and
		// b[et.inner]. I've introduced a 'typ' type in the meantime that would've caught
		// this.
		if et := b[ei.outer]; et != empty && et != t {
			return false
		}
		if et := b[ei.inner]; et != empty && et != t {
			return false
		}
		// Additional optimization: Don't move to the outer part of a room if the inner
		// part is empty; it's pointless.
		if end == ei.outer && b[ei.inner] == empty {
			return false
		}
	}

	// Check that the path isn't blocked.
	for _, p := range getPath(start, end) {
		if b[p] != 0 {
			return false
		}
	}

	return true
}
