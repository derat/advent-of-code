package lib

// Heap implements a binary heap.
// See https://www.cs.princeton.edu/~wayne/cs423/lectures/heaps-4up.pdf.
type Heap struct {
	items  []interface{}
	before HeapFunc
}

// HeapFunc returns true if a is ordered before b.
type HeapFunc func(a, b interface{}) bool

func NewHeap(before HeapFunc) *Heap {
	return &Heap{before: before}
}

func (h *Heap) Len() int {
	return len(h.items)
}

func (h *Heap) Insert(x interface{}) {
	// Add new item as rightmost leaf.
	h.items = append(h.items, x)

	// Bubble item up.
	for i := len(h.items) - 1; i > 0; i /= 2 {
		parent := (i+1)/2 - 1
		if !h.before(x, h.items[parent]) {
			break
		}
		h.items[i], h.items[parent] = h.items[parent], h.items[i]
	}
}

func (h *Heap) Pop() interface{} {
	if len(h.items) == 0 {
		panic("Heap is empty")
	}

	v := h.items[0]

	// Promote rightmost leaf to root.
	h.items[0] = h.items[len(h.items)-1]
	h.items = h.items[:len(h.items)-1]

	// Bubble root down.
	i := 0
	for 2*(i+1)-1 < len(h.items) {
		left := 2*(i+1) - 1
		right := 2 * (i + 1)
		if right >= len(h.items) || h.before(h.items[left], h.items[right]) {
			if h.before(h.items[left], h.items[i]) {
				h.items[left], h.items[i] = h.items[i], h.items[left]
				i = left
			} else {
				break
			}
		} else if right < len(h.items) {
			if h.before(h.items[right], h.items[i]) {
				h.items[right], h.items[i] = h.items[i], h.items[right]
				i = right
			} else {
				break
			}
		}
	}

	return v
}

type heapNode struct {
	data    interface{}
	parent  *heapNode // nil for root
	child   *heapNode // leftmost child
	sibling *heapNode // right sibling, nil if rightmost child
	degree  int       // number of children
}
