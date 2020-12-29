package main

import (
	"fmt"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	statSet := make(map[string]struct{})
	ingreds := make(map[string]map[string]int)
	for _, ln := range lib.InputLines("2015/15") {
		var name, rest string
		lib.Extract(ln, `^(\w+): (.+)$`, &name, &rest)
		stats := make(map[string]int)
		for _, s := range strings.Split(rest, ", ") {
			var n string
			var v int
			lib.Extract(s, `^(\w+) (-?\d+)$`, &n, &v)
			stats[n] = v
			statSet[n] = struct{}{}
		}
		ingreds[name] = stats
	}

	ingredNames := lib.MapStringKeys(ingreds)
	ingredBits := 64 / len(ingreds)

	statNames := lib.MapStringKeys(statSet)
	statBits := 64/len(statSet) - 1 // signed

	computeScore := func(stats uint64) int64 {
		prod := int64(1)
		for i, n := range statNames {
			if n == "calories" {
				continue // calories aren't included in score
			}
			v := lib.UnpackIntSigned(stats, statBits, i)
			if v <= 0 {
				return 0
			}
			prod *= int64(v)
		}
		return prod
	}

	// Optimization idea: memoize combinations as we go.
	var findRecipe func(uint64, uint64, int, int) (uint64, uint64)
	findRecipe = func(amounts, stats uint64, idx, rem int) (topAmounts, topStats uint64) {
		name := ingredNames[idx]
		ingrStats, ok := ingreds[name]
		lib.Assert(ok)

		var topScore int64

		// We need to fill the remainder if we're the final ingredient.
		final := idx == len(ingredNames)-1
		minAmt := 0
		if final {
			minAmt = rem
		}
		for amt := minAmt; amt <= rem; amt++ {
			newAmounts := lib.PackInt(amounts, amt, ingredBits, idx)
			newStats := stats
			for i, sname := range statNames {
				v := lib.UnpackIntSigned(newStats, statBits, i)
				v += ingrStats[sname] * amt
				newStats = lib.PackInt(newStats, v, statBits, i)
			}
			// If we need to add more ingredients, recurse.
			if amt < rem && !final {
				newAmounts, newStats = findRecipe(newAmounts, newStats, idx+1, rem-amt)
			}
			if score := computeScore(newStats); score > topScore {
				topScore = score
				topAmounts = newAmounts
				topStats = newStats
			}
		}

		return topAmounts, topStats
	}

	const total = 100 // given in problem
	_, topStats := findRecipe(0, 0, 0, total)
	fmt.Println(computeScore(topStats))
}
