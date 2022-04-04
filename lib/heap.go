package lib

// Heap implements a binary heap.
// See https://runestone.academy/runestone/books/published/pythonds/Trees/BinaryHeapImplementation.html.
// (I originally used the description from https://www.cs.princeton.edu/~wayne/cs423/lectures/heaps-4up.pdf
// but I ended up with buggy code; probably I messed up while writing it.)
type Heap[T any] struct {
	items  []T
	size   int
	before HeapFunc[T]
}

// HeapFunc returns true if a is ordered before b.
type HeapFunc[T any] func(a, b T) bool

func NewHeap[T any](before HeapFunc[T]) *Heap[T] {
	var zero T
	return &Heap[T]{
		items:  []T{zero}, // first item is unused but makes code simpler
		before: before,
	}
}

func (h *Heap[T]) Len() int {
	return h.size
}

func (h *Heap[T]) Insert(x T) {
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

func (h *Heap[T]) Pop() T {
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

type heapNode[T any] struct {
	data    T
	parent  *heapNode[T] // nil for root
	child   *heapNode[T] // leftmost child
	sibling *heapNode[T] // right sibling, nil if rightmost child
	degree  int          // number of children
}
