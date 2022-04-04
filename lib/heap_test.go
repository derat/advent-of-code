package lib

import "testing"

func TestHeap(t *testing.T) {
	h := NewHeap(func(a, b int) bool { return a < b })

	for _, tc := range []struct {
		op  string
		val int
	}{
		{"len", 0},
		{"insert", 3},
		{"insert", 1},
		{"insert", 10},
		{"insert", 0},
		{"insert", 1},
		{"insert", 2},
		{"len", 6},
		{"pop", 0},
		{"pop", 1},
		{"pop", 1},
		{"pop", 2},
		{"pop", 3},
		{"insert", 20},
		{"insert", 5},
		{"insert", 10},
		{"pop", 5},
		{"pop", 10},
		{"pop", 10},
		{"pop", 20},
		{"len", 0},
	} {
		switch tc.op {
		case "insert":
			h.Insert(tc.val)
		case "len":
			if got := h.Len(); got != tc.val {
				t.Fatalf("Len() = %v; want %v", got, tc.val)
			}
		case "pop":
			if got := h.Pop(); got != tc.val {
				t.Fatalf("Pop() = %v; want %v", got, tc.val)
			}
		}
	}
}
