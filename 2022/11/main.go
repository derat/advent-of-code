package main

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var monkeys []*monkey
	var orig [][]int // save orig items for part 2
	for i, pg := range lib.InputParagraphs("2022/11") {
		m := newMonkey(pg)
		lib.AssertEq(m.id, i)
		monkeys = append(monkeys, m)
		orig = append(orig, append([]int(nil), m.items...))
	}

	// Part 1: Worry level is divided by 3 after monkey inspects item.
	// Compute product of total inspections by two most active monkeys after 20 rounds.
	ins := make([]int, len(monkeys))
	for round := 0; round < 20; round++ {
		for i, m := range monkeys {
			ins[i] += len(m.items)
			m.act(monkeys, 3 /* div */, 0 /* mod */)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ins)))
	fmt.Println(int64(ins[0]) * int64(ins[1]))

	// Part 2: Worry level doesn't decrease and 10000 rounds.
	tests := make([]int, len(monkeys))
	for i := range monkeys {
		monkeys[i].items = orig[i]
		tests[i] = monkeys[i].test
	}
	lcm := lib.LCM(tests...)
	ins = make([]int, len(monkeys))
	for round := 0; round < 10000; round++ {
		for i, m := range monkeys {
			ins[i] += len(m.items)
			m.act(monkeys, 0 /* div */, lcm /* mod */)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ins)))
	fmt.Println(int64(ins[0]) * int64(ins[1]))
}

type monkey struct {
	id      int
	items   []int  // worry levels of items
	op      byte   // '*' or '+'
	operand string // int or "old"
	test    int    // check if worry level is evenly divided by this
	dt, df  int    // dest if test is true/false
}

// Parses lines like the following:
//
//  Monkey 0:
//    Starting items: 79, 98
//    Operation: new = old * 19
//    Test: divisible by 23
//  	If true: throw to monkey 2
//  	If false: throw to monkey 3
func newMonkey(lines []string) *monkey {
	lib.AssertEq(len(lines), 6)
	var m monkey
	lib.Extract(lines[0], `^Monkey (\d+):$`, &m.id)
	m.items = lib.ExtractInts(lines[1]) // e.g. "Starting items: 79, 98"
	lib.Extract(lines[2], `^  Operation: new = old ([+*]) (.+)$`, &m.op, &m.operand)
	lib.Extract(lines[3], `^  Test: divisible by (\d+)$`, &m.test)
	lib.Extract(lines[4], `^    If true: throw to monkey (\d+)$`, &m.dt)
	lib.Extract(lines[5], `^    If false: throw to monkey (\d+)$`, &m.df)
	return &m
}

func (m *monkey) act(monkeys []*monkey, div, mod int) {
	for _, item := range m.items {
		// Monkey inspects the item, increasing worry level.
		var operand int
		if m.operand == "old" {
			operand = item
		} else {
			var err error
			operand, err = strconv.Atoi(m.operand)
			lib.Assert(err == nil)
		}
		switch m.op {
		case '+':
			item += operand
		case '*':
			item *= operand
		default:
			lib.Panicf("Invalid operand %q", m.op)
		}

		// Monkey gets bored with the item, decreasing worry level.
		if div != 0 {
			item /= div
		}
		if mod != 0 {
			item %= mod
		}

		// Monkey tests and throws the item.
		var dst = m.df
		if item%m.test == 0 {
			dst = m.dt
		}
		lib.Assert(dst != m.id)
		monkeys[dst].items = append(monkeys[dst].items, item)
	}
	m.items = nil
}
