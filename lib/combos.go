package lib

import "fmt"

// FindCombos returns all combinations of the supplied items that sum exactly to target.
// Initial is a bitfield specifying available items, e.g. if 0x1 is set then items[0]
// can be used. Pass 1<<len(items)-1 to use all items.
func FindCombos(items []int, initial uint64, target int) []uint64 {
	type set map[uint64]struct{}            // combos of items (specified by index)
	AssertLessEq(len(items), 64)            // packing into uint64
	sumRemain := make([]set, target+1)      // combos of remaining items for each sum
	sumRemain[0] = set{initial: struct{}{}} // base case

	for sum := 1; sum < len(sumRemain); sum++ {
		combos := make(set)
		for i, val := range items {
			// Carry forward any earlier combos where this item was available.
			if prev := sum - val; prev >= 0 {
				for rem := range sumRemain[prev] {
					if HasBit(rem, i) {
						combos[SetBit(rem, i, false)] = struct{}{} // use the item
					}
				}
			}
		}
		sumRemain[sum] = combos
		fmt.Println(sum, len(combos))
	}

	res := make([]uint64, 0, len(sumRemain[target]))
	for s := range sumRemain[target] {
		res = append(res, s)
	}
	return res
}
