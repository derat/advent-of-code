package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	// My input has 60 valves, of which 15 have non-zero flow rates.
	// So, there are 60 * 2**15 = 1966080 possible states for a given
	// round (actually, we probably don't need to track the position
	// once all valves are open, since there's no point in moving then).
	// For a given round, we only care about the best score that's possible
	// for each state. Maybe there are additional optimizations to drop
	// purely-worse states.

	rates := make(map[string]int)
	tunnels := make(map[string][]string)
	for _, ln := range lib.InputLines("2022/16") {
		var name, dsts string
		var rate int
		lib.Extract(ln, `^Valve ([A-Z][A-Z]) has flow rate=(\d+); tunnels? leads? to valves? (.+)$`,
			&name, &rate, &dsts)
		if rate > 0 {
			rates[name] = rate
		}
		for _, dst := range strings.Split(dsts, ", ") {
			tunnels[name] = append(tunnels[name], dst)
		}
	}

	lib.AssertLessEq(len(tunnels), 64)
	lib.AssertLessEq(len(rates), 16)

	// Assign an ID to each valve. Give lower IDs to valves with non-zero flow rates
	// so they can be used in a smaller bitfield.
	names := lib.MapKeys(tunnels)
	sort.Slice(names, func(a, b int) bool {
		na, nb := names[a], names[b]
		if ra, rb := rates[na], rates[nb]; ra > 0 && rb == 0 {
			return true
		} else if rb > 0 && ra == 0 {
			return false
		}
		return na < nb
	})
	valveIDs := make(map[string]uint8, len(names))
	for i, n := range names {
		valveIDs[n] = uint8(i)
	}
	idValves := lib.InvertMap(valveIDs)

	var allOpen uint16
	for n, r := range rates {
		lib.AssertGreater(r, 0)
		allOpen |= 1 << valveIDs[n]
	}

	// Returns the flow rate for one minute with the specified valves open.
	minuteRate := func(open uint16) int {
		var rate int
		for n, r := range rates {
			if open&(1<<valveIDs[n]) != 0 {
				rate += r
			}
		}
		return rate
	}

	// Part 1: Find the maximum pressure that can be released in 30 minutes.
	// TODO: This produces the wrong value for the example input (1639 instead of 1651),
	// but somehow produces the correct answer when given my actual input for the day.
	start := state{valveIDs["AA"], 0x0, 30}
	fmt.Println(-1 * lib.AStar(
		[]state{start},
		func(s state) bool { return s.rem == 0 },
		func(s state, next map[state]int) {
			name := idValves[s.cur]
			lib.Assert(name != "")
			cost := -1 * minuteRate(s.open)

			// If the current valve has a non-zero rate and is closed,
			// we can open it during this turn.
			if r := rates[name]; r > 0 && s.open&(1<<s.cur) == 0 {
				ns := s
				ns.open |= 1 << s.cur
				ns.rem--
				next[ns] = cost
			}

			// If all the valves aren't already open, also add states where we move
			// to surrounding valves.
			if s.open != allOpen {
				for _, dn := range tunnels[name] {
					ns := s
					ns.cur = valveIDs[dn]
					ns.rem--
					next[ns] = cost
				}
			}
		},
		func(s state) int {
			// As a bad lower bound on cost, take the released pressure from the already-open
			// valves plus the pressure we would get if we opened the remaining valves on the
			// next turn (which we can't do).
			closed := allOpen & ^s.open
			return -int(s.rem)*minuteRate(s.open) - int(s.rem-1)*minuteRate(closed)
		}))

	// Part 2: One elephant is also moving/opening, and you only have 26 minutes.
	// We and the elephant act the same, so just track positions in a single bitfield
	// to treat us interchangeably to reduce the number of states.
	// TODO: Sigh, this also produces the wrong output for the example input
	// (1675 rather than 1707) but succeeds with my actual input.
	// It's also extremely slow (30s), maybe because of my lazy estimate function
	// that I just reused from part 1.
	start2 := state2{1 << valveIDs["AA"], 0x0, 26}
	fmt.Println(-1 * lib.AStar(
		[]state2{start2},
		func(s state2) bool { return s.rem == 0 },
		func(s state2, next map[state2]int) {
			// Extract the current positions from the state.
			var name1, name2 string
			for id, n := range idValves {
				if s.cur&(1<<id) != 0 {
					if name1 == "" {
						name1 = n
					} else {
						lib.AssertEq(name2, "")
						name2 = n
					}
				}
			}
			lib.Assert(name1 != "")
			if name2 == "" {
				name2 = name1
			}

			id1, id2 := valveIDs[name1], valveIDs[name2]
			r1, r2 := rates[name1], rates[name2]
			o1 := r1 > 0 && s.open&(1<<id1) != 0
			o2 := r2 > 0 && s.open&(1<<id2) != 0
			cost := -1 * minuteRate(s.open)

			// Either we and the elephant can both open valves (if we're in different locations),
			// we and the elephant can both move, or one of us can open a valve and the other can
			// move. Copy-and-paste city.

			// Both move.
			if s.open != allOpen {
				for _, dn1 := range tunnels[name1] {
					for _, dn2 := range tunnels[name2] {
						ns := s
						ns.cur = (1 << valveIDs[dn1]) | (1 << valveIDs[dn2])
						ns.rem--
						next[ns] = cost
					}
				}
			}

			// Both open.
			if id1 != id2 && r1 > 0 && !o1 && r2 > 0 && !o2 {
				ns := s
				ns.open |= (1 << id1) | (1 << id2)
				ns.rem--
				next[ns] = cost
			}

			// One opens, other moves.
			if r1 > 0 && !o1 {
				ns := s
				ns.open |= (1 << id1)
				ns.rem--
				next[ns] = cost
				for _, dn := range tunnels[name2] {
					tmp := ns
					tmp.cur = (1 << id1) | (1 << valveIDs[dn])
					next[tmp] = cost
				}
			}
			if r2 > 0 && !o2 {
				ns := s
				ns.open |= (1 << id2)
				ns.rem--
				next[ns] = cost
				for _, dn := range tunnels[name1] {
					tmp := ns
					tmp.cur = (1 << id2) | (1 << valveIDs[dn])
					next[tmp] = cost
				}
			}
		},
		func(s state2) int {
			closed := allOpen & ^s.open
			return -int(s.rem)*minuteRate(s.open) - int(s.rem-1)*minuteRate(closed)
		}))
}

type state struct {
	cur  uint8  // id of current valve
	open uint16 // bitfield of open valves (1 << id)
	rem  uint8  // remaining minutes
}

type state2 struct {
	cur  uint64 // bitfield of my location and the elephant's (1 << id)
	open uint16 // bitfield of open valves (1 << id)
	rem  uint8  // remaining minutes
}
