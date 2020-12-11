package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var ins []instr
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		in, err := newInstr(sc.Text())
		if err != nil {
			log.Fatalf("bad line %q: %v", sc.Text(), err)
		}
		ins = append(ins, in)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
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

type op int

const (
	acc op = iota
	jmp
	nop
)

func (o op) String() string {
	switch o {
	case acc:
		return "acc"
	case jmp:
		return "jmp"
	case nop:
		return "nop"
	}
	return fmt.Sprintf("[%d]", o)
}

type instr struct {
	op  op
	val int
}

func newInstr(ln string) (instr, error) {
	ps := strings.Fields(ln)
	if len(ps) != 2 {
		return instr{}, fmt.Errorf("%v part(s) instead of 2", len(ps))
	}

	var in instr
	switch ps[0] {
	case "acc":
		in.op = acc
	case "jmp":
		in.op = jmp
	case "nop":
		in.op = nop
	default:
		return in, fmt.Errorf("invalid op %q", ps[0])
	}

	var err error
	if in.val, err = strconv.Atoi(ps[1]); err != nil {
		return in, err
	}

	return in, nil
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
