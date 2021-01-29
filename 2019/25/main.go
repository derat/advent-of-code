package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	vm := lib.NewIntcode(lib.InputInt64s("2019/25"))

	vm.In = make(chan int64, 2048) // characters read from stdin
	go func(stdin io.Reader) {
		r := bufio.NewReader(stdin)
		for {
			cmd, err := r.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				lib.Panicf("Input failed: %v", err)
			}
			cmd = strings.TrimSpace(cmd)

			switch {
			case cmd == "north" || cmd == "south" || cmd == "east" || cmd == "west" || cmd == "inv":
			case lib.ExtractMaybe(cmd, `^(?:drop|take) .+$`):
			case cmd == "combos":
				cmd = strings.TrimSpace(combos())
			case cmd == "yoink":
				cmd = strings.TrimSpace(yoink)
			default:
				fmt.Printf("Invalid command %q\n", cmd)
				continue
			}

			for _, ch := range cmd {
				vm.In <- int64(ch)
			}
			vm.In <- '\n'
		}
		vm.Halt()
	}(os.Stdin)

	done := make(chan struct{}) // closed when program halts
	go func() {
		for v := range vm.Out {
			fmt.Print(string(rune(v)))
		}
		close(done)
	}()

	vm.Start()
	<-done
	lib.Assert(vm.Wait())
}

// yoink runs commands to grab all non-leathal items and navigate
// to Security Checkpoint. It should be run from Hull Breach.
const yoink = `
east
take loom
east
take fixed point
north
take spool of cat6
west
take shell
east
north
take weather machine
south
south
west
south
take ornament
west
north
take candy cane
south
east
north
west
north
take wreath
north
east`

// combos runs commands for trying all combinations of items.
// All of the items should be initially in your inventory (i.e.
// after running the "yoink" command), and it should run in
// Security Checkpoint (directly north of Pressure-Sensitive
// Floor in my input).
func combos() string {
	items := []string{
		"candy cane",
		"fixed point",
		"loom",
		"ornament",
		"shell",
		"spool of cat6",
		"weather machine",
		"wreath",
	}

	var cmds, have []string
	for _, s := range items {
		cmds = append(cmds, "drop "+s)
	}
	for i := 0; i < 1<<len(items); i++ {
		for _, s := range have {
			cmds = append(cmds, "drop "+s)
		}
		have = nil
		for j, s := range items {
			if i&(1<<j) != 0 {
				cmds = append(cmds, "take "+s)
				have = append(have, s)
			}
		}
		cmds = append(cmds, "south")
	}
	return strings.Join(cmds, "\n")
}
