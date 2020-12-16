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
	rules := make(map[string]rule) // keyed by field name
	var yours ticket
	var nearby []ticket

	sect := rulesSect
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		switch ln {
		case "": // ignore blank lines
		case "your ticket:":
			sect = yoursSect
		case "nearby tickets:":
			sect = nearbySect
		default:
			switch sect {
			case rulesSect:
				parts := strings.Split(ln, ": ")
				if len(parts) != 2 {
					log.Fatalf("bad rules line %q", ln)
				}
				rule, err := newRule(parts[1])
				if err != nil {
					log.Fatalf("bad rule in line %q: %v", ln, err)
				}
				rules[parts[0]] = rule
			case yoursSect:
				var err error
				if yours, err = newTicket(ln); err != nil {
					log.Fatalf("bad ticket line %q: %v", ln, err)
				}
			case nearbySect:
				if t, err := newTicket(ln); err != nil {
					log.Fatalf("bad ticket line %q: %v", ln, err)
				} else {
					nearby = append(nearby, t)
				}
			}
		}
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	errRate := 0
	var valid []ticket
	for _, t := range nearby {
		tval := true
		for _, f := range t.fields {
			fval := false
			for _, r := range rules {
				if r.valid(f) {
					fval = true
					break
				}
			}
			if !fval {
				tval = false
				errRate += f
			}
		}
		if tval {
			valid = append(valid, t)
		}
	}
	fmt.Println(errRate)

	// Make a map from field name to possible indexes.
	poss := make(map[string]map[int]struct{}, len(rules))
	for n := range rules {
		fm := make(map[int]struct{})
		for i := range yours.fields {
			fm[i] = struct{}{}
		}
		poss[n] = fm
	}

	// Drop indexes with values that are out of range.
	for _, t := range valid {
		for i, f := range t.fields {
			for n, r := range rules {
				if !r.valid(f) {
					delete(poss[n], i)
				}
			}
		}
	}

	// Determine the mapping from field name to index.
	indexes := make(map[string]int, len(rules))
	for len(poss) != 0 {
		plen := len(poss)
		for n, p := range poss {
			if len(p) == 1 {
				var v int
				for v = range p {
				}
				indexes[n] = v
				for o := range poss {
					delete(poss[o], v)
				}
				delete(poss, n)
				break
			}
		}
		if len(poss) == plen {
			log.Fatal("didn't find an index")
		}
	}

	prod := 1
	for n, i := range indexes {
		if strings.HasPrefix(n, "departure") {
			prod *= yours.fields[i]
		}
	}
	fmt.Println(prod)
}

type sect int

const (
	rulesSect sect = iota
	yoursSect
	nearbySect
)

type rule struct {
	ranges [][]int
}

func newRule(s string) (rule, error) {
	var r rule
	for _, p := range strings.Split(s, " or ") {
		parts := strings.Split(p, "-")
		if len(parts) != 2 {
			return r, fmt.Errorf("bad range %q", p)
		}
		min, err := strconv.Atoi(parts[0])
		if err != nil {
			return r, err
		}
		max, err := strconv.Atoi(parts[1])
		if err != nil {
			return r, err
		}
		r.ranges = append(r.ranges, []int{min, max})
	}
	return r, nil
}

func (r *rule) valid(n int) bool {
	for _, ra := range r.ranges {
		if n >= ra[0] && n <= ra[1] {
			return true
		}
	}
	return false
}

type ticket struct {
	fields []int
}

func newTicket(s string) (ticket, error) {
	var t ticket
	for _, f := range strings.Split(s, ",") {
		v, err := strconv.Atoi(f)
		if err != nil {
			return t, err
		}
		t.fields = append(t.fields, v)
	}
	return t, nil
}
