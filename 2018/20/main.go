package main

import (
	"fmt"
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
	//lib.AssertEq(input, root.String())
	lib.AssertEq(n, len(tokens))

	doors := make(map[uint64]struct{}) // keys are packed x, y, dir
	visit(0, 0, root.seqs, doors)

	// Perform a BFS to find the shortest path to each room.
	rooms := lib.BFS(lib.PackInts(0, 0), func(s uint64) []uint64 {
		x, y := lib.UnpackIntSigned2(s)
		var next []uint64
		for _, d := range []dir{north, south, east, west} {
			if lib.MapHasKey(doors, lib.PackInts(x, y, int(d))) {
				next = append(next, lib.PackInts(x+d.dx(), y+d.dy()))
			}
		}
		return next
	}, nil, -1)

	// Part 1: Largest number of doors required to pass through to reach a room,
	// i.e. the longest of all rooms' shortest paths.
	fmt.Println(lib.Max(lib.MapIntVals(rooms)...))

	// Part 2: Number of rooms with a shortest path of at least 1000 doors.
	fmt.Println(lib.AtLeast(1000, lib.MapIntVals(rooms)...))
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
		first := toks[cons]
		switch first {
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
					simplify(grp)
					s.seqs = append(s.seqs, grp)
					continue Loop
				default:
					lib.Panicf("Invalid token %q after option", next)
				}
			}
		case "|", ")":
			break Loop // let these be consumed by whoever called us
		default:
			s.seqs = append(s.seqs, &seq{str: &first})
			cons++ // consume literal string
		}
	}

	// If we didn't add anything, then this is an empty string.
	if len(s.seqs) == 0 && len(s.opts) == 0 && s.str == nil {
		var empty string
		s.str = &empty
	}

	// Simplify single-item sequences.
	if len(s.seqs) == 1 {
		*s = *s.seqs[0]
	}

	return s, cons
}

// simplify reduces grp to a single string if it was formerly
// a pointless two-option group like "(NWSE|)".
func simplify(grp *seq) {
	if len(grp.opts) != 2 {
		return
	}
	a, b := grp.opts[0], grp.opts[1]
	if a.str == nil || b.str == nil {
		return
	}
	as, bs := *a.str, *b.str
	if loop(as) && bs == "" {
		grp.str = &as
		grp.opts = nil
	} else if loop(bs) && as == "" {
		grp.str = &bs
		grp.opts = nil
	}
}

// loop returns true if str loops around to its starting point,
// e.g. "NS" or "NSWWENSE".
func loop(str string) bool {
	n := strings.Count(str, "N")
	s := strings.Count(str, "S")
	e := strings.Count(str, "E")
	w := strings.Count(str, "W")
	return n == s && e == w
}

// visit recursively follows all paths reachable using seqs starting at x, y.
// It records doors that are used in the supplied map.
func visit(x, y int, seqs []*seq, doors map[uint64]struct{}) {
	if len(seqs) == 0 {
		return
	}

	first, rest := seqs[0], seqs[1:]

	switch {
	case len(first.opts) > 0:
		for _, o := range first.opts {
			visit(x, y, append([]*seq{o}, rest...), doors)
		}
	case len(first.seqs) > 0:
		var next []*seq
		next = append(next, first.seqs...)
		next = append(next, rest...)
		visit(x, y, next, doors)
	case first.str != nil:
		for _, ch := range *first.str {
			var d dir
			switch ch {
			case 'N':
				d = north
			case 'S':
				d = south
			case 'W':
				d = west
			case 'E':
				d = east
			}
			// Add the door in both directions.
			doors[lib.PackInts(x, y, int(d))] = struct{}{}
			x, y = x+d.dx(), y+d.dy()
			doors[lib.PackInts(x, y, int(d.opp()))] = struct{}{}
		}
		visit(x, y, rest, doors)
	default:
		panic("Nothing to do")
	}
}

type dir int

const (
	north dir = iota
	east
	south
	west
)

func (d dir) dx() int {
	switch d {
	case west:
		return -1
	case east:
		return 1
	default:
		return 0
	}
}

func (d dir) dy() int {
	switch d {
	case north:
		return -1
	case south:
		return 1
	default:
		return 0
	}
}

func (d dir) opp() dir {
	return (d + 2) % 4
}
