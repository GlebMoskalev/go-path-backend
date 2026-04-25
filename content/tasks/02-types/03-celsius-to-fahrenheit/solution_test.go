package solution

import "testing"

func TestToFahrenheitFreezing(t *testing.T) {
	got := ToFahrenheit(0)
	want := Fahrenheit(32)
	if got != want {
		t.Errorf("ToFahrenheit(0) = %v, want %v", got, want)
	}
}

func TestToFahrenheitBoiling(t *testing.T) {
	got := ToFahrenheit(100)
	want := Fahrenheit(212)
	if got != want {
		t.Errorf("ToFahrenheit(100) = %v, want %v", got, want)
	}
}

func TestToFahrenheitNegative(t *testing.T) {
	got := ToFahrenheit(-40)
	want := Fahrenheit(-40)
	if got != want {
		t.Errorf("ToFahrenheit(-40) = %v, want %v", got, want)
	}
}

func TestToFahrenheitBody(t *testing.T) {
	got := ToFahrenheit(37)
	want := Fahrenheit(98.6)
	if got != want {
		t.Errorf("ToFahrenheit(37) = %v, want %v", got, want)
	}
}
