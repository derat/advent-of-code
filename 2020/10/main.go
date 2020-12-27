package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	var outs []int
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Fatal(err)
		}
		outs = append(outs, v)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}
	sort.Ints(outs)
	outs = append(outs, outs[len(outs)-1]+3)

	diffs := make(map[int]int)
	last := 0
	for _, v := range outs {
		d := v - last
		if d < 1 || d > 3 {
			log.Fatalf("bad diff %v between %v and %v", d, last, v)
		}
		diffs[d]++
		last = v
	}
	fmt.Printf("%v * %v = %v\n", diffs[1], diffs[3], diffs[1]*diffs[3])

	paths := map[int]int64{0: 1}
	for _, v := range outs {
		paths[v] = paths[v-3] + paths[v-2] + paths[v-1]
	}
	fmt.Println(paths[outs[len(outs)-1]])
}
