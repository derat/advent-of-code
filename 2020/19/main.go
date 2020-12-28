package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	pgs := lib.ReadParagraphs()
	lib.AssertEq(len(pgs), 2)

	// Parse rules.
	var rules []rule
	for _, ln := range pgs[0] {
		var rl rule
		var num int
		var rest string
		lib.Parse(ln, `^(\d+): (.+)$`, &num, &rest)
		if len(rest) == 3 && rest[0] == '"' && rest[2] == '"' {
			rl.ch = rune(rest[1])
		} else {
			for _, part := range strings.Split(rest, "|") {
				rl.opts = append(rl.opts, lib.ExtractInts(part))
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

	// The second paragraph contains messages to validate.
	msgs := pgs[1]

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
