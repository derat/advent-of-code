package lib

// PackInts packs vals into a uint64, using the supplied number of bits for each element.
// Values must fit within the supplied bits. Use UnpackIntSigned to unpack signed ints.
func PackInts(vals []int, bits int) uint64 {
	AssertLessEq(len(vals)*bits, 64)
	var packed uint64
	for i, v := range vals {
		packed = PackInt(packed, v, bits, i)
	}
	return packed
}

// UnpackInts unpacks n unsigned values previously packed using PackInts.
func UnpackInts(packed uint64, bits, n int) []int {
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
