package main

import (
	"fmt"
	"math"

	"github.com/derat/advent-of-code/lib"
)

func main() {
	init := state{hp: 50, mana: 500} // given in description
	lib.Extract(lib.Input("2015/22"), `^Hit Points: (\d+)\nDamage: (\d+)\n$`, &init.bossHP, &init.bossDamage)

	best := state{spent: math.MaxInt32}

	var playerTurn, bossTurn func(state)

	playerTurn = func(s state) {
		// Part 2: Hard mode
		s.hp -= s.turnDamage
		if s.lost() {
			return
		}

		s.update()
		if s.won() {
			if s.spent < best.spent {
				best = s
			}
			return
		}
		for _, sp := range []spell{magicMissle, drain, shield, poison, recharge} {
			if ns, ok := sp(s); ok && ns.spent < best.spent {
				bossTurn(ns)
			}
		}
	}

	bossTurn = func(s state) {
		s.update()
		if s.won() {
			if s.spent < best.spent {
				best = s
			}
			return
		}
		s.hp -= lib.Max(s.bossDamage-s.armor(), 1)
		if !s.lost() {
			playerTurn(s)
		}
	}

	playerTurn(init)
	fmt.Println(best.spent)

	// Part 2: Player loses one HP at the beginning of each player turn.
	init.turnDamage = 1
	best = state{spent: math.MaxInt32}
	playerTurn(init)
	fmt.Println(best.spent)
}

type state struct {
	hp, mana, spent          int
	bossHP, bossDamage       int
	shield, poison, recharge int // turns remaining for spells (0 if inactive)
	turnDamage               int // hp lost at start of player turns (part 2)
}

func (s *state) won() bool {
	return s.hp > 0 && s.bossHP <= 0
}

func (s *state) lost() bool {
	return s.hp <= 0
}

func (s *state) spend(mana int) bool {
	if s.mana < mana {
		return false
	}
	s.mana -= mana
	s.spent += mana
	return true
}

func (s *state) armor() int {
	if s.shield > 0 {
		return 7
	}
	return 0
}

// update updates s for the start of a turn.
func (s *state) update() {
	if s.poison > 0 {
		s.bossHP -= 3
	}
	if s.recharge > 0 {
		s.mana += 101
	}

	if s.shield > 0 {
		s.shield--
	}
	if s.poison > 0 {
		s.poison--
	}
	if s.recharge > 0 {
		s.recharge--
	}
}

type spell func(s state) (state, bool)

func magicMissle(s state) (state, bool) {
	if !s.spend(53) {
		return s, false
	}
	s.bossHP -= 4
	return s, true
}

func drain(s state) (state, bool) {
	if !s.spend(73) {
		return s, false
	}
	s.hp += 2
	s.bossHP -= 2
	return s, true
}

func shield(s state) (state, bool) {
	if s.shield > 0 || !s.spend(113) {
		return s, false
	}
	s.shield = 6
	return s, true
}

func poison(s state) (state, bool) {
	if s.poison > 0 || !s.spend(173) {
		return s, false
	}
	s.poison = 6
	return s, true
}

func recharge(s state) (state, bool) {
	if s.recharge > 0 || !s.spend(229) {
		return s, false
	}
	s.recharge = 5
	return s, true
}
