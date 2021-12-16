package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	lns := lib.InputLines("2021/16")
	lib.AssertEq(len(lns), 1)
	b, err := hex.DecodeString(lns[0])
	if err != nil {
		lib.Panicf("Failed decoding input: %v", err)
	}

	r := lib.NewBitReader(b, 0)
	root := readPacket(r)
	fmt.Println(root.vsum())
	fmt.Println(root.eval())
}

type packet interface {
	String() string
	vsum() int
	eval() int
}

// literal is a packet that holds an integer value.
type literal struct {
	ver, val int
}

func (p *literal) String() string { return fmt.Sprintf("%vv%v", p.val, p.ver) }
func (p *literal) vsum() int      { return p.ver }
func (p *literal) eval() int      { return p.val }

// operator a packet that holds an opcode and a variable number of subpackets.
type operator struct {
	ver, op int
	sub     []packet
}

func (p *operator) String() string {
	ss := make([]string, len(p.sub))
	for i, s := range p.sub {
		ss[i] = s.String()
	}
	return fmt.Sprintf("op%vv%v(%v)", p.op, p.ver, strings.Join(ss, " "))
}

func (p *operator) vsum() int {
	sum := p.ver
	for _, s := range p.sub {
		sum += s.vsum()
	}
	return sum
}

func (p *operator) eval() int {
	lib.AssertLess(0, len(p.sub))
	vals := make([]int, len(p.sub))
	for i, s := range p.sub {
		vals[i] = s.eval()
	}

	switch p.op {
	case 0:
		return lib.Sum(vals...)
	case 1:
		return lib.Product(vals...)
	case 2:
		return lib.Min(vals...)
	case 3:
		return lib.Max(vals...)
	case 5:
		lib.AssertEq(len(vals), 2)
		return lib.If(vals[0] > vals[1], 1, 0)
	case 6:
		lib.AssertEq(len(vals), 2)
		return lib.If(vals[0] < vals[1], 1, 0)
	case 7:
		lib.AssertEq(len(vals), 2)
		return lib.If(vals[0] == vals[1], 1, 0)
	default:
		lib.Panicf("Invalid operation %v", p.op)
		return 0
	}
}

// readPacket reads the next packet (and any subpackets) from r.
func readPacket(r *lib.BitReader) packet {
	ver := r.ReadInt(3)
	tid := r.ReadInt(3)
	switch tid {
	case 4:
		// "Packets with type ID 4 represent a literal value. Literal value packets encode a single
		// binary number. To do this, the binary number is padded with leading zeroes until its
		// length is a multiple of four bits, and then it is broken into groups of four bits. Each
		// group is prefixed by a 1 bit except the last group, which is prefixed by a 0 bit. These
		// groups of five bits immediately follow the packet header."
		var val int
		for true {
			n := r.ReadInt(5)
			val = (val << 4) | n&0xf
			if n&0x10 == 0 {
				break
			}
		}
		return &literal{ver, val}

	default:
		// "An operator packet contains one or more packets. To indicate which subsequent binary
		// data represents its sub-packets, an operator packet can use one of two modes indicated by
		// the bit immediately after the packet header; this is called the length type ID. ...
		// Finally, after the length type ID bit and the 15-bit or 11-bit field, the sub-packets
		// appear."
		var subs []packet
		if r.ReadInt(1) == 0 {
			// "If the length type ID is 0, then the next 15 bits are a number that represents the
			// total length in bits of the sub-packets contained by this packet."
			slen := r.ReadInt(15)
			next := r.Offset() + slen
			for r.Offset() < next {
				subs = append(subs, readPacket(r))
			}
			lib.AssertEq(r.Offset(), next)
		} else {
			// "If the length type ID is 1, then the next 11 bits are a number that represents the
			// number of sub-packets immediately contained by this packet."
			nsubs := r.ReadInt(11)
			for i := 0; i < nsubs; i++ {
				subs = append(subs, readPacket(r))
			}
		}
		return &operator{ver, tid, subs}
	}
}
