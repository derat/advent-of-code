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
	sc := bufio.NewScanner(os.Stdin)

	// Read rules.
	var rules []rule
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			break
		}

		col := strings.IndexRune(ln, ':')
		if col == -1 {
			log.Fatalf("Missing colon in line %q", ln)
		}

		num, err := strconv.Atoi(ln[:col])
		if err != nil {
			log.Fatalf("Bad rule number in line %q: %v", ln, err)
		}

		var rl rule
		rest := strings.TrimSpace(ln[col+1:])
		if len(rest) == 3 && rest[0] == '"' && rest[2] == '"' {
			rl.ch = rune(rest[1])
		} else {
			for _, part := range strings.Split(rest, "|") {
				var list []int
				for _, s := range strings.Fields(part) {
					v, err := strconv.Atoi(s)
					if err != nil {
						log.Fatalf("Bad subrule number in line %q: %v", ln, err)
					}
					list = append(list, v)
				}
				rl.opts = append(rl.opts, list)
			}
			if len(rl.opts) == 0 {
				log.Fatalf("Bad subrule list in %q", ln)
			}
		}

		if num >= len(rules) {
			nr := make([]rule, num+1)
			copy(nr, rules)
			rules = nr
		}
		rules[num] = rl
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	// Read messages.
	var msgs []string
	for sc.Scan() {
		ln := strings.TrimSpace(sc.Text())
		if ln == "" {
			continue
		}
		msgs = append(msgs, ln)
	}
	if sc.Err() != nil {
		log.Fatal(sc.Err())
	}

	// Part 1:
	cnt := 0
	for _, msg := range msgs {
		if valid, err := check(msg, rules); err != nil {
			log.Fatalf("Failed checking %q: %v", msg, err)
		} else if valid {
			cnt++
		}
	}
	fmt.Println(cnt)

	// Part 2:
	rules[8] = rule{opts: [][]int{{42}, {42, 8}}}
	rules[11] = rule{opts: [][]int{{42, 31}, {42, 11, 31}}}
	cnt = 0
	for _, msg := range msgs {
		if valid, err := check(msg, rules); err != nil {
			log.Fatalf("Failed checking %q: %v", msg, err)
		} else if valid {
			cnt++
		}
	}
	fmt.Println(cnt)
}

type rule struct {
	opts [][]int // rule lists for alternation
	ch   rune
}

func check(msg string, rules []rule) (bool, error) {
	seqs := [][]int{{0}} // valid sequences

	for i, ch := range msg {
		// We ran out of valid sequences before the end of the message.
		if len(seqs) == 0 {
			return false, nil
		}

		var next [][]int // possible sequences for next char
		for len(seqs) > 0 {
			var exp [][]int // expanded sequences still under consideration for this char
			for _, s := range seqs {
				first := rules[s[0]]
				switch {
				case first.ch != 0:
					if ch == first.ch {
						// If we matched the last character in the rule, we're done.
						if len(s) == 1 && i == len(msg)-1 {
							return true, nil
						}
						if len(s) > 1 {
							next = append(next, s[1:])
						}
					}
				case len(first.opts) > 0:
					// Expand the options from the first rule's sequences.
					for _, o := range first.opts {
						exp = append(exp, append(o[:], s[1:]...))
					}
				default:
					return false, fmt.Errorf("Rule %d contains neither char nor rules", s[0])
				}
			}
			seqs = exp
		}
		seqs = next
	}

	// Almost valid, but the last character was wrong.
	return false, nil
}
