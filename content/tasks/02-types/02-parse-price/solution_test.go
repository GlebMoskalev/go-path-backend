package solution

import "testing"

func TestParsePriceFloat(t *testing.T) {
	got, ok := ParsePrice("19.99")
	if !ok || got != 19.99 {
		t.Errorf("ParsePrice(%q) = (%v, %v), want (19.99, true)", "19.99", got, ok)
	}
}

func TestParsePriceInteger(t *testing.T) {
	got, ok := ParsePrice("100")
	if !ok || got != 100.0 {
		t.Errorf("ParsePrice(%q) = (%v, %v), want (100.0, true)", "100", got, ok)
	}
}

func TestParsePriceInvalidString(t *testing.T) {
	got, ok := ParsePrice("abc")
	if ok || got != 0 {
		t.Errorf("ParsePrice(%q) = (%v, %v), want (0, false)", "abc", got, ok)
	}
}

func TestParsePriceEmpty(t *testing.T) {
	got, ok := ParsePrice("")
	if ok || got != 0 {
		t.Errorf("ParsePrice(%q) = (%v, %v), want (0, false)", "", got, ok)
	}
}
