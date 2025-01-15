package main

import (
	"fmt"
	"slices"
	"strings"
)

type action int
type sequence []action

func (s sequence) String() string {
	ss := make([]string, len(s))
	for i, action := range s {
		ss[i] = action.String()
	}
	return strings.Join(ss, " -> ")
}


const (
	/* player actions */
	RegularAttack action = iota
	RegularBlock
	CriticalAttack
	CriticalBlock
	Pass  // the player cannot take any action
	/* control signals */
	Dead
	EOS  // End of Sequence
	Invalid  // No valid sequences with this prefix remain
)

func nextAction(a action) action {
	switch a {
	case RegularAttack:
		return RegularBlock
	case RegularBlock:
		return CriticalAttack
	case CriticalAttack:
		return CriticalBlock
	case CriticalBlock:
		return Pass
	default:
		return Invalid
	}
}

func (a action) String() string {
	switch a {
	case RegularAttack:
		return "🗡️"
	case RegularBlock:
		return "🛡️"
	case CriticalAttack:
		return "🔥🗡️"
	case CriticalBlock:
		return "🔥🛡️"
	case Pass:
		return "🫥"
	case Dead:
		return "💀"
	case EOS:
		return "✅"
	default:
		return "⛔"
	}

}

type weapon struct {
	normal, critical int
}

type roll struct {
	hits, crits int
}

type fighter struct {
	health int
	initiative int
	weapon weapon
	roll roll
}

func (f fighter) copy() fighter {
	return f
}

func (f fighter) checkApplyAction(a action) bool {
	switch a {
	case RegularAttack, RegularBlock:
		return f.roll.hits > 0
	case CriticalAttack, CriticalBlock:
		return f.roll.crits > 0
	case Pass:
		return f.roll.hits + f.roll.crits == 0
	default:
		return true
	}
}

func (f fighter) actionsRemaining() int {
	return f.roll.hits + f.roll.crits
}

func (f *fighter) applyAction(a action, other fighter) {
	if a == RegularAttack {
		if other.weapon.normal >= f.health {
			f.health = 0
		} else {
			f.health -= other.weapon.normal
		}
	} else if a == CriticalAttack {
		if other.weapon.critical >= f.health {
			f.health = 0
		} else {
			f.health -= other.weapon.critical
		}
	} else if a == RegularBlock {
		if f.roll.hits > 0 {
			f.roll.hits--
		}
	} else if a == CriticalBlock {
		if f.roll.crits > 0 {
			f.roll.crits--
		} else if f.roll.hits > 0 {
			f.roll.hits--
		}
	}
}


func (f *fighter) performAction(a action) {
	if a == RegularAttack || a == RegularBlock {
		f.roll.hits--
	} else if a == CriticalAttack || a == CriticalBlock {
		f.roll.crits--
	}
}


func newFighter(health, initiative, normal, critical int) fighter {
	return fighter{health, initiative, weapon{normal, critical}, *new(roll)}
}

func (f fighter) String() string {
	return fmt.Sprintf(
		"❤️ %d | ⚡ %d | 🗡️ %d 🔥 %d | 🎲 🗡️ %d 🔥 %d",
		f.health,
		f.initiative,
		f.weapon.normal,
		f.weapon.critical,
		f.roll.hits,
		f.roll.crits,
	)
}

func runPossiblities(f1, f2 *fighter, previous, current sequence) sequence {
	var f *fighter
	var o *fighter
	idx := len(current)
	if idx % 2 == 0 {
		f, o = f1, f2
	} else {
		f, o = f2, f1
	}
	if f.health == 0 {
		return append(current, Dead)
	} else if f.actionsRemaining() + o.actionsRemaining() == 0 {
		return append(current, EOS)
	}
	if previous == nil || len(previous) <= idx + 1 {
		for a := RegularAttack; a != Invalid; a = nextAction(a) {
			if ! f.checkApplyAction(a) {
				continue
			}
			current = append(current, a)
			f.performAction(a)
			o.applyAction(a, *f)
			// fmt.Println("current", current)
			return runPossiblities(f1, f2, previous, current)
		}
	} else {
		switch previous[idx+1] {
		case Dead, EOS, Invalid:
			for a := nextAction(previous[idx]); a != Invalid; a = nextAction(a) {
				if ! f.checkApplyAction(a) {
					continue
				}
				current = append(current, a)
				f.performAction(a)
				o.applyAction(a, *f)
				// fmt.Println("current", current)
				return runPossiblities(f1, f2, previous, current)
			}
		default:
			a := previous[idx]
			current = append(current, a)
			f.performAction(a)
			o.applyAction(a, *f)
			// fmt.Println("current", current)
			return runPossiblities(f1, f2, previous, current)

		}

	}
	return append(current, Invalid)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkAllPossibilities(f1, f2 fighter) []sequence {
	var fr1 fighter
	var fr2 fighter
	if f2.initiative > f1.initiative {
		fr2, fr1 = f1.copy(), f2.copy()
	} else {
		fr1, fr2 = f1.copy(), f2.copy()
	}
	allSequences := make([]sequence, 0, 8)
	var previous sequence = nil
	endState := sequence{Invalid}

	for seq := runPossiblities(&fr1, &fr2, previous, make(sequence, 0)); ! slices.Equal(seq, endState); {

		// fmt.Println(" final", seq)
		// fmt.Println()
		if seq[len(seq)-1] != Invalid {
			allSequences = append(allSequences, seq)
		}
		if f2.initiative > f1.initiative {
			fr2, fr1 = f1.copy(), f2.copy()
		} else {
			fr1, fr2 = f1.copy(), f2.copy()
		}
		seq = runPossiblities(&fr1, &fr2, seq, make(sequence, 0))
	}

	return allSequences
}


func main() {
	f1 := fighter {
		health: 100,
		initiative: 1,
		weapon: weapon { 3, 5 },
		roll: roll {2, 1},
	}
	f2 := fighter {
		health: 100,
		initiative: 1,
		weapon: weapon { 3, 5 },
		roll: roll {2, 2},
	}
	allPossibilities := checkAllPossibilities(f1, f2)
	for _, possible := range allPossibilities {
		fmt.Println(possible)
	}
}

