package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	monkeys := make(map[string]*monkey)
	const namePat = `[a-z]{4}`
	for _, ln := range lib.InputLines("2022/21") {
		var name string
		var m monkey
		switch {
		case lib.TryExtract(ln, `^(`+namePat+`): (`+namePat+`) ([-+*/]) (`+namePat+`)`,
			&name, &m.lhs, &m.op, &m.rhs):
		case lib.TryExtract(ln, `^(`+namePat+`): (\d+)$`, &name, &m.val):
		default:
			lib.Panicf("Bad line %q", ln)
		}
		monkeys[name] = &m
	}

	fmt.Println(monkeys["root"].yell(monkeys))
}

type monkey struct {
	val      int64  // 0 if not yet computed
	lhs, rhs string // monkey names
	op       byte   // '+', '-', '*', '/', or 0 for literal
}

func (m *monkey) yell(all map[string]*monkey) int64 {
	if m.val != 0 {
		return m.val // literal or already cached
	}
	lhs := all[m.lhs].yell(all)
	rhs := all[m.rhs].yell(all)
	switch m.op {
	case '+':
		m.val = lhs + rhs
	case '-':
		m.val = lhs - rhs
	case '*':
		m.val = lhs * rhs
	case '/':
		m.val = lhs / rhs
	default:
		lib.Panicf("Invalid operator %q", m.op)
	}
	return m.val
}
