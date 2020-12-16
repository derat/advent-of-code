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
	rules := make(map[string]rule)
	var your ticket
	var nearby []ticket

	sect := rulesSect
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		switch ln {
		case "": // ignore blank lines
		case "your ticket:":
			sect = yourSect
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
			case yourSect:
				var err error
				if your, err = newTicket(ln); err != nil {
					log.Fatalf("bad ticket line %q: %v", ln, err)
				}
				_ = your // unused in part 1
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
	for _, t := range nearby {
		for _, f := range t.fields {
			valid := false
			for _, r := range rules {
				if r.valid(f) {
					valid = true
					break
				}
			}
			if !valid {
				errRate += f
			}
		}
	}
	fmt.Println(errRate)
}

type sect int

const (
	rulesSect sect = iota
	yourSect
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
