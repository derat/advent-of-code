package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	changes := lib.InputInts("2018/1")
	fmt.Println(lib.Sum(changes...))
}
