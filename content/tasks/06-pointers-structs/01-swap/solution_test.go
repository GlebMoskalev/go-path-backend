package solution

import "testing"

func TestSwapBasic(t *testing.T) {
	a, b := 10, 20
	Swap(&a, &b)
	if a != 20 || b != 10 {
		t.Errorf("Swap: got a=%v, b=%v, want a=20, b=10", a, b)
	}
}

func TestSwapEqual(t *testing.T) {
	a, b := 5, 5
	Swap(&a, &b)
	if a != 5 || b != 5 {
		t.Errorf("Swap: got a=%v, b=%v, want a=5, b=5", a, b)
	}
}

func TestSwapZeroAndOne(t *testing.T) {
	a, b := 0, 1
	Swap(&a, &b)
	if a != 1 || b != 0 {
		t.Errorf("Swap: got a=%v, b=%v, want a=1, b=0", a, b)
	}
}
