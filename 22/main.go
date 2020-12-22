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
	var d1, d2 []int

	var recip *[]int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		switch ln {
		case "":
		case "Player 1:":
			recip = &d1
		case "Player 2:":
			recip = &d2
		default:
			v, err := strconv.Atoi(ln)
			if err != nil {
				log.Fatalf("Bad card %q: %v", ln, err)
			}
			*recip = append(*recip, v)
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	cnt := 0
	for len(d1) > 0 && len(d2) > 0 {
		cnt++
		v1, v2 := d1[0], d2[0]
		d1, d2 = d1[1:], d2[1:]
		if v1 > v2 {
			d1 = append(d1, v1, v2)
		} else {
			d2 = append(d2, v2, v1)
		}
	}

	// Part 1:
	var wd *[]int
	if len(d1) > 0 {
		wd = &d1
	} else {
		wd = &d2
	}
	score := 0
	for i, v := range *wd {
		score += v * (len(*wd) - i)
	}
	fmt.Println(score)
}
