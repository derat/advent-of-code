package lib

// If returns a if cond is true and b otherwise.
// It does not short-circuit.
func If(cond bool, a, b int) int {
	if cond {
		return a
	}
	return b
}

// If returns a if cond is true and b otherwise.
// It does not short-circuit.
func IfByte(cond bool, a, b byte) byte {
	if cond {
		return a
	}
	return b
}
