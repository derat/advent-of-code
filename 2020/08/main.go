package main

import (
	"fmt"
	"log"

	"github.com/derat/advent-of-code/lib"
)

const (
	acc = iota
	jmp
	nop
)

func main() {
	var ins []lib.Instr
	for _, ln := range lib.InputLines("2020/8") {
		ins = append(ins, lib.NewInstr(ln, 0, 0, map[uint8]string{
			acc: `^acc ([-+]\d+)$`,
			jmp: `^jmp ([-+]\d+)$`,
			nop: `^nop ([-+]\d+)$`,
		}))
	}

	r, accum := run(ins)
	if r != loop {
		log.Fatal("didn't loop")
	}
	fmt.Println(accum)

	swap := func(o uint8) uint8 {
		if o == jmp {
			return nop
		} else if o == nop {
			return jmp
		}
		return o
	}

	for i := 0; i < len(ins); i++ {
		in := &ins[i]
		if in.Op == acc {
			continue
		}
		in.Op = swap(in.Op) // swap jmp and nop
		r, accum := run(ins)
		in.Op = swap(in.Op) // swap back
		if r == ok {
			fmt.Println(accum)
			break
		}
	}
}

type res int

const (
	ok res = iota
	loop
	segv
)

func run(ins []lib.Instr) (res, int64) {
	var accum int64
	var ip int
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
		switch in.Op {
		case acc:
			accum += in.Val(0, nil)
			ip++
		case jmp:
			ip += int(in.Val(0, nil))
		case nop:
			ip++
		}
	}
}
