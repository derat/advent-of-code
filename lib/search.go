package lib

// AStarVarCost uses the A* algorithm to find the minimum number of steps from the initial
// state(s) to a state where the done function returns true. The next function should
// return all states reachable in a single step from the state passed to it along with the
// corresponding cost, and the estimate function should return a lower bound on the remaining
// cost to reach a target state.
// See https://www.redblobgames.com/pathfinding/a-star/introduction.html.
func AStarVarCost(initial []uint64,
	done func(uint64) bool,
	next func(uint64) map[uint64]int,
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

		for ns, nc := range next(cur) {
			newCost := cost + nc
			if oldCost, ok := costs[ns]; !ok || newCost < oldCost {
				costs[ns] = newCost
				pri := newCost + estimate(ns)
				frontier.Insert(asNode{ns, pri})
			}
		}
	}

	return -1
}

// AStar is a simplified version of AStarVarCost that adapts the supplied next function
// to report a cost of 1 to move to each next state.
func AStar(initial []uint64,
	done func(uint64) bool,
	next func(uint64) []uint64,
	estimate func(uint64) int) int {
	return AStarVarCost(initial, done, func(s uint64) map[uint64]int {
		ns := next(s)
		nm := make(map[uint64]int, len(ns))
		for _, n := range ns {
			nm[n] = 1
		}
		return nm
	}, estimate)
}

type asNode struct {
	state uint64
	pri   int // lower is better
}

// BFS performs a breadth-first search to discover paths to states reachable from start.
// If opts is non-nil, it is used to configure the search.
// The returned steps map contains the minimum number of steps to each state.
// The returned from map contains the state preceding each destination state.
func BFS(start uint64, next func(s uint64) []uint64, opts *BFSOptions) (steps map[uint64]int, from map[uint64]uint64) {
	queue := []uint64{start} // next states to check
	steps = map[uint64]int{start: 0}
	from = map[uint64]uint64{start: start}

	var remain map[uint64]struct{}
	if opts != nil && len(opts.AllDests) > 0 {
		remain = make(map[uint64]struct{})
		for _, d := range opts.AllDests {
			remain[d] = struct{}{}
		}
	}

Loop:
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		cost := steps[cur] + 1

		// Early exit if we've exceeded the maximum number of steps.
		if opts != nil && opts.MaxSteps > 0 && cost > opts.MaxSteps {
			break Loop
		}

		for _, st := range next(cur) {
			// Skip already-visited states.
			if _, ok := from[st]; ok {
				continue
			}

			queue = append(queue, st)
			from[st] = cur
			steps[st] = cost

			// Early exit if we've reached one of the "any" destinations.
			if opts != nil && MapHasKey(opts.AnyDests, st) {
				break Loop
			}

			// Early exit if we've reached all required destinations.
			if remain != nil {
				delete(remain, st)
				if len(remain) == 0 {
					break Loop
				}
			}
		}
	}

	return steps, from
}

// BFSOptions specifies optional configuration for BFS.
type BFSOptions struct {
	// AllDests contains states that must all be reached before exiting.
	AllDests []uint64
	// AnyDests contains states of which just one must be reached before exiting.
	AnyDests map[uint64]struct{}
	// MaxSteps contains the maximum number of steps before exiting.
	MaxSteps int
}
