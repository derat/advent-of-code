package main

import (
	"regexp"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	input := lib.InputLines("2018/20")[0]
	lib.Assert(strings.HasPrefix(input, "^"))
	lib.Assert(strings.HasSuffix(input, "$"))
	input = input[1 : len(input)-1]

	// Parse the tokens and check that the resulting sequence matches the original regexp.
	tokens := lib.Tokenize(input, "(", ")", "|", regexp.MustCompile(`^[NSEW]+`))
	root, n := parse(tokens)
	lib.AssertEq(input, root.String())
	lib.AssertEq(n, len(tokens))
}

// seq describes a portion of a regular expression.
// Exactly one of its fields is non-nil.
type seq struct {
	str  *string // literal string, e.g. "NEEWS" or ""
	opts []*seq  // options in a group
	seqs []*seq  // consecutive sequences
}

func (s *seq) String() string {
	switch {
	case len(s.opts) > 0:
		lib.Assert(s.seqs == nil)
		lib.Assert(s.str == nil)
		str := "("
		for i, o := range s.opts {
			str += o.String()
			if i < len(s.opts)-1 {
				str += "|"
			}
		}
		return str + ")"
	case len(s.seqs) > 0:
		lib.Assert(s.opts == nil)
		lib.Assert(s.str == nil)
		var str string
		for _, o := range s.seqs {
			str += o.String()
		}
		return str
	case s.str != nil:
		lib.Assert(s.opts == nil)
		lib.Assert(s.seqs == nil)
		return *s.str
	default:
		panic("Empty sequence")
	}
}

// parse parses a single full sequence or option from the beginning of toks.
// Returns the parsed sequence and the number of tokens that were consumed.
func parse(toks []string) (*seq, int) {
	s := &seq{}
	var cons int // consumed tokens
Loop:
	for cons < len(toks) {
		start := toks[cons]
		switch start {
		case "(":
			grp := &seq{}
			cons++ // consume opening paren
			for {
				opt, n := parse(toks[cons:]) // consume group
				grp.opts = append(grp.opts, opt)
				cons += n
				if cons >= len(toks) {
					panic("Unclosed group")
				}
				next := toks[cons]
				switch next {
				case "|":
					cons++ // consume separator
				case ")":
					cons++ // consume closing paren
					s.seqs = append(s.seqs, grp)
					continue Loop
				default:
					lib.Panicf("Invalid token %q after option", next)
				}
			}
		case "|", ")":
			break Loop // let these be consumed by whoever called us
		default:
			s.seqs = append(s.seqs, &seq{str: &start})
			cons++ // consume literal string
		}
	}

	// If we didn't add anything, then this is an empty string.
	if len(s.seqs) == 0 && len(s.opts) == 0 && s.str == nil {
		var empty string
		s.str = &empty
	}

	return s, cons
}
