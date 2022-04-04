package lib

// If returns a if cond is true and b otherwise.
// It does not short-circuit.
func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
