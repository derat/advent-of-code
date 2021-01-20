package lib

// Heap implements a binary heap.
// See https://runestone.academy/runestone/books/published/pythonds/Trees/BinaryHeapImplementation.html.
// (I originally used the description from https://www.cs.princeton.edu/~wayne/cs423/lectures/heaps-4up.pdf
// but I ended up with buggy code; probably I messed up while writing it.)
type Heap struct {
	items  []interface{}
	size   int
	before HeapFunc
}

// HeapFunc returns true if a is ordered before b.
type HeapFunc func(a, b interface{}) bool

func NewHeap(before HeapFunc) *Heap {
	return &Heap{
		items:  []interface{}{nil}, // first item is unused but makes code simpler
		before: before,
	}
}

func (h *Heap) Len() int {
	return h.size
}

func (h *Heap) Insert(x interface{}) {
	// Add new item as rightmost leaf.
	h.items = append(h.items, x)
	h.size++

	// Bubble item up.
	for i := h.size; i/2 > 0; i /= 2 {
		parent := i / 2
		if h.before(x, h.items[parent]) {
			h.items[i], h.items[parent] = h.items[parent], h.items[i]
		}
	}
}

func (h *Heap) Pop() interface{} {
	if h.size == 0 {
		panic("Heap is empty")
	}

	v := h.items[1]

	// Promote rightmost leaf to root.
	h.items[1] = h.items[h.size]
	h.items = h.items[:len(h.items)-1]
	h.size--

	// Bubble root down.
	i := 1
	for i*2 <= h.size {
		var mc int
		if i*2+1 > h.size || h.before(h.items[i*2], h.items[i*2+1]) {
			mc = i * 2
		} else {
			mc = i*2 + 1
		}
		if h.before(h.items[mc], h.items[i]) {
			h.items[mc], h.items[i] = h.items[i], h.items[mc]
		}
		i = mc
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
