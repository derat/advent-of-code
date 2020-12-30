package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	var bossHP, bossDamage, bossArmor int
	lib.Extract(lib.Input("2015/21"), `^Hit Points: (\d+)\nDamage: (\d+)\nArmor: (\d+)\n$`,
		&bossHP, &bossDamage, &bossArmor)

	check := func(items ...item) (win bool, cost int) {
		var damage, armor int
		for _, it := range items {
			cost += it.cost
			damage += it.damage
			armor += it.armor
		}

		hp := 100 // given in description
		damage = lib.Max(damage-bossArmor, 1)

		bhp := bossHP
		bdam := lib.Max(bossDamage-armor, 1)

		turns := (hp + bdam - 1) / bdam
		bturns := (bhp + damage - 1) / damage

		return turns >= bturns, cost
	}

	// Part 1: Minimum cost to win
	minCost := math.MaxInt32
	for _, weapon := range weapons {
		for _, armor := range armors {
			for r1, ring1 := range rings {
				for _, ring2 := range rings[r1+1:] {
					if win, cost := check(weapon, armor, ring1, ring2); win && cost < minCost {
						minCost = cost
					}
				}
			}
			if win, cost := check(weapon, armor); win && cost < minCost { // no rings
				minCost = cost
			}
		}
	}
	fmt.Println(minCost)

	// Part 2: Maximum cost to lose
	maxCost := 0
	for _, weapon := range weapons {
		for _, armor := range armors {
			for r1, ring1 := range rings {
				for _, ring2 := range rings[r1+1:] {
					if win, cost := check(weapon, armor, ring1, ring2); !win && cost > maxCost {
						maxCost = cost
					}
				}
			}
			if win, cost := check(weapon, armor); !win && cost > maxCost { // no rings
				maxCost = cost
			}
		}
	}
	fmt.Println(maxCost)
}

type item struct {
	name                string
	cost, damage, armor int
}

// Given in description.
var weapons = []item{
	{"Dagger", 8, 4, 0},
	{"Shortsword", 10, 5, 0},
	{"Warhammer", 25, 6, 0},
	{"Longsword", 40, 7, 0},
	{"Greataxe", 74, 8, 0},
}
var armors = []item{
	{"[no armor]", 0, 0, 0}, // added
	{"Leather", 13, 0, 1},
	{"Chainmail", 31, 0, 2},
	{"Splintmail", 53, 0, 3},
	{"Bandedmail", 75, 0, 4},
	{"Platemail", 102, 0, 5},
}
var rings = []item{
	{"[no ring]", 0, 0, 0}, // added
	{"Damage +1", 25, 1, 0},
	{"Damage +2", 50, 2, 0},
	{"Damage +3", 100, 3, 0},
	{"Defense +1", 20, 0, 1},
	{"Defense +2", 40, 0, 2},
	{"Defense +3", 80, 0, 3},
}
