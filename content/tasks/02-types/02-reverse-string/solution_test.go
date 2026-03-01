package solution

import "testing"

func TestReverseEnglish(t *testing.T) {
	got := Reverse("hello")
	if got != "olleh" {
		t.Errorf("Reverse(\"hello\") = %q, want \"olleh\"", got)
	}
}

func TestReverseRussian(t *testing.T) {
	got := Reverse("Привет")
	if got != "тевирП" {
		t.Errorf("Reverse(\"Привет\") = %q, want \"тевирП\"", got)
	}
}

func TestReverseEmpty(t *testing.T) {
	got := Reverse("")
	if got != "" {
		t.Errorf("Reverse(\"\") = %q, want \"\"", got)
	}
}

func TestReverseSingleChar(t *testing.T) {
	got := Reverse("a")
	if got != "a" {
		t.Errorf("Reverse(\"a\") = %q, want \"a\"", got)
	}
}

func TestReversePalindrome(t *testing.T) {
	got := Reverse("abcba")
	if got != "abcba" {
		t.Errorf("Reverse(\"abcba\") = %q, want \"abcba\"", got)
	}
}
