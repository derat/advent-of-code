// Copyright 2022 Daniel Erat.
// All rights reserved.

package lib

import (
	"fmt"
	"reflect"
	"testing"
)

func TestPackInts_RoundTrip(t *testing.T) {
	const bits = 11
	orig := []int{(1 << bits) - 1, 54, 0, 1 << (bits - 1), 85}
	packed := PackInts(orig...)
	got := UnpackInts[int](packed, len(orig))
	if !reflect.DeepEqual(got, orig) {
		t.Errorf("UnpackInts(%#x, %d) = %v; want %v", packed, len(orig), got, orig)
	}
}

func TestUnpackIntSigned(t *testing.T) {
	orig := []int{-7, 5, -1, 0, 7, -2}
	bits := 64 / len(orig)
	packed := PackInts(orig...)
	for i, want := range orig {
		if got := UnpackIntSigned[int](packed, bits, i*bits); got != want {
			t.Errorf("UnpackIntSigned(%#b, %d, %d) = %v; want %v",
				packed, bits, i*bits, got, want)
		}
	}
}

func TestUnpackInt2(t *testing.T) {
	for _, vals := range [][2]int{
		{0, 0},
		{0, 1},
		{1, 0},
		{1, 1},
		{1 << 31, 1 << 31},
	} {
		a, b := UnpackInt2(PackInts(vals[0], vals[1]))
		if a != vals[0] || b != vals[1] {
			t.Errorf("UnpackInt2() = %v, %v; want %v, %v", a, b, vals[0], vals[1])
		}
	}
}

func TestSetBit(t *testing.T) {
	for _, tc := range []struct {
		init uint64
		idx  int
		val  bool
		want uint64
	}{
		{0, 0, true, 0x1},
		{0, 63, true, 1 << 63},
		{0, 0, false, 0x0},
		{1<<64 - 1, 0, true, 1<<64 - 1},
		{1<<64 - 1, 0, false, 1<<64 - 2},
		{1<<64 - 1, 63, true, 1<<64 - 1},
		{1<<64 - 1, 63, false, 0x7fffffffffffffff},
	} {
		if got := SetBit(tc.init, tc.idx, tc.val); got != tc.want {
			t.Errorf("SetBit(%#x, %v, %v) = %#x; want %#x", tc.init, tc.idx, tc.val, got, tc.want)
		}
	}
}

func TestExtractBits(t *testing.T) {
	var v uint64 = 0xdeadbeefcafeface
	for i := 0; i <= 63; i++ {
		want := uint64(If(HasBit(v, i), 1, 0))
		if got := ExtractBits(v, i, i); got != want {
			t.Errorf("ExtractBits(%#b, %v, %v) = %#b; want %#b", v, i, i, got, want)
		}
	}

	for _, tc := range []struct {
		v         uint64
		high, low int
		want      uint64
	}{
		// fedcba9876543210
		{0b1101101101010010, 0xf, 0x0, 0b1101101101010010},
		{0b1101101101010010, 0xf, 0xa, 0b110110},
		{0b1101101101010010, 0xb, 0x6, 0b101101},
		{0b1101101101010010, 0x5, 0x0, 0b010010},
	} {
		if got := ExtractBits(tc.v, tc.high, tc.low); got != tc.want {
			t.Errorf("ExtractBits(%#b, %v, %v) = %#b; want %#b", tc.v, tc.high, tc.low, got, tc.want)
		}
	}
}

func TestBitReader(t *testing.T) {
	for _, tc := range []struct {
		b     []byte
		off   int
		reads []int
		want  []uint64
	}{
		{[]byte{0b00000000}, 0, []int{8}, []uint64{0b00000000}},
		{[]byte{0b11011010}, 0, []int{8}, []uint64{0b11011010}},
		{[]byte{0b11111111}, 0, []int{8}, []uint64{0b11111111}},
		{[]byte{0b00000000}, 0, []int{2, 3, 2, 1}, []uint64{0b00, 0b000, 0b00, 0b0}},
		{[]byte{0b11011010}, 0, []int{2, 3, 2, 1}, []uint64{0b11, 0b011, 0b01, 0b0}},
		{[]byte{0b11111111}, 0, []int{2, 3, 2, 1}, []uint64{0b11, 0b111, 0b11, 0b1}},
		{[]byte{0b11011010}, 0, []int{1, 1, 1, 1, 1, 1, 1, 1}, []uint64{1, 1, 0, 1, 1, 0, 1, 0}},
		{[]byte{0b11011010}, 2, []int{2, 3, 1}, []uint64{0b01, 0b101, 0b0}},
		{[]byte{0b10110011, 0b01001101}, 0, []int{16}, []uint64{0b10110011_01001101}},
		{[]byte{0b10110011, 0b01001101}, 0, []int{2, 10, 4}, []uint64{0b10, 0b1100110100, 0b1101}},
		{[]byte{0b10110011, 0b01001101}, 9, []int{3, 4}, []uint64{0b100, 0b1101}},
		{[]byte{0b10110011, 0b01001101, 0b00101101, 0b00101111, 0b11110000, 0b11010001, 0b11011010, 0b00110100}, 0,
			[]int{64}, []uint64{0b10110011_01001101_00101101_00101111_11110000_11010001_11011010_00110100}},
	} {
		t.Run(fmt.Sprintf("%+v", tc), func(t *testing.T) {
			br := NewBitReader(tc.b, tc.off)
			for i, bits := range tc.reads {
				if got := br.Read(bits); got != tc.want[i] {
					t.Errorf("%v: Read(%d) = %b; want %b", i, bits, got, tc.want[i])
				}
			}
		})
	}
}
