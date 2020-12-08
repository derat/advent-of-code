package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^(\d+) (.+) bag(?:s?)(?:\.?)$`)

func main() {
	holders := make(map[string][]string)
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
			inner := m[2]
			holders[inner] = append(holders[inner], outer)
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
}
