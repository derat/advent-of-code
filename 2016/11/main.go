package main

import (
	"fmt"
	"regexp"

	"github.com/derat/advent-of-code/lib"
)

const bits = 7 // bits used to represent a full set of elements

func main() {
	// Parsing this input is painful.
	sepRegexp := regexp.MustCompile(`(, and |, | and )`)
	elements := make(map[string]int)
	var floors int
	initial := setFloorNum(0, 0) // we start on the bottom floor
	for floor, ln := range lib.InputLines("2016/11") {
		floors++
		var rest string
		lib.Extract(ln, `^The \w+ floor contains (.+)\.$`, &rest)
		var chips, gens int
		if rest != "nothing relevant" {
			for _, p := range sepRegexp.Split(rest, -1) {
				var element, thing string
				lib.Extract(p, `^a (\w+)(?:-compatible)? (microchip|generator)$`, &element, &thing)
				id, ok := elements[element]
				if !ok {
					id = len(elements)
					elements[element] = id
				}
				if thing == "microchip" {
					chips |= 1 << id
				} else {
					gens |= 1 << id
				}
			}
		}
		initial = setFloorItems(initial, floor, chips, gens)
	}

	fmt.Println(solve(initial, floors, len(elements)))

	// Part 2: Chips and generators for two more elements on first floor.
	elerium := len(elements)
	elements["elerium"] = elerium
	dilithium := len(elements)
	elements["dilithium"] = dilithium
	updated, ok := update(initial, 0, true, thing{elerium, true}, thing{elerium, false},
		thing{dilithium, true}, thing{dilithium, false})
	lib.Assertf(ok, "Invalid initial state for part 2")
	fmt.Println(solve(updated, floors, len(elements)))
}

// solve returns the number of moves necessary to get from the supplied initial
// state to having all elements on the top floor.
// This approach (BFS) is terrible and slow (it takes a few minutes to complete
// for part 2), but it doesn't blow up memory so I'm running with it.
func solve(initial uint64, floors, numElements int) int {
	lib.AssertLessEq(numElements, bits)

	// We want to get all items to the top floor.
	targetMask := setFloorItems(0, floors-1, 1<<numElements-1, 1<<numElements-1)

	seen := map[uint64]struct{}{initial: struct{}{}}
	alreadySeen := func(st uint64) bool {
		_, ok := seen[st]
		return ok
	}

	todo := map[uint64]struct{}{initial: struct{}{}}
	var moves int
	for ; ; moves++ {
		lib.Assertf(len(todo) != 0, "No new states to check")
		nextTodo := make(map[uint64]struct{})
		for st := range todo {
			// If everything is on the top floor, we're done.
			if st&targetMask == targetMask {
				return moves
			}

			floor := getFloorNum(st)
			chips, gens := getFloorItems(st, floor)
			things := getThings(chips, gens)
			lib.Assertf(len(things) > 0, "No items on floor %d", floor)

			// Go down with one or two items from the current floor.
			if floor > 0 {
				st := setFloorNum(st, floor-1)
				for i, th1 := range things {
					if st, ok := move(st, floor, floor-1, th1); ok && !alreadySeen(st) {
						nextTodo[st] = struct{}{}
					}
					for _, th2 := range things[i+1:] {
						if st, ok := move(st, floor, floor-1, th1, th2); ok && !alreadySeen(st) {
							nextTodo[st] = struct{}{}
						}
					}
				}
			}

			// Go up with one or two items from the current floor.
			if floor < floors-1 {
				st := setFloorNum(st, floor+1)
				for i, th1 := range things {
					if st, ok := move(st, floor, floor+1, th1); ok && !alreadySeen(st) {
						nextTodo[st] = struct{}{}
					}
					for _, th2 := range things[i+1:] {
						if st, ok := move(st, floor, floor+1, th1, th2); ok && !alreadySeen(st) {
							nextTodo[st] = struct{}{}
						}
					}
				}
			}
		}
		todo = nextTodo
	}
	return moves
}

func getFloorItems(state uint64, floor int) (chips, gens int) {
	chips = lib.UnpackInt(state, bits, (floor*2+1)*bits) // bottom field is floor num
	gens = lib.UnpackInt(state, bits, ((floor+1)*2)*bits)
	return chips, gens
}

func setFloorItems(state uint64, floor, chips, gens int) uint64 {
	state = lib.PackInt(state, chips, bits, (floor*2+1)*bits)
	state = lib.PackInt(state, gens, bits, ((floor+1)*2)*bits)
	return state
}

func getFloorNum(state uint64) int {
	return lib.UnpackInt(state, bits, 0)
}

func setFloorNum(state uint64, floor int) uint64 {
	return lib.PackInt(state, floor, bits, 0)
}

type thing struct {
	id   int
	chip bool
}

func getThings(chips, gens int) []thing {
	var things []thing
	for id := 0; chips != 0 || gens != 0; id++ {
		if chips&0x1 != 0 {
			things = append(things, thing{id, true})
		}
		if gens&0x1 != 0 {
			things = append(things, thing{id, false})
		}
		chips >>= 1
		gens >>= 1
	}
	return things
}

func update(st uint64, floor int, add bool, things ...thing) (uint64, bool) {
	chips, gens := getFloorItems(st, floor)
	for _, th := range things {
		if add {
			if th.chip {
				chips |= 1 << th.id
			} else {
				gens |= 1 << th.id
			}
		} else {
			if th.chip {
				chips &= ^(1 << th.id)
			} else {
				gens &= ^(1 << th.id)
			}
		}
	}
	st = setFloorItems(st, floor, chips, gens)

	protected := chips & gens
	valid := gens == 0 || chips & ^protected == 0

	return st, valid
}

func move(st uint64, oldFloor, newFloor int, things ...thing) (uint64, bool) {
	st, ok := update(st, oldFloor, false, things...)
	if !ok {
		return st, false
	}
	return update(st, newFloor, true, things...)
}
