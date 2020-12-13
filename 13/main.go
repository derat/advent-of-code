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
	r := bufio.NewReader(os.Stdin)

	s, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	est, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		log.Fatal(err)
	}

	var ids []int
	if s, err = r.ReadString('\n'); err != nil {
		log.Fatal(err)
	}
	for _, p := range strings.Split(s[:len(s)-1], ",") {
		if p == "x" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			log.Fatal(err)
		}
		ids = append(ids, n)
	}

Loop:
	for i := est; true; i++ {
		for _, id := range ids {
			if i%id == 0 {
				wait := i - est
				fmt.Printf("%d * %d = %d\n", id, wait, id*wait)
				break Loop
			}
		}
	}
}
