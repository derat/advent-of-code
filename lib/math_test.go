// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import "testing"

func TestPow(t *testing.T) {
	for x := 1; x <= 10; x++ {
		want := 1
		for y := 0; y <= 10; y++ {
			if y > 0 {
				want *= x
			}
			if got := Pow(x, y); got != want {
				t.Errorf("Pow(%v, %v) = %v; want %v", x, y, got, want)
			}
		}
	}
}

func TestGCD(t *testing.T) {
	for _, tc := range []struct{ a, b, want int }{
		{1, 1, 1},
		{4, 4, 4},
		{2, 4, 2},
		{4, 6, 2},
		{4, 16, 4},
		{3, 5, 1},
	} {
		if got := GCD(tc.a, tc.b); got != tc.want {
			t.Errorf("GCD(%v, %v) = %v; want %v", tc.a, tc.b, got, tc.want)
		}
	}
}

func TestLCM(t *testing.T) {
	for _, tc := range []struct {
		in   []int
		want int
	}{
		{[]int{1, 1}, 1},
		{[]int{1, 2}, 2},
		{[]int{2, 1}, 2},
		{[]int{2, 4}, 4},
		{[]int{2, 6}, 6},
		{[]int{3, 5}, 15},
		{[]int{6, 9}, 18},
		{[]int{3, 2, 6}, 6},
		{[]int{3, 2, 5, 6, 10}, 30},
	} {
		if got := LCM(tc.in...); got != tc.want {
			t.Errorf("LCM(%v) = %v; want %v", tc.in, got, tc.want)
		}
	}
}
