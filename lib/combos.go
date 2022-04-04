package lib

// FindCombos returns all combinations of the supplied items that sum exactly to target.
// Initial is a bitfield specifying available items, e.g. if 0x1 is set then items[0]
// can be used. Pass 1<<len(items)-1 to use all items.
func FindCombos(items []int, initial uint64, target int) []uint64 {
	type set map[uint64]struct{}            // combos of items (specified by index)
	AssertLessEq(len(items), 64)            // packing into uint64
	sumRemain := make([]set, target+1)      // combos of remaining items for each sum
	sumRemain[0] = set{initial: struct{}{}} // base case

	for sum := 1; sum < len(sumRemain); sum++ {
		rems := make(set)
		for i, val := range items {
			// Carry forward any earlier combos where this item was available.
			if prev := sum - val; prev >= 0 {
				for rem := range sumRemain[prev] {
					if HasBit(rem, i) {
						rems[SetBit(rem, i, false)] = struct{}{} // use the item
					}
				}
			}
		}
		sumRemain[sum] = rems
	}

	res := make([]uint64, 0, len(sumRemain[target]))
	for rem := range sumRemain[target] {
		res = append(res, initial&^rem)
	}
	return res
}

// Perms sends all permutations of the supplied slice to ch and closes it.
// This is the non-recursive version of https://en.wikipedia.org/wiki/Heap%27s_algorithm.
// s is modified in-place.
func Perms[T any](s []T, ch chan []T) {
	send := func() {
		c := make([]T, len(s))
		copy(c, s)
		ch <- c
	}

	state := make([]int, len(s))
	send()

	var i int
	for i < len(s) {
		if state[i] < i {
			if i%2 == 0 {
				s[0], s[i] = s[i], s[0]
			} else {
				s[state[i]], s[i] = s[i], s[state[i]]
			}
			send()
			state[i]++
			i = 0
		} else {
			state[i] = 0
			i++
		}
	}
	close(ch)
}
