package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	sum := &snafu{}
	for _, ln := range lib.InputLines("2022/25") {
		sum.add(parse(ln))
	}
	fmt.Println(sum)
}

type snafu struct {
	c []int // coefficients, e.g. c[0]*1 + c[1]*5 + c[2]*25 + ...
}

func parse(str string) *snafu {
	s := snafu{make([]int, len(str))}
	for i := 0; i < len(str); i++ {
		ch := str[len(str)-i-1]
		switch ch {
		case '2':
			s.c[i] = 2
		case '1':
			s.c[i] = 1
		case '0':
			s.c[i] = 0
		case '-':
			s.c[i] = -1
		case '=':
			s.c[i] = -2
		default:
			lib.Panicf("Invalid char %q in %q", ch, s)
		}
	}
	return &s
}

func (s *snafu) add(o *snafu) {
	addc := func(i, n int) {
		if i >= len(s.c) {
			alloc := make([]int, i+1)
			copy(alloc, s.c)
			s.c = alloc
		}
		s.c[i] += n
	}
	for i, n := range o.c {
		addc(i, n)
	}

	// Reduce coefficients that are outside of the [-2, 2] range by
	// incrementing the coefficients for the next highest powers.
	for i := 0; i < len(s.c); i++ {
		for s.c[i] > 2 {
			addc(i+1, 1)
			s.c[i] -= 5
		}
		for s.c[i] < -2 {
			addc(i+1, -1)
			s.c[i] += 5
		}
	}
}

func (s *snafu) String() string {
	var str string
	for i := 0; i < len(s.c); i++ {
		n := s.c[len(s.c)-i-1]
		switch n {
		case 2:
			str += "2"
		case 1:
			str += "1"
		case 0:
			str += "0"
		case -1:
			str += "-"
		case -2:
			str += "="
		default:
			lib.Panicf("Non-reduced coefficient %v at position %v", n, i)
		}
	}
	return str
}
