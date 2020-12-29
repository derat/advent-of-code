package lib

import (
	"reflect"
	"testing"
)

func TestPackInt_RoundTrip(t *testing.T) {
	const bits = 11
	orig := []int{(1 << bits) - 1, 54, 0, 1 << (bits - 1), 85}
	packed := PackInts(orig, bits)
	got := UnpackInts(packed, bits, len(orig))
	if !reflect.DeepEqual(got, orig) {
		t.Errorf("UnpackInts(%#x, %d, %d) = %v; want %v", packed, bits, len(orig), got, orig)
	}
}

func TestUnpackIntSigned(t *testing.T) {
	const bits = 4
	orig := []int{-7, 5, -1, 0, 7, -2}
	packed := PackInts(orig, bits)
	for i, want := range orig {
		if got := UnpackIntSigned(packed, bits, i); got != want {
			t.Errorf("UnpackIntSigned(%#b, %d, %d) = %v; want %v", packed, bits, i, got, want)
		}
	}
}