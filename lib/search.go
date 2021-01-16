package lib

// AStar uses the A* algorithm to find the minimum number of steps from the initial
// state(s) to a state where the done function returns true. The next function should
// return all states reachable in a single step from the state passed to it, and the
// estimate function should return a lower bound on the number of steps to reach
// a target state. See https://www.redblobgames.com/pathfinding/a-star/introduction.html.
func AStar(initial []uint64,
	done func(uint64) bool,
	next func(uint64) []uint64,
	estimate func(uint64) int) int {
	// TODO: Add some way to track the path if needed.
	frontier := NewHeap(func(a, b interface{}) bool { return a.(asNode).pri < b.(asNode).pri })
	costs := make(map[uint64]int)
	for _, init := range initial {
		frontier.Insert(asNode{init, 0})
		costs[init] = 0
	}

	for frontier.Len() != 0 {
		cur := frontier.Pop().(asNode).state
		cost := costs[cur]

		// Check if we're done.
		if done(cur) {
			return cost
		}

		for _, next := range next(cur) {
			newCost := cost + 1
			if old, ok := costs[next]; ok && old <= newCost {
				continue // already visited with equal or lower cost
			}
			pri := newCost + estimate(next)
			frontier.Insert(asNode{next, pri})
			costs[next] = newCost
		}
	}

	return -1
}

type asNode struct {
	state uint64
	pri   int // lower is better
}

// BFS returns a map of the minimum number of steps to go from start to reachable states.
// If dests is non-empty, returns when all have been reached. Otherwise, goes to all reachable states.
// If max is non-negative, returns as soon as the cost exceeds it.
func BFS(start uint64, next func(s uint64) []uint64, dests []uint64, max int) map[uint64]int {
	queue := []uint64{start}
	seen := map[uint64]struct{}{start: struct{}{}}
	costs := map[uint64]int{start: 0}

	remain := make(map[uint64]struct{}, len(dests))
	for _, d := range dests {
		remain[d] = struct{}{}
	}

Loop:
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		cost := costs[cur] + 1
		if max >= 0 && cost > max {
			break
		}
		for _, st := range next(cur) {
			if _, ok := seen[st]; !ok {
				queue = append(queue, st)
				seen[st] = struct{}{}
				costs[st] = cost

				// Early exit if we've reached all specified destinations.
				delete(remain, st)
				if len(dests) > 0 && len(remain) == 0 {
					break Loop
				}
			}
		}
	}

	return costs
}
