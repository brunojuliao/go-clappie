package parties

import (
	"strings"
	"testing"
)

func TestRollDice(t *testing.T) {
	result, err := Roll([]string{"1d6"})
	if err != nil {
		t.Fatalf("Roll(1d6): %v", err)
	}
	if !strings.HasPrefix(result, "1d6:") {
		t.Errorf("unexpected result format: %q", result)
	}
}

func TestRollDiceWithModifier(t *testing.T) {
	result, err := Roll([]string{"2d6+3"})
	if err != nil {
		t.Fatalf("Roll(2d6+3): %v", err)
	}
	if !strings.Contains(result, "+3") {
		t.Errorf("expected modifier in result: %q", result)
	}
}

func TestRollCoin(t *testing.T) {
	result, err := Roll([]string{"coin"})
	if err != nil {
		t.Fatalf("Roll(coin): %v", err)
	}
	if result != "Heads" && result != "Tails" {
		t.Errorf("unexpected coin result: %q", result)
	}
}

func TestRollPick(t *testing.T) {
	result, err := Roll([]string{"pick", "a", "b", "c"})
	if err != nil {
		t.Fatalf("Roll(pick): %v", err)
	}
	if result != "a" && result != "b" && result != "c" {
		t.Errorf("unexpected pick result: %q", result)
	}
}

func TestRollInvalid(t *testing.T) {
	_, err := Roll([]string{"invalid"})
	if err == nil {
		t.Error("expected error for invalid dice expression")
	}
}

func TestRollEmpty(t *testing.T) {
	_, err := Roll([]string{})
	if err == nil {
		t.Error("expected error for empty args")
	}
}
