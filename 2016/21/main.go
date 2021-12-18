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
		case lib.TryExtract(ln, `^(swap position) (\d+) with position (\d+)$`, &in.op, &in.p1, &in.p2):
		case lib.TryExtract(ln, `^(swap letter) ([a-z]) with letter ([a-z])$`, &in.op, &in.l1, &in.l2):
		case lib.TryExtract(ln, `^(rotate) (left|right) (\d+) steps?$`, &in.op, &arg, &in.p1):
			in.p1 *= lib.If(arg == "left", -1, 1)
		case lib.TryExtract(ln, `^(rotate) based on position of letter ([a-z])$`, &in.op, &in.l1):
		case lib.TryExtract(ln, `^(reverse) positions (\d+) through (\d+)$`, &in.op, &in.p1, &in.p2):
		case lib.TryExtract(ln, `^(move) position (\d+) to position (\d+)$`, &in.op, &in.p1, &in.p2):
		default:
			lib.Assertf(false, "Bad line %q", ln)
		}
		ins = append(ins, in)
	}

	// Part 1: scramble
	pw := []byte("abcdefgh")
	for _, in := range ins {
		in.do(pw)
	}
	fmt.Println(string(pw))

	// Part 2: unscramble
	pw2 := []byte("fbgdceah")
	for i := len(ins) - 1; i >= 0; i-- {
		ins[i].undo(pw2)
	}
	fmt.Println(string(pw2))
}

type instr struct {
	op     string
	p1, p2 int  // positions
	l1, l2 byte // letters
}

func (in *instr) do(b []byte) {
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
	default:
		lib.Assertf(false, "Bad op %q", in.op)
	}
}

func (in *instr) undo(b []byte) {
	switch in.op {
	case "swap position", "swap letter", "reverse":
		in.do(b) // same in either direction
	case "rotate":
		if in.l1 == 0 {
			lib.RotateSlice(b, -in.p1) // rotate in reverse direction
		} else {
			// This is the only tricky part. We have the position of the letter after the
			// forward operation was performed, so we need to figure out where it was initially
			// to undo the operation by rotating back.
			//
			//  index  rotate  end
			//      0       1    1
			//      1       2    3
			//      2       3    5
			//      3       4    7
			//      4       6    10 (2)
			//      5       7    12 (4)
			//      6       8    14 (6)
			//      7       9    16 (0)
			//
			// For numbers less than 4, we rotate by idx+1, so the resulting index (2*idx+1) will be odd.
			// Otherwise, we rotate by idx+2, so the resulting index (2*idx+2) will be even.
			idx := bytes.IndexByte(b, in.l1)
			lib.Assertf(idx != -1, "Failed to find %q in %q", in.l1, b)
			orig := lib.If(idx%2 == 1, (idx-1)/2, (idx+len(b)-2)/2)
			lib.RotateSlice(b, orig-idx)
		}
	case "move":
		rev := instr{op: in.op, p1: in.p2, p2: in.p1}
		rev.do(b) // swap args
	default:
		lib.Assertf(false, "Bad op %q", in.op)
	}
}
