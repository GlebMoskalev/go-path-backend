package solution

import "testing"

func TestRuneCountASCII(t *testing.T) {
	got := RuneCount("hello")
	want := 5
	if got != want {
		t.Errorf("RuneCount(%q) = %v, want %v", "hello", got, want)
	}
}

func TestRuneCountCyrillic(t *testing.T) {
	got := RuneCount("Привет")
	want := 6
	if got != want {
		t.Errorf("RuneCount(%q) = %v, want %v", "Привет", got, want)
	}
}

func TestRuneCountEmpty(t *testing.T) {
	got := RuneCount("")
	want := 0
	if got != want {
		t.Errorf("RuneCount(%q) = %v, want %v", "", got, want)
	}
}

func TestRuneCountEmoji(t *testing.T) {
	got := RuneCount("Go 🚀")
	want := 5
	if got != want {
		t.Errorf("RuneCount(%q) = %v, want %v", "Go 🚀", got, want)
	}
}
