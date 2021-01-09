package lib

// PackInts packs vals into a uint64, dividing the total bits evenly across the values.
// Values must fit within the supplied bits. Use UnpackIntSigned to unpack signed ints.
func PackInts(vals ...int) uint64 {
	bits := 64 / len(vals)
	var packed uint64
	for i, v := range vals {
		packed = PackInt(packed, v, bits, i)
	}
	return packed
}

// UnpackInts unpacks n unsigned values previously packed using PackInts.
func UnpackInts(packed uint64, n int) []int {
	bits := 64 / n
	vals := make([]int, n)
	for i := range vals {
		vals[i] = UnpackInt(packed, bits, i)
	}
	return vals
}

// PackInt sets packed's i-th position to val.
func PackInt(packed uint64, val, bits, i int) uint64 {
	mask := uint64(1<<bits - 1)
	packed &= ^(mask << (i * bits))
	return packed | (uint64(val&int(mask)) << (i * bits))
}

// UnpackInt unpacks and returns the unsigned value at position i.
func UnpackInt(packed uint64, bits, i int) int {
	return int((packed >> (i * bits)) & (1<<bits - 1))
}

// UnpackIntSigned is like UnpackInt but with support for negative numbers.
func UnpackIntSigned(packed uint64, bits, i int) int {
	val := UnpackInt(packed, bits, i)
	shift := 64 - bits
	return int((int64(val) << shift) >> shift) // extend sign bit
}

// UnpackInt2 is a convenience function that unpacks two 32-bit values from p.
func UnpackInt2(p uint64) (a, b int) {
	return UnpackInt(p, 32, 0), UnpackInt(p, 32, 1)
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
