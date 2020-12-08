package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var re = regexp.MustCompile(`^(\d+) (.+) bag(?:s?)(?:\.?)$`)

type bagInfo struct {
	color string
	num   int
}

func main() {
	holders := make(map[string][]string)
	bags := make(map[string][]bagInfo)

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		s := sc.Text()
		ps := strings.Split(s, " bags contain ")
		if len(ps) != 2 {
			log.Fatalf("bad line %q", s)
		}
		outer, lst := ps[0], ps[1]
		if lst == "no other bags." {
			continue
		}
		for _, p := range strings.Split(lst, ", ") {
			m := re.FindStringSubmatch(p)
			if m == nil {
				log.Fatalf("failed parsing %q in %q", p, s)
			}
			cnt, _ := strconv.Atoi(m[1])
			inner := m[2]
			holders[inner] = append(holders[inner], outer)
			bags[outer] = append(bags[outer], bagInfo{color: inner, num: cnt})
		}
	}
	if sc.Err() != nil {
		panic(sc.Err())
	}

	seen := make(map[string]struct{})
	var add func(col string)
	add = func(col string) {
		for _, c := range holders[col] {
			if _, ok := seen[c]; !ok {
				seen[c] = struct{}{}
				add(c)
			}
		}
	}
	add("shiny gold")
	println(len(seen))

	var count func(col string) int
	count = func(col string) int {
		total := 0
		for _, b := range bags[col] {
			total += b.num * (1 + count(b.color))
		}
		return total
	}
	println(count("shiny gold"))
}
