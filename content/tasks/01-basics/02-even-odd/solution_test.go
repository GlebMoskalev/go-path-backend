package solution

import "testing"

func TestEvenOrOddEven(t *testing.T) {
	got := EvenOrOdd(2)
	if got != "even" {
		t.Errorf("EvenOrOdd(2) = %q, want \"even\"", got)
	}
}

func TestEvenOrOddOdd(t *testing.T) {
	got := EvenOrOdd(3)
	if got != "odd" {
		t.Errorf("EvenOrOdd(3) = %q, want \"odd\"", got)
	}
}

func TestEvenOrOddZero(t *testing.T) {
	got := EvenOrOdd(0)
	if got != "even" {
		t.Errorf("EvenOrOdd(0) = %q, want \"even\"", got)
	}
}

func TestEvenOrOddNegative(t *testing.T) {
	got := EvenOrOdd(-1)
	if got != "odd" {
		t.Errorf("EvenOrOdd(-1) = %q, want \"odd\"", got)
	}
}

func TestEvenOrOddNegativeEven(t *testing.T) {
	got := EvenOrOdd(-4)
	if got != "even" {
		t.Errorf("EvenOrOdd(-4) = %q, want \"even\"", got)
	}
}
