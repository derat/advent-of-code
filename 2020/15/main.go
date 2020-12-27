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
	const (
		end  = 2020
		end2 = 30000000
	)

	r := bufio.NewReader(os.Stdin)
	s, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	turn := 1
	seen := make(map[int]int) // turn in which number last spoken
	last := -1                // last number spoken
	for _, ns := range strings.Split(strings.TrimSpace(s), ",") {
		n, err := strconv.Atoi(ns)
		if err != nil {
			log.Fatal(err)
		}
		last = n
		seen[n] = turn
		turn++
	}

	next := 0 // next number to speak (i.e. age of last number)
	for {
		last = next
		if lastTurn, ok := seen[last]; ok {
			next = turn - lastTurn
		} else {
			next = 0
		}
		seen[last] = turn

		if turn == end {
			fmt.Println(last)
		}
		if turn == end2 {
			fmt.Println(last)
			break
		}
		turn++
	}
}
