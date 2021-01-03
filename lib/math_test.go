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
