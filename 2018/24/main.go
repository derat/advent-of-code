package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/derat/advent-of-code/lib"
)

const debug = false

func main() {
	var grps []*group
	for _, lns := range lib.InputParagraphs("2018/24") {
		team := team(strings.TrimRight(lns[0], ":"))
		for i, ln := range lns[1:] {
			grp := &group{
				team:   team,
				id:     i + 1,
				weak:   make(map[string]struct{}),
				immune: make(map[string]struct{}),
			}
			var attrs string // e.g. "weak to bludgeoning; immune to slashing, cold"
			lib.Extract(ln, `^(\d+) units each with (\d+) hit points(?: \((.+)\))? `+
				`with an attack that does (\d+) ([a-z]+) damage at initiative (\d+)$`,
				&grp.units, &grp.hp, &attrs, &grp.dp, &grp.dt, &grp.init)
			if attrs != "" {
				for _, p := range strings.Split(attrs, "; ") {
					var prop, types string
					lib.Extract(p, `^(weak|immune) to (.+)$`, &prop, &types)
					m := &grp.weak
					if prop == "immune" {
						m = &grp.immune
					}
					for _, t := range strings.Split(types, ", ") {
						(*m)[t] = struct{}{}
					}
				}
			}
			grps = append(grps, grp)
		}
	}

	var rounds int
	for {
		rounds++

		// Check if a team has been wiped out.
		var done bool
		for _, t := range []team{immune, infect} {
			if debug {
				fmt.Println(t + ":")
			}
			var cnt int
			for _, g := range grps {
				if g.alive() && g.team == t {
					cnt++
					if debug {
						fmt.Printf("Group %d contains %d units\n", g.id, g.units)
					}
				}
			}
			if cnt == 0 {
				done = true
				if debug {
					fmt.Println("No groups remain.")
				}
			}
		}
		if done {
			break
		}
		if debug {
			fmt.Println()
		}

		// Target selection phase.
		// The order between teams shouldn't matter here, but I'm trying to match
		// the example output from the problem to make comparisons easier.
		targets := make(map[*group]*group)    // keys are attackers, vals are defenders
		selected := make(map[*group]struct{}) // already-targeted groups
		for _, t := range []team{infect, immune} {
			for _, att := range targetOrder(grps, t) {
				def, dam := att.target(grps, selected)
				if def == nil {
					continue
				}
				// Can't cache the damage since some units might be killed by the
				// time that we get to attack.
				targets[att] = def
				selected[def] = struct{}{}
				if debug {
					fmt.Printf("%s group %d would deal defending group %d %d damage\n",
						att.team, att.id, def.id, dam)
				}
			}
		}
		if debug {
			fmt.Println()
		}

		// Attacking phase.
		for _, att := range attackOrder(grps) {
			if !att.alive() {
				continue // attacker was killed earlier in the round
			}
			def := targets[att]
			if def == nil || !def.alive() {
				continue // defender was killed earlier in the round
			}
			// "The defending group only loses whole units from damage; damage is always dealt in
			// such a way that it kills the most units possible, and any remaining damage to a unit
			// that does not immediately kill it is ignored."
			dam := att.damageAgainst(def)
			killed := lib.Min(dam/def.hp, def.units)
			def.units -= killed
			if debug {
				fmt.Printf("%s group %d attacks defending group %d, killing %d units\n",
					att.team, att.id, def.id, killed)
			}
		}
		if debug {
			fmt.Println()
		}
	}

	// Part 1: Print number of units belonging to winning army.
	var cnt int
	for _, g := range grps {
		if g.alive() {
			cnt += g.units
		}
	}
	fmt.Println(cnt)
}

type team string

const (
	immune team = "Immune System"
	infect      = "Infection"
)

func (t *team) String() string {
	return string(*t)
}

func (t *team) opp() team {
	if *t == immune {
		return infect
	}
	return immune
}

type group struct {
	team  team
	id    int
	units int
	hp    int    // per unit
	init  int    // initiative
	dt    string // damage type, e.g. "fire"
	dp    int    // damage points

	weak, immune map[string]struct{}
}

func (g *group) alive() bool {
	return g.units > 0
}

func (g *group) ep() int {
	// "Each group also has an effective power: the number of units in that group multiplied by
	// their attack damage."
	return g.units * g.dp
}

func (g *group) damageAgainst(def *group) int {
	// "The damage an attacking group deals to a defending group depends on the attacking group's
	// attack type and the defending group's immunities and weaknesses. By default, an attacking
	// group would deal damage equal to its effective power to the defending group. However, if the
	// defending group is immune to the attacking group's attack type, the defending group instead
	// takes no damage; if the defending group is weak to the attacking group's attack type, the
	// defending group instead takes double damage."
	if lib.MapHasKey(def.immune, g.dt) {
		return 0
	}
	dam := g.ep()
	if lib.MapHasKey(def.weak, g.dt) {
		dam *= 2
	}
	return dam
}

// target returns the (possibly nil) group from grps that g will target.
// Already-targeted groups in sel are excluded from consideration.
func (g *group) target(grps []*group, sel map[*group]struct{}) (*group, int) {
	var ret []*group
	for _, def := range grps {
		if def.alive() && def.team != g.team && !lib.MapHasKey(sel, def) {
			ret = append(ret, def)
		}
	}
	if len(ret) == 0 {
		return nil, 0
	}

	// "The attacking group chooses to target the group in the enemy army to which it would deal the
	// most damage (after accounting for weaknesses and immunities, but not accounting for whether
	// the defending group has enough units to actually receive all of that damage).
	//
	// If an attacking group is considering two defending groups to which it would deal equal
	// damage, it chooses to target the defending group with the largest effective power; if there
	// is still a tie, it chooses the defending group with the highest initiative."
	dam := make(map[*group]int, len(ret))
	for _, def := range ret {
		dam[def] = g.damageAgainst(def)
	}
	sort.Slice(ret, func(i, j int) bool {
		a, b := ret[i], ret[j]
		if da, db := dam[a], dam[b]; da > db {
			return true
		} else if da < db {
			return false
		}
		if ae, be := a.ep(), b.ep(); ae > be {
			return true
		} else if ae < be {
			return false
		}
		return a.init > b.init
	})
	return ret[0], dam[ret[0]]
}

// targetOrder returns living groups from the specified team in the order
// in which they should choose their targets during the target selection phase.
func targetOrder(grps []*group, t team) []*group {
	var ret []*group
	for _, g := range grps {
		if g.alive() && g.team == t {
			ret = append(ret, g)
		}
	}
	// "In decreasing order of effective power, groups choose their targets;
	// in a tie, the group with the higher initiative chooses first."
	sort.Slice(ret, func(i, j int) bool {
		a, b := ret[i], ret[j]
		if ae, be := a.ep(), b.ep(); ae > be {
			return true
		} else if ae < be {
			return false
		}
		return a.init > b.init
	})
	return ret
}

// attackOrder returns living groups in the order in which they should attack
// during the attacking phase.
func attackOrder(grps []*group) []*group {
	var ret []*group
	for _, g := range grps {
		if g.alive() {
			ret = append(ret, g)
		}
	}
	// "Groups attack in decreasing order of initiative, regardless of whether they are part of the
	// infection or the immune system. (If a group contains no units, it cannot attack.)"
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].init > ret[j].init
	})
	return ret
}
