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

	var accum, ip int
	seen := make(map[int]struct{})
	for {
		if ip < 0 || ip >= len(ins) {
			log.Fatalf("ip %v not in [0,%v)", ip, len(ins))
		}

		if _, ok := seen[ip]; ok {
			println(accum)
			break
		}
		seen[ip] = struct{}{}

		in := &ins[ip]
		fmt.Printf("%d %s %d\n", ip, in.op, in.val)
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
