package main

import (
	"fmt"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	monkeys := make(map[string]*monkey)
	children := make(map[*monkey][2]string)
	const namePat = `[a-z]{4}`
	for _, ln := range lib.InputLines("2022/21") {
		var m monkey
		var a, b string
		switch {
		case lib.TryExtract(ln, `^(`+namePat+`): (`+namePat+`) ([-+*/]) (`+namePat+`)`,
			&m.name, &a, &m.op, &b):
			children[&m] = [2]string{a, b}
		case lib.TryExtract(ln, `^(`+namePat+`): (\d+)$`, &m.name, &m.val):
		default:
			lib.Panicf("Bad line %q", ln)
		}
		monkeys[m.name] = &m
	}
	for m, ch := range children {
		m.a = monkeys[ch[0]]
		m.b = monkeys[ch[1]]
	}

	// Part 1: Figure out what "root" evalutes to.
	root := monkeys["root"]
	fmt.Println(root.yell())

	// Part 2: "humn" is now a variable instead of returning a fixed value.
	// What does it need to be to make root's operands be equal?
	// With both the example and my input, I get an easy-to-solve expression like a*x + b = c!
	a := root.a.yell2()
	b := root.b.yell2()
	if len(b) > len(a) {
		a, b = b, a
	}
	if len(a) != 2 || len(b) != 1 {
		lib.Panicf("Got unexpected polynomials %q and %q", a, b)
	}
	fmt.Println(b.sub(poly{a[0]}).div(poly{a[1]}))

}

type monkey struct {
	name string
	val  int64   // 0 if not yet computed
	val2 poly    // nil if not yet computed
	op   byte    // '+', '-', '*', '/', or 0 for literal
	a, b *monkey // operands
}

func (m *monkey) yell() int64 {
	if m.val != 0 {
		return m.val // literal or already cached
	}
	a, b := m.a.yell(), m.b.yell()
	switch m.op {
	case '+':
		m.val = a + b
	case '-':
		m.val = a - b
	case '*':
		m.val = a * b
	case '/':
		m.val = a / b
	default:
		lib.Panicf("Invalid operator %q", m.op)
	}
	return m.val
}

func (m *monkey) yell2() poly {
	if m.val2 != nil {
		return m.val2
	}
	switch {
	case m.name == "humn":
		m.val2 = poly([]fraction{zero, {1, 1}}) // 0 + 1*x
	case m.op == 0:
		m.val2 = poly{fraction{m.val, 1}}
	case m.op == '+':
		m.val2 = m.a.yell2().add(m.b.yell2())
	case m.op == '-':
		m.val2 = m.a.yell2().sub(m.b.yell2())
	case m.op == '*':
		m.val2 = m.a.yell2().mul(m.b.yell2())
	case m.op == '/':
		m.val2 = m.a.yell2().div(m.b.yell2())
	default:
		lib.Panicf("Invalid operator %q", m.op)
	}
	return m.val2
}

type fraction struct{ num, den int64 }

var zero = fraction{0, 1}

func (f fraction) simplify() fraction {
	lib.Assert(f.den != 0)
	if f.num == 0 {
		return fraction{0, 1}
	}
	gcd := lib.GCD(f.num, f.den)
	return fraction{f.num / gcd, f.den / gcd}
}

func (f fraction) invert() fraction { return fraction{f.den, f.num}.simplify() }
func (f fraction) add(o fraction) fraction {
	if f.den == o.den {
		return fraction{f.num + o.num, f.den}.simplify()
	}
	return fraction{f.num*o.den + o.num*f.den, f.den * o.den}.simplify()
}
func (f fraction) mul(o fraction) fraction { return fraction{f.num * o.num, f.den * o.den}.simplify() }
func (f fraction) sub(o fraction) fraction { return f.add(fraction{-o.num, o.den}) }
func (f fraction) div(o fraction) fraction { return f.mul(fraction{o.den, o.num}) }

// This was apparently overkill: I guess I could've gotten away with just scalar
// and constant fields. Should've analyzed my input before writing code...
type poly []fraction // s[0]*1 + s[1]*x + s[2]*(x**2) + ...

func (p poly) simplify() poly {
	// Drop trailing (higher-order) zero terms.
	for i := len(p) - 1; i >= 0; i-- {
		if p[i].num != 0 {
			break
		}
		p = p[:i]
	}
	return p
}

func (p poly) add(o poly) poly {
	if len(o) > len(p) {
		o, p = p, o
	}
	res := append(poly(nil), p...)
	for i, f := range o {
		res[i] = res[i].add(f)
	}
	return res.simplify()
}

func (p poly) sub(o poly) poly {
	neg := make(poly, len(o))
	for i, f := range o {
		neg[i] = fraction{-f.num, f.den}
	}
	return p.add(neg)
}

func (p poly) mul(o poly) poly {
	if len(o) > 1 {
		if len(p) > 1 {
			// Multiplication by anything beyond a constant apparently isn't needed.
			lib.Panicf("Can't multiply %s by %s", p, o)
		}
		p, o = o, p
	}
	res := make(poly, len(p))
	for i, f := range p {
		res[i] = f.mul(o[0])
	}
	return res.simplify()
}

func (p poly) div(o poly) poly {
	if len(o) > 1 {
		// Division by anything beyond a constant apparently isn't needed.
		lib.Panicf("Can't divide %s by %s", p, o)
	}
	return p.mul(poly{o[0].invert()})
}

// Wrote this for debugging, but it didn't come in particularly useful.
func (p poly) String() string {
	if len(p) == 0 {
		return "0"
	}

	var s string
	for i := len(p) - 1; i >= 0; i-- {
		f := p[i]
		neg := (f.num < 0 && f.den > 0) || (f.num > 0 && f.den < 0)
		first := i == len(p)-1

		if f.num == 0 {
			lib.Assert(!first)
			continue
		}

		switch {
		case first && neg:
			s += "-"
		case !first && neg:
			s += " - "
		case !first && !neg:
			s += " + "
		}

		switch {
		case f.den == 1:
			s += fmt.Sprintf("%v", lib.Abs(f.num))
		default:
			s += fmt.Sprintf("%v/%v", lib.Abs(f.num), lib.Abs(f.den))
		}

		switch {
		case i == 1:
			s += "*x"
		case i > 1:
			s += fmt.Sprintf("*x^%d", i)
		}
	}
	return s
}
