package solution

import "testing"

func TestSumBasic(t *testing.T) {
	got := Sum([]int{1, 2, 3})
	if got != 6 {
		t.Errorf("Sum([1,2,3]) = %d, want 6", got)
	}
}

func TestSumEmpty(t *testing.T) {
	got := Sum([]int{})
	if got != 0 {
		t.Errorf("Sum([]) = %d, want 0", got)
	}
}

func TestSumNil(t *testing.T) {
	got := Sum(nil)
	if got != 0 {
		t.Errorf("Sum(nil) = %d, want 0", got)
	}
}

func TestSumNegative(t *testing.T) {
	got := Sum([]int{-1, -2, -3})
	if got != -6 {
		t.Errorf("Sum([-1,-2,-3]) = %d, want -6", got)
	}
}

func TestSumMixed(t *testing.T) {
	got := Sum([]int{-10, 5, 5})
	if got != 0 {
		t.Errorf("Sum([-10,5,5]) = %d, want 0", got)
	}
}
