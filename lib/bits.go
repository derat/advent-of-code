package lib

import (
	"math"
)

// PackInts packs vals into a uint64, dividing the total bits evenly across the values.
// Values must fit within the supplied bits. Use UnpackIntSigned to unpack signed ints.
func PackInts(vals ...int) uint64 {
	bits := 64 / len(vals)
	var packed uint64
	for i, v := range vals {
		packed = PackInt(packed, v, bits, i*bits)
	}
	return packed
}

// UnpackInts unpacks n unsigned values previously packed using PackInts.
func UnpackInts(packed uint64, n int) []int {
	bits := 64 / n
	vals := make([]int, n)
	for i := range vals {
		vals[i] = UnpackInt(packed, bits, i*bits)
	}
	return vals
}

// PackInt sets a size-bit region at the supplied offset (from the LSB)
// in packed to val.
func PackInt(packed uint64, val, size, offset int) uint64 {
	mask := uint64(1<<size - 1)
	packed &= ^(mask << offset)
	return packed | (uint64(val&int(mask)) << offset)
}

// UnpackInt unpacks and returns an unsigned value of size bits at the supplied
// offset (from the LSB) from packed.
func UnpackInt(packed uint64, size, offset int) int {
	return int((packed >> offset) & (1<<size - 1))
}

// UnpackIntSigned is like UnpackInt but with support for negative numbers.
func UnpackIntSigned(packed uint64, size, offset int) int {
	val := UnpackInt(packed, size, offset)
	shift := 64 - size
	return int((int64(val) << shift) >> shift) // extend sign bit
}

// UnpackInt2 is a convenience function that unpacks two 32-bit values from p.
func UnpackInt2(p uint64) (a, b int) {
	return UnpackInt(p, 32, 0), UnpackInt(p, 32, 32)
}

// UnpackIntSigned2 is like UnpackInt2 but for signed ints.
func UnpackIntSigned2(p uint64) (a, b int) {
	return UnpackIntSigned(p, 32, 0), UnpackIntSigned(p, 32, 32)
}

// SetBit sets the i-th bit in field to v and returns the field.
func SetBit(field uint64, i int, v bool) uint64 {
	if v {
		return field | (1 << i)
	}
	return field &^ (1 << i)
}

// HasBit returns true if the i-th bit in field is set.
func HasBit(field uint64, i int) bool {
	return field&(1<<i) != 0
}

// Hi returns the top 4 bits of b.
func Hi(b byte) byte {
	return (b >> 4) & 0xf
}

// Hi returns the bottom 4 bits of b.
func Lo(b byte) byte {
	return b & 0xf
}

// HiIsLo returns true if Hi(b) equals Lo(b).
func HiIsLo(b byte) bool {
	return Hi(b) == Lo(b)
}

// ExtractBits zeros all bits in v except the ones between indexes high and low
// (inclusive) and right-shifts the resulting value by low.
//
// In other words:
//
//  val: abcdefgh
//  pos: 76543210
//
// If high is 7 and low is 0, returns abcdefgh.
// If high is 5 and low is 1, returns 000cdefg.
// If high is 3 and low is 3, returns 0000000e.
func ExtractBits(v uint64, high, low int) uint64 {
	AssertLessEq(low, high)
	nbits := high - low + 1
	mask := uint64(1<<nbits - 1)
	return (v >> low) & mask
}

// BitReader reads variable numbers of bits from a byte slice.
type BitReader struct {
	b   []byte
	off int // current bit offset into b
}

// NewBitReader returns a BitReader that reads bits from b starting at the
// specified bit offset (with 0 being the MSB of the first byte in the slice).
func NewBitReader(b []byte, off int) *BitReader {
	return &BitReader{b, off}
}

// Offset returns the current offset into the slice that was passed to NewBitReader.
func (br *BitReader) Offset() int { return br.off }

// Read reads and returns the specified number of bits, advancing the current offset.
func (br *BitReader) Read(nbits int) uint64 {
	AssertLessEq(1, nbits)
	AssertLessEq(nbits, 64)
	size := 8 * len(br.b)
	if br.off+nbits > size {
		Panicf("%v-bit read at %v overruns %v-bit buffer", nbits, br.off, size)
	}
	var v uint64
	for nbits > 0 {
		cur := br.b[br.off/8]       // byte to extract bits from
		high := 7 - (br.off % 8)    // index of high bit (7 is MSB, 0 is LSB)
		low := Max(high-nbits+1, 0) // index of low bit
		nb := high - low + 1        // number of bits to extract from cur
		v = (v << nb) | ExtractBits(uint64(cur), high, low)
		nbits -= nb
		br.off += nb
	}
	AssertEq(nbits, 0)
	return v
}

// ReadInt is a convenience method that calls Read and casts the returned value to an int.
func (br *BitReader) ReadInt(nbits int) int {
	// TODO: This ought to be math.MaxInt, but it wasn't added until Go 1.17:
	// https://github.com/golang/go/issues/28538
	AssertLessEq(nbits, math.MaxInt64)
	return int(br.Read(nbits))
}
