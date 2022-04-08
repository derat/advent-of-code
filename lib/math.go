// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import "golang.org/x/exp/constraints"

// Min returns the minimum of the supplied values.
func Min[T constraints.Ordered](vals ...T) T {
	Assertf(len(vals) > 0, "No values given")
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

// Max returns the maximum of the supplied values.
func Max[T constraints.Ordered](vals ...T) T {
	Assertf(len(vals) > 0, "No values given")
	max := vals[0]
	for _, v := range vals[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// AtLeast returns the number of values greater than or equal to n.
func AtLeast[T constraints.Ordered](n T, vals ...T) int {
	var cnt int
	for _, v := range vals {
		if v >= n {
			cnt++
		}
	}
	return cnt
}

// Clamp clamps val within [min, max].
func Clamp[T constraints.Ordered](val, min, max T) T {
	return Min(Max(val, min), max)
}

// Sum returns the sum of the supplied values.
func Sum[T constraints.Integer | constraints.Float](vals ...T) T {
	var sum T
	for _, v := range vals {
		sum += v
	}
	return sum
}

// Product returns the product of the supplied values.
func Product[T constraints.Integer | constraints.Float](vals ...T) T {
	Assertf(len(vals) > 0, "No values given")
	prod := vals[0]
	for _, v := range vals[1:] {
		prod *= v
	}
	return prod
}

// Abs returns the absolute value of v.
func Abs[T constraints.Signed | constraints.Float](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// Pow returns x to the power of n.
func Pow[T constraints.Integer](x T, n int) T {
	return powInt(1, x, n)
}

// https://en.wikipedia.org/wiki/Exponentiation_by_squaring
func powInt[T constraints.Integer](y, x T, n int) T {
	switch {
	case n < 0:
		panic("Negative exponent")
	case n == 0:
		return y
	case n == 1:
		return x * y
	case n%2 == 0:
		return powInt(y, x*x, n/2)
	default:
		return powInt(x*y, x*x, (n-1)/2)
	}
}

// GCD returns the greatest common denominator of a and b using the Euclidean algorithm.
// See https://www.khanacademy.org/computing/computer-science/cryptography/modarithmetic/a/the-euclidean-algorithm.
func GCD[T constraints.Integer](a, b T) T {
	// From Khan Academy:
	//  If A = 0 then GCD(A,B)=B, since the GCD(0,B)=B, and we can stop.
	//  If B = 0 then GCD(A,B)=A, since the GCD(A,0)=A, and we can stop.
	//  Write A in quotient remainder form (A = B*Q + R)
	//  Find GCD(B,R) using the Euclidean Algorithm since GCD(A,B) = GCD(B,R)
	for b != 0 {
		if a == 0 {
			return b
		}
		b0 := b
		b = a % b
		a = b0
	}
	return a
}

// LCM reaturns the least common multiple of the supplied integers.
func LCM[T constraints.Integer](vals ...T) T {
	Assert(len(vals) > 0)
	if len(vals) == 1 {
		return vals[0]
	}
	res := vals[0] * vals[1] / GCD(vals[0], vals[1])
	for i := 2; i < len(vals); i++ {
		res = LCM(res, vals[i])
	}
	return res
}
