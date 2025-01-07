package main

import (
	"fmt"
)

type action int
type sequence []action

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
	case Dead:
		return "💀"
	default:
		return "-"
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

func runPossiblities(f1, f2 fighter, allSequences *[]sequence, current *sequence) {
	var f fighter
	if len(*current) % 2 == 0 {
		f = f1
	} else {
		f = f2
	}
	f.performAction(Pass)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func checkAllPossibilities(f1, f2 fighter) []sequence {
	if f2.initiative > f1.initiative {
		f2, f1 = f1, f2
	}
	allSequences := make([]sequence, 0, 8)
	maxPossibleSequence := max(f1.roll.hits + f1.roll.crits, f2.roll.hits + f2.roll.crits) * 2 + 1
	currentSequence := make(sequence, 0, maxPossibleSequence)
	runPossiblities(f1, f2, &allSequences, &currentSequence)

	return allSequences
}


func main() {
}

