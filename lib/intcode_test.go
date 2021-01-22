package lib

import (
	"strconv"
	"strings"
	"testing"
)

func TestIntcode_Basic(t *testing.T) {
	for _, tc := range []struct {
		prog, mem string // orig program and final mem state
		in, out   string // input and output values
	}{
		// 2019/2 examples:
		{"1,9,10,3,2,3,11,0,99,30,40,50", "3500,9,10,70,2,3,11,0,99,30,40,50", "", ""},
		{"1,0,0,0,99", "2,0,0,0,99", "", ""},
		{"2,3,0,3,99", "2,3,0,6,99", "", ""},
		{"2,4,4,5,99,0", "2,4,4,5,99,9801", "", ""},
		{"1,1,1,4,99,5,6,0,99", "30,1,1,4,2,5,6,0,99", "", ""},

		// 2019/5 examples:
		{"3,0,99", "40,0,99", "40", ""},                 // output
		{"4,2,99", "4,2,99", "", "99"},                  // input
		{"1002,4,3,4,33", "1002,4,3,4,99", "", ""},      // multiply pos/immed
		{"1101,100,-1,4,0", "1101,100,-1,4,99", "", ""}, // negative
		// Jump test: output 1 if input is 8 and 0 otherwise (positional).
		{"3,9,8,9,10,9,4,9,99,-1,8", "", "7", "0"},
		{"3,9,8,9,10,9,4,9,99,-1,8", "", "8", "1"},
		{"3,9,8,9,10,9,4,9,99,-1,8", "", "9", "0"},
		// Jump test: output 1 if input less than 8 and 0 otherwise (positional).
		{"3,9,7,9,10,9,4,9,99,-1,8", "", "7", "1"},
		{"3,9,7,9,10,9,4,9,99,-1,8", "", "8", "0"},
		{"3,9,7,9,10,9,4,9,99,-1,8", "", "9", "0"},
		// Jump test: output 1 if input is 8 and 0 otherwise (immediate).
		{"3,3,1108,-1,8,3,4,3,99", "", "7", "0"},
		{"3,3,1108,-1,8,3,4,3,99", "", "8", "1"},
		{"3,3,1108,-1,8,3,4,3,99", "", "9", "0"},
		// Jump test: output 1 if input less than 8 and 0 otherwise (immediate).
		{"3,3,1107,-1,8,3,4,3,99", "", "7", "1"},
		{"3,3,1107,-1,8,3,4,3,99", "", "8", "0"},
		{"3,3,1107,-1,8,3,4,3,99", "", "9", "0"},
		// "... uses an input instruction to ask for a single number. The program will then output
		// 999 if the input value is below 8, output 1000 if the input value is equal to 8, or
		// output 1001 if the input value is greater than 8."
		{`3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,` +
			`1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,` +
			`999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99`, "", "7", "999"},
		{`3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,` +
			`1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,` +
			`999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99`, "", "8", "1000"},
		{`3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,` +
			`1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,` +
			`999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99`, "", "9", "1001"},
	} {
		prog := ExtractInt64s(tc.prog)
		vm := NewIntcode(prog)
		vm.Start()
		if tc.in != "" {
			for _, v := range ExtractInt64s(tc.in) {
				vm.In <- v
			}
		}
		var out []string
		for v := range vm.Out {
			out = append(out, strconv.FormatInt(v, 10))
		}
		if !vm.Wait() {
			t.Errorf("%q with input %q failed", tc.prog, tc.in)
			continue
		}

		if tc.mem != "" {
			mem := make([]string, len(prog))
			for i := range mem {
				mem[i] = strconv.FormatInt(vm.Mem[int64(i)], 10)
			}
			if got := strings.Join(mem, ","); got != tc.mem {
				t.Errorf("%q with input %q produced %q; want %q", tc.prog, tc.in, got, tc.mem)
			}
		}
		if got := strings.Join(out, ","); got != tc.out {
			t.Errorf("%q with input %q output %q; want %q", tc.prog, tc.in, got, tc.out)
		}
	}
}

func TestIntcode_Parallel(t *testing.T) {
	for _, tc := range []struct {
		prog, ins string
		feedback  bool
		want      int64
	}{
		// Day 7 examples:
		{"3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0", "4,3,2,1,0", false, 43210},
		{`3,23,3,24,1002,24,10,24,1002,23,-1,23,` +
			`101,5,23,23,1,24,23,23,4,23,99,0,0`, "0,1,2,3,4", false, 54321},
		{`3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,` +
			`1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0`, "1,0,4,3,2", false, 65210},
		{`3,26,1001,26,-4,26,3,27,1002,27,2,27,1,27,26,` +
			`27,4,27,1001,28,-1,28,1005,28,6,99,0,0,5`, "9,8,7,6,5", true, 139629729},
		{`3,52,1001,52,-5,52,3,53,1,52,56,54,1007,54,5,55,1005,55,26,1001,54,` +
			`-5,54,1105,1,12,1,53,54,53,1008,54,0,55,1001,55,1,55,2,53,55,53,4,` +
			`53,1001,56,-1,56,1005,56,6,99,0,0,0,0,10`, "9,7,8,5,6", true, 18216},
	} {
		// Create the VMs and wire them together.
		prog := ExtractInt64s(tc.prog)
		ins := ExtractInt64s(tc.ins)
		vms := make([]*Intcode, len(ins))
		for i := range ins {
			vms[i] = NewIntcode(prog)
			if i > 0 {
				vms[i].In = vms[i-1].Out
			}
		}
		if tc.feedback {
			vms[0].In = vms[len(vms)-1].Out
		}

		// Start the VMs and feed them their inputs.
		for i, vm := range vms {
			vm.Start()
			vm.In <- ins[i]
		}

		// Send the input signal to the first and read the output from the last.
		vms[0].In <- 0
		if !vms[0].Wait() {
			t.Errorf("%q with inputs %q failed", tc.prog, tc.ins)
		} else if got := <-vms[len(vms)-1].Out; got != tc.want {
			t.Errorf("%q with inputs %q produced %v; want %v", tc.prog, tc.ins, got, tc.want)
		}
	}
}
