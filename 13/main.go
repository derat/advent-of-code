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
	est, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	var ids []int64
	if s, err = r.ReadString('\n'); err != nil {
		log.Fatal(err)
	}
	for _, p := range strings.Split(s[:len(s)-1], ",") {
		if p == "x" {
			ids = append(ids, 0)
		} else {
			n, err := strconv.ParseInt(p, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			ids = append(ids, n)
		}
	}

Loop:
	for i := est; true; i++ {
		for _, id := range ids {
			if id > 0 && i%id == 0 {
				wait := i - est
				fmt.Printf("%d * %d = %d\n", id, wait, id*wait)
				break Loop
			}
		}
	}

	// Find t such that (t+n) % ids[n] = 0 for all n where ids[n] != 0.
	// IDs seem to all be coprime to each other.

	// Finds the first x >= t such that (t+off)%mod == 0.
	// The supplied step is used as an increment.
	find := func(t, mod, off, step int64) int64 {
		for {
			if (t+off)%mod == 0 {
				return t
			}
			t += step
		}
	}

	// There's probably a simpler solution but I've spent too much time on this already. :-(
	t := ids[0]
	step := ids[0]
	for i := 1; i < len(ids); i++ {
		mod := ids[i]
		if mod == 0 {
			continue
		}
		t = find(t, mod, int64(i), step)
		// Since the IDs are coprime, we can multiply them together to find the period.
		step *= mod
	}
	fmt.Println(t)
}
