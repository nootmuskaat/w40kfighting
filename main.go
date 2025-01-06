package main

import (
	"fmt"
)

type action int

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


func main() {
	f := newFighter(22, 3, 3, 5)
	fmt.Println(f)
}

