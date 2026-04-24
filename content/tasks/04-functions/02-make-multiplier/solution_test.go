package solution

import "testing"

func TestMakeMultiplierDouble(t *testing.T) {
	double := MakeMultiplier(2)
	got := double(5)
	want := 10
	if got != want {
		t.Errorf("MakeMultiplier(2)(5) = %v, want %v", got, want)
	}
}

func TestMakeMultiplierTriple(t *testing.T) {
	triple := MakeMultiplier(3)
	got := triple(5)
	want := 15
	if got != want {
		t.Errorf("MakeMultiplier(3)(5) = %v, want %v", got, want)
	}
}

func TestMakeMultiplierZeroFactor(t *testing.T) {
	f := MakeMultiplier(0)
	got := f(99)
	want := 0
	if got != want {
		t.Errorf("MakeMultiplier(0)(99) = %v, want %v", got, want)
	}
}

func TestMakeMultiplierZeroInput(t *testing.T) {
	double := MakeMultiplier(2)
	got := double(0)
	want := 0
	if got != want {
		t.Errorf("MakeMultiplier(2)(0) = %v, want %v", got, want)
	}
}
