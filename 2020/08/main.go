package main

import (
	"log"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []instr
	for _, ln := range lib.ReadLines() {
		var in instr
		lib.Parse(ln, `^(acc|jmp|nop) ([-+]\d+)$`, (*string)(&in.op), &in.val)
		ins = append(ins, in)
	}

	r, accum := run(ins)
	if r != loop {
		log.Fatal("didn't loop")
	}
	println(accum)

	swap := func(o op) op {
		if o == jmp {
			return nop
		} else if o == nop {
			return jmp
		}
		return o
	}

	for i := 0; i < len(ins); i++ {
		in := &ins[i]
		if in.op == acc {
			continue
		}
		in.op = swap(in.op) // swap jmp and nop
		r, accum := run(ins)
		in.op = swap(in.op) // swap back
		if r == ok {
			println(accum)
			break
		}
	}
}

type op string

const (
	acc op = "acc"
	jmp    = "jmp"
	nop    = "nop"
)

type instr struct {
	op  op
	val int
}

type res int

const (
	ok res = iota
	loop
	segv
)

func run(ins []instr) (res, int) {
	var accum, ip int
	seen := make(map[int]struct{})
	for {
		if ip == len(ins) {
			return ok, accum
		} else if ip < 0 || ip > len(ins) {
			return segv, accum
		}

		if _, ok := seen[ip]; ok {
			return loop, accum
		}
		seen[ip] = struct{}{}

		in := &ins[ip]
		//fmt.Printf("%d %s %d\n", ip, in.op, in.val)
		switch in.op {
		case acc:
			accum += in.val
			ip++
		case jmp:
			ip += in.val
		case nop:
			ip++
		}
	}
}
