package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputInt64s("2019/5")

	// Part 1: Run diagnostic program with input of 1 and print diagnostic code
	// (the final output after a stream of 0s).
	vm := lib.NewIntcode(input)
	go func() { vm.In <- 1 }()
	vm.Start()
	var last int64
	for v := range vm.Out {
		last = v
	}
	lib.Assert(vm.Wait())
	fmt.Println(last)

	// Part 2: Provide 5 as input and print diagnostic code (only output).
	vm = lib.NewIntcode(input)
	done := make(chan struct{})
	go func() {
		vm.In <- 5
		fmt.Println(<-vm.Out)
		close(done)
	}()
	lib.Assert(vm.Run())
	<-done
}
