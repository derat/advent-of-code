package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	regs := make(map[string]int)
	var max int // part 2
	for _, ln := range lib.InputLines("2017/8") {
		var dr, op, cr, co string
		var amt, cv int
		lib.Extract(ln, `^(\w+) (inc|dec) (-?\d+) if (\w+) (<|<=|==|>=|>|!=) (-?\d+)$`,
			&dr, &op, &amt, &cr, &co, &cv)
		if cond(co, regs[cr], cv) {
			regs[dr] += lib.If(op == "inc", amt, -amt)
			max = lib.Max(max, regs[dr])
		}
	}
	fmt.Println(lib.Max(lib.MapIntVals(regs)...))
	fmt.Println(max)
}

func cond(op string, v1, v2 int) bool {
	switch op {
	case "<":
		return v1 < v2
	case "<=":
		return v1 <= v2
	case "==":
		return v1 == v2
	case ">=":
		return v1 >= v2
	case ">":
		return v1 > v2
	case "!=":
		return v1 != v2
	default:
		lib.Panicf("Invalid op %q", op)
		return false
	}
}
