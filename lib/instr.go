package lib

import (
	"fmt"
	"strconv"
)

// NewInstr parses an instruction from ln.
// rmin and rmax specify the start and end characters for registers.
// ops maps from opcode to a regular expression matching the operation,
// with 0, 1, or 2 subgroups used to extract arguments.
func NewInstr(ln string, rmin, rmax byte, ops map[uint8]string) Instr {
	in := Instr{r: [2]int8{-1, -1}}
	for op, re := range ops {
		nargs := getRegexp(re).NumSubexp()
		args := make([]string, 2)
		argps := []interface{}{&args[0], &args[1]}
		if _, ok := ExtractMaybe(ln, re, argps[:nargs]...); !ok {
			continue
		}
		in.Op = op
		for i := 0; i < nargs; i++ {
			arg := args[i]
			if len(arg) == 1 && arg[0] >= rmin && arg[0] <= rmax {
				in.r[i] = int8(arg[0] - rmin)
			} else {
				var err error
				if in.v[i], err = strconv.ParseInt(arg, 10, 64); err != nil {
					Panicf("Failed converting numeric arg: %v", err)
				}
				in.r[i] = -1
			}
		}
		return in
	}
	panic(fmt.Sprintf("Invalid instruction %q", ln))
}

// Instr represents a single instruction.
type Instr struct {
	Op uint8
	r  [2]int8 // -1 if unset
	v  [2]int64
}

// Val returns the value of in's i-th argument.
func (in *Instr) Val(i int, regs []int64) int64 {
	if in.r[i] >= 0 {
		return regs[in.r[i]]
	}
	return in.v[i]
}

// Ptr returns a pointer to in's i-th argument, which must be a register.
func (in *Instr) Ptr(i int, regs []int64) *int64 {
	if in.r[i] < 0 {
		Panicf("Can't get pointer to non-register arg")
	}
	return &(regs[in.r[i]])
}
