package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.InputParagraphs("2017/25")

	var cur string
	var steps int
	lib.Extract(pgs[0][0], `^Begin in state ([A-Z])\.$`, &cur)
	lib.Extract(pgs[0][1], `^Perform a diagnostic checksum after (\d+) steps\.$`, &steps)

	// Each state description looks like the following:
	//
	//  In state D:
	//    If the current value is 0:
	//  	- Write the value 1.
	//  	- Move one slot to the left.
	//  	- Continue with state E.
	//    If the current value is 1:
	//  	- Write the value 0.
	//  	- Move one slot to the left.
	//  	- Continue with state F.
	states := make(map[string]state)
	for _, pg := range pgs[1:] {
		var name, dir0, dir1 string
		var st state
		lib.Extract(strings.Join(pg, " "), `^In state ([A-Z]):\s+`+
			`If the current value is 0:\s+`+
			`- Write the value (0|1)\.\s+`+
			`- Move one slot to the (left|right)\.\s+`+
			`- Continue with state ([A-Z])\.\s+`+
			`If the current value is 1:\s+`+
			`- Write the value (0|1)\.\s+`+
			`- Move one slot to the (left|right)\.\s+`+
			`- Continue with state ([A-Z])\.$`,
			&name, &st.zero.write, &dir0, &st.zero.next,
			&st.one.write, &dir1, &st.one.next)
		st.zero.off = lib.If(dir0 == "left", -1, 1)
		st.one.off = lib.If(dir1 == "left", -1, 1)
		states[name] = st
	}

	// I was expecting to need to do some sort of analysis,
	// e.g. detecting loops to be able to predict the final count
	// without needing execute all steps. Running them all only
	// takes a few seconds, though.
	tape := make(map[int]int)
	var pos int
	for step := 0; step < steps; step++ {
		var tr trans
		if tape[pos] == 0 {
			tr = states[cur].zero
		} else {
			tr = states[cur].one
		}
		if tr.write == 1 {
			tape[pos] = 1
		} else {
			delete(tape, pos)
		}
		pos += tr.off
		cur = tr.next
	}
	fmt.Println(len(tape))
}

type state struct {
	zero, one trans
}

type trans struct {
	write int    // 0 or 1
	off   int    // -1 or 1
	next  string // e.g. "B"
}
