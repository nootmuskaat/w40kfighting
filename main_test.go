package main

import (
	"testing"
)

func TestBasicFighterApplyingActions(t *testing.T) {
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

	f1.applyAction(RegularAttack, f2)
	if f1.health != 97 {
		t.Errorf("Expected health of 97, got %d", f1.health)
	}

	f1.applyAction(CriticalAttack, f2)
	if f1.health != 92 {
		t.Errorf("Expected health of 92, got %d", f1.health)
	}

	f1.applyAction(CriticalBlock, f2)
	if f1.roll.crits != 0 && f1.roll.hits != 2 {
		t.Errorf("Expected roll to be 0&2, got %d&%d", f1.roll.hits, f1.roll.crits)
	}

	f1.applyAction(CriticalBlock, f2)
	if f1.roll.crits != 0 && f1.roll.hits != 1 {
		t.Errorf("Expected roll to be 0&1, got %d&%d", f1.roll.hits, f1.roll.crits)
	}

	f1.applyAction(RegularBlock, f2)
	if f1.roll.crits != 0 && f1.roll.hits != 0 {
		t.Errorf("Expected roll to be 0&0, got %d&%d", f1.roll.hits, f1.roll.crits)
	}
}

func TestBasicFighterPerformingActions(t *testing.T) {
	f := fighter {
		health: 100,
		initiative: 1,
		weapon: weapon { 3, 5 },
		roll: roll {10, 10},
	}

	f.performAction(RegularAttack)
	if f.roll.hits != 9 {
		t.Errorf("Expected rolled hits to be at 9, got %d", f.roll.hits)
	}

	f.performAction(RegularBlock)
	if f.roll.hits != 8 {
		t.Errorf("Expected rolled hits to be at 8, got %d", f.roll.hits)
	}

	f.performAction(CriticalAttack)
	if f.roll.crits != 9 {
		t.Errorf("Expected rolled crits to be at 9, got %d", f.roll.crits)
	}

	f.performAction(CriticalBlock)
	if f.roll.crits != 8 {
		t.Errorf("Expected rolled crits to be at 8, got %d", f.roll.crits)
	}
}
