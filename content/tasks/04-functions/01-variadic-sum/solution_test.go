package solution

import "testing"

func TestSumMultiple(t *testing.T) {
	got := Sum(1, 2, 3)
	want := 6
	if got != want {
		t.Errorf("Sum(1, 2, 3) = %v, want %v", got, want)
	}
}

func TestSumWithNegative(t *testing.T) {
	got := Sum(10, -5, 3)
	want := 8
	if got != want {
		t.Errorf("Sum(10, -5, 3) = %v, want %v", got, want)
	}
}

func TestSumEmpty(t *testing.T) {
	got := Sum()
	want := 0
	if got != want {
		t.Errorf("Sum() = %v, want %v", got, want)
	}
}

func TestSumSingle(t *testing.T) {
	got := Sum(42)
	want := 42
	if got != want {
		t.Errorf("Sum(42) = %v, want %v", got, want)
	}
}
