package main

import (
	"bytes"
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var ins []instr
	for _, ln := range lib.InputLines("2016/21") {
		var in instr
		var arg string
		switch {
		case lib.ExtractMaybe(ln, `^(swap position) (\d+) with position (\d+)$`, &in.op, &in.p1, &in.p2):
		case lib.ExtractMaybe(ln, `^(swap letter) ([a-z]) with letter ([a-z])$`, &in.op, &in.l1, &in.l2):
		case lib.ExtractMaybe(ln, `^(rotate) (left|right) (\d+) steps?$`, &in.op, &arg, &in.p1):
			in.p1 *= lib.If(arg == "left", -1, 1)
		case lib.ExtractMaybe(ln, `^(rotate) based on position of letter ([a-z])$`, &in.op, &in.l1):
		case lib.ExtractMaybe(ln, `^(reverse) positions (\d+) through (\d+)$`, &in.op, &in.p1, &in.p2):
		case lib.ExtractMaybe(ln, `^(move) position (\d+) to position (\d+)$`, &in.op, &in.p1, &in.p2):
		default:
			lib.Assertf(false, "Bad line %q", ln)
		}
		ins = append(ins, in)
	}

	const init = "abcdefgh"
	pw := []byte(init)
	scramble(pw, ins)
	fmt.Println(string(pw))
}

type instr struct {
	op     string
	p1, p2 int  // positions
	l1, l2 byte // letters
}

func scramble(b []byte, ins []instr) {
	for _, in := range ins {
		switch in.op {
		case "swap position":
			b[in.p1], b[in.p2] = b[in.p2], b[in.p1]
		case "swap letter":
			for i, ch := range b {
				if ch == in.l1 {
					b[i] = in.l2
				} else if ch == in.l2 {
					b[i] = in.l1
				}
			}
		case "rotate":
			amt := in.p1
			if in.l1 != 0 {
				idx := bytes.IndexByte(b, in.l1)
				lib.Assertf(idx != -1, "Failed to find %q in %q", in.l1, b)
				// "Once the index is determined, rotate the string to the right one time, plus a
				// number of times equal to that index, plus one additional time if the index was at
				// least 4."
				amt = 1 + idx + lib.If(idx >= 4, 1, 0)
			}
			lib.RotateSlice(b, amt)
		case "reverse":
			for i, j := in.p1, in.p2; i < j; i, j = i+1, j-1 {
				b[i], b[j] = b[j], b[i]
			}
		case "move":
			lib.Move(b, in.p1, in.p1+1, in.p2)
		}
	}
}
