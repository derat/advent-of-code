package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

type input struct {
	op, lhs, rhs string // rhs unused for SET, NOT
}

func main() {
	inputs := make(map[string]input)   // keyed by dest wire id
	signals := make(map[string]uint16) // final wire values keyed by id
	deps := make(map[string][]string)  // wire id to ids depending on it
	check := make(map[string]struct{}) // newly-resolved wires

	// Records src depending on dst.
	// If src is a literal value, also adds it to signals and check.
	addDep := func(src, dst string) {
		deps[src] = append(deps[src], dst)
		if v, err := strconv.Atoi(src); err == nil {
			signals[src] = uint16(v)
			check[src] = struct{}{}
		}
	}

	for _, ln := range lib.InputLines("2015/7") {
		switch {
		case strings.HasPrefix(ln, "NOT"):
			var src, dst string
			lib.Extract(ln, `^NOT (\w+) -> ([a-z]+)$`, &src, &dst)
			inputs[dst] = input{"NOT", src, ""}
			addDep(src, dst)
		case strings.Contains(ln, "AND") || strings.Contains(ln, "OR") || strings.Contains(ln, "SHIFT"):
			var lhs, op, rhs, dst string
			lib.Extract(ln, `^(\w+) (AND|OR|LSHIFT|RSHIFT) (\w+) -> ([a-z]+)$`, &lhs, &op, &rhs, &dst)
			inputs[dst] = input{op, lhs, rhs}
			addDep(lhs, dst)
			addDep(rhs, dst)
		default:
			var src, dst string
			lib.Extract(ln, `^(\w+) -> (\w+)$`, &src, &dst)
			inputs[dst] = input{"SET", src, ""}
			addDep(src, dst)
		}
	}

	// Returns in's value, or false if it can't be computed yet.
	compute := func(in input) (uint16, bool) {
		lhs, lok := signals[in.lhs]
		rhs, rok := signals[in.rhs]
		switch in.op {
		case "SET":
			if lok {
				return lhs, true
			}
			return 0, false
		case "AND":
			if lok && rok {
				return lhs & rhs, true
			}
			return 0, false
		case "OR":
			if lok && rok {
				return lhs | rhs, true
			}
			return 0, false
		case "NOT":
			if lok {
				return ^lhs, true
			}
			return 0, false
		case "LSHIFT":
			if lok && rok {
				return lhs << int(rhs), true
			}
			return 0, false
		case "RSHIFT":
			if lok && rok {
				return lhs >> int(rhs), true
			}
			return 0, false
		default:
			panic(fmt.Sprintf("Invalid op %q", in.op))
		}
	}

	for len(check) > 0 {
		newCheck := make(map[string]struct{})
		for src := range check {
			for _, dst := range deps[src] {
				if _, ok := signals[dst]; ok {
					continue // already computed it
				}
				in := inputs[dst]
				if v, ok := compute(in); ok {
					signals[dst] = v
					newCheck[dst] = struct{}{}
				}
			}
		}
		check = newCheck
	}

	fmt.Println(signals["a"])
}
