package solution

import "testing"

func TestFormatGreetingBasic(t *testing.T) {
	got := FormatGreeting("Иван", 25)
	want := "Привет, Иван! Тебе 25 лет."
	if got != want {
		t.Errorf("FormatGreeting(%q, 25) = %q, want %q", "Иван", got, want)
	}
}

func TestFormatGreetingAnotherName(t *testing.T) {
	got := FormatGreeting("Анна", 1)
	want := "Привет, Анна! Тебе 1 лет."
	if got != want {
		t.Errorf("FormatGreeting(%q, 1) = %q, want %q", "Анна", got, want)
	}
}

func TestFormatGreetingEmptyName(t *testing.T) {
	got := FormatGreeting("", 0)
	want := "Привет, ! Тебе 0 лет."
	if got != want {
		t.Errorf("FormatGreeting(%q, 0) = %q, want %q", "", got, want)
	}
}
