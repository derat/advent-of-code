package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	monkeys := make(map[string]*monkey)
	children := make(map[*monkey][2]string)
	const namePat = `[a-z]{4}`
	for _, ln := range lib.InputLines("2022/21") {
		var m monkey
		var lhs, rhs string
		switch {
		case lib.TryExtract(ln, `^(`+namePat+`): (`+namePat+`) ([-+*/]) (`+namePat+`)`,
			&m.name, &lhs, &m.op, &rhs):
			children[&m] = [2]string{lhs, rhs}
		case lib.TryExtract(ln, `^(`+namePat+`): (\d+)$`, &m.name, &m.val):
		default:
			lib.Panicf("Bad line %q", ln)
		}
		monkeys[m.name] = &m
	}
	for m, ch := range children {
		m.lhs = monkeys[ch[0]]
		m.rhs = monkeys[ch[1]]
	}

	fmt.Println(monkeys["root"].yell())
}

type monkey struct {
	name     string
	val      int64 // 0 if not yet computed
	lhs, rhs *monkey
	op       byte // '+', '-', '*', '/', or 0 for literal
}

func (m *monkey) yell() int64 {
	if m.val != 0 {
		return m.val // literal or already cached
	}
	lhs, rhs := m.lhs.yell(), m.rhs.yell()
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
