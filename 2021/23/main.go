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
	start := [p1len]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		tr(init[0]), tr(init[1]), tr(init[2]), tr(init[3]),
		tr(init[4]), tr(init[5]), tr(init[6]), tr(init[7])}
	end := [p1len]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 1, 2, 3, 4}

	// Part 1: "What is the least energy required to organize the amphipods?"
	cost := lib.AStar([]interface{}{start},
		func(si interface{}) bool { return si.([p1len]typ) == end },
		func(si interface{}, next map[interface{}]int) {
			s := si.([p1len]typ)
			for start, t := range s {
				if t == empty {
					continue
				}
				for end := range s {
					if legal(start, end, s[:]) {
						ns := s
						ns[start] = empty
						ns[end] = t
						cost := tinfos[t].cost * len(getPath(start, end))
						next[ns] = cost
					}
				}
			}
		},
		func(si interface{}) int {
			// Return the sum of the minimum cost to move each amphipod to its room.
			var sum int
			for start, t := range si.([p1len]typ) {
				if t == empty {
					continue
				}
				ti := tinfos[t]
				if start == ti.outer || start == ti.inner {
					continue
				}
				sum += ti.cost * len(getPath(start, ti.outer))
			}
			return sum
		})
	fmt.Println(cost)

	// Part 2: "Using the initial configuration from the full diagram, what is the least energy
	// required to organize the amphipods?"
	//
	// This is almost identical to the previous call, but I had to change the
	// type conversions and also add innerer/innerest to the already-there check
	// in the lower-bound function.
	//
	// My annoying bug for part 2 was that I missed updating getRoomInfo to additionally
	// check innerer and innerest, resulting in no paths being found.
	start2 := [p2len]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		start[11], start[12], start[13], start[14],
		4, 3, 2, 1, // given in puzzle
		4, 2, 1, 3, // given in puzzle
		start[15], start[16], start[17], start[18]}
	end2 := [p2len]typ{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4, 1, 2, 3, 4}
	cost2 := lib.AStar([]interface{}{start2},
		func(si interface{}) bool { return si.([p2len]typ) == end2 },
		func(si interface{}, next map[interface{}]int) {
			s := si.([p2len]typ)
			for start, t := range s {
				if t == empty {
					continue
				}
				for end := range s {
					if legal(start, end, s[:]) {
						ns := s
						ns[start] = empty
						ns[end] = t
						cost := tinfos[t].cost * len(getPath(start, end))
						next[ns] = cost
					}
				}
			}
		},
		func(si interface{}) int {
			var sum int
			for start, t := range si.([p2len]typ) {
				if t == empty {
					continue
				}
				ti := tinfos[t]
				if start == ti.outer || start == ti.inner || start == ti.innerer || start == ti.innerest {
					continue
				}
				sum += ti.cost * len(getPath(start, ti.outer))
			}
			return sum
		})
	fmt.Println(cost2)
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
//  0  1  2  3  4  5  6  7  8  9  10
//       11    12    13    14
//       15    16    17    18
//
// Update after looking at part 2: I again see that I have bit-packed into
// a uint64 prematurely; I won't be able to fit the additional 8 positions
// using my previous scheme. It was also completely pointless, as my part 1
// runtime (around 150 ms) doesn't change if I just use [19]typ as map keys
// instead of uint64.
//
//  0  1  2  3  4  5  6  7  8  9  10
//       11    12    13    14
//       15    16    17    18
//       19    20    21    22
//       23    24    25    26

const (
	p1len = 19
	p2len = 27
)

type typ int

const empty typ = 0

func dump(b []typ) string {
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
	if len(b) == p2len {
		s += "  #" + ch(b[19]) + "#" + ch(b[20]) + "#" + ch(b[21]) + "#" + ch(b[22]) + "#\n"
		s += "  #" + ch(b[23]) + "#" + ch(b[24]) + "#" + ch(b[25]) + "#" + ch(b[26]) + "#\n"
	}
	s += "  #########\n"
	return s
}

type tinfo struct {
	t                 typ
	door              int // index of hallway just outside room
	outer, inner      int // index of outer and inner room position
	innerer, innerest int // part 2: like inner but more so
	cost              int // cost of moving a single step

}

var tinfos = []tinfo{
	tinfo{0, -1, -1, -1, -1, -1, 0},   // make it easy to index with type
	tinfo{1, 2, 11, 15, 19, 23, 1},    // "Amber amphipods require 1 energy per step"
	tinfo{2, 4, 12, 16, 20, 24, 10},   // "Bronze amphipods require 10 energy"
	tinfo{3, 6, 13, 17, 21, 25, 100},  // "Copper amphipods require 100"
	tinfo{4, 8, 14, 18, 22, 26, 1000}, // "Desert ones require 1000"
}

var paths = map[[2]int][]int{} // map from start/end to steps (including end)

// getPath returns the shortest path from start to end.
func getPath(start, end int) []int {
	p, ok := paths[[2]int{start, end}]
	lib.Assertf(ok, "No path from %d to %d", start, end)
	return p
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

	// Flattens ps into a new slice.
	concat := func(ps ...[]int) []int {
		var path []int
		for _, p := range ps {
			path = append(path, p...)
		}
		return path
	}
	// Adds bidirectional links directly between start and end.
	direct := func(start, end int) {
		paths[[2]int{start, end}] = []int{end}
		paths[[2]int{end, start}] = []int{start}
	}
	// Adds bidirectional links between start and end, passing through mid.
	join := func(start, mid, end int) {
		paths[[2]int{start, end}] = concat(getPath(start, mid), getPath(mid, end))
		paths[[2]int{end, start}] = concat(getPath(end, mid), getPath(mid, start))
	}

	for _, t := range tinfos[1:] {
		// Generate paths between rooms and their doorways.
		direct(t.door, t.outer)
		direct(t.outer, t.inner)
		direct(t.inner, t.innerer)
		direct(t.innerer, t.innerest)
		join(t.inner, t.innerer, t.innerest)
		join(t.outer, t.inner, t.innerer)
		join(t.outer, t.innerer, t.innerest)
		join(t.door, t.outer, t.inner)
		join(t.door, t.inner, t.innerer)
		join(t.door, t.inner, t.innerest)

		// Generate paths between rooms and all hallway positions besides their doors.
		for i := 0; i <= 10; i++ {
			if i == t.door {
				continue
			}
			join(t.outer, t.door, i)
			join(t.inner, t.door, i)
			join(t.innerer, t.door, i)
			join(t.innerest, t.door, i)
		}
	}

	// Finally, generate paths between different rooms.
	for i, st := range tinfos[1:] {
		for j, et := range tinfos[1:] {
			if i == j {
				continue
			}
			// We'll double-add some paths here, but it shouldn't matter.
			for _, s := range []int{st.outer, st.inner, st.innerer, st.innerest} {
				for _, e := range []int{et.outer, et.inner, et.innerer, et.innerest} {
					join(s, st.door, e)
				}
			}
		}
	}
}

// getRoomInfo returns information for the room at position p.
// The returned bool is false if p is not part of a room.
func getRoomInfo(p int) (tinfo, bool) {
	for _, ti := range tinfos[1:] {
		if p == ti.outer || p == ti.inner || p == ti.innerer || p == ti.innerest {
			return ti, true
		}
	}
	return tinfo{}, false
}

// legal returns true if the amphipod at position start in b can move to position end.
func legal(start, end int, b []typ) bool {
	lib.AssertLess(start, len(b))
	lib.AssertLess(end, len(b))

	t := b[start]
	if t == empty || start == end || b[end] != empty {
		return false
	}

	// "Amphipods will never stop on the space immediately outside any room.
	// They can move into that space so long as they immediately continue moving."
	for _, ti := range tinfos[1:] {
		if end == ti.door {
			return false
		}
	}

	si, sr := getRoomInfo(start)
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

		// Additional optimization: don't move around once we're already in the right
		// destination room.
		if sr && si.t == ei.t {
			return false
		}

		// Sigh. I wasted a bunch of time staring at my cost-estimation code when I actually
		// had a bug here: I was checking et.outer and et.inner instead of b[et.outer] and
		// b[et.inner]. I've introduced a 'typ' type in the meantime that would've caught
		// this.
		//
		// Adding the type paid off, since I immediately made the same mistake again when
		// updating the 'additional optimization' section for part 2.
		if et := b[ei.outer]; et != empty && et != t {
			return false
		}
		if et := b[ei.inner]; et != empty && et != t {
			return false
		}

		// Part 2: Also check deeper areas of rooms.
		if len(b) == p2len {
			if et := b[ei.innerer]; et != empty && et != t {
				return false
			}
			if et := b[ei.innerest]; et != empty && et != t {
				return false
			}
		}

		// Additional optimization: Don't move to the outer part of a room if the inner
		// part is empty; it's pointless.
		if end == ei.outer && b[ei.inner] == empty {
			return false
		}
		if len(b) == p2len {
			if end == ei.inner && b[ei.innerer] == empty {
				return false
			}
			if end == ei.innerer && b[ei.innerest] == empty {
				return false
			}
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
