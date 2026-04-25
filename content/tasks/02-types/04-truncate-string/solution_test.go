package solution

import "testing"

func TestTruncateShort(t *testing.T) {
	got := Truncate("Hello", 10)
	want := "Hello"
	if got != want {
		t.Errorf("Truncate(%q, 10) = %q, want %q", "Hello", got, want)
	}
}

func TestTruncateLong(t *testing.T) {
	got := Truncate("Hello, World!", 5)
	want := "He..."
	if got != want {
		t.Errorf("Truncate(%q, 5) = %q, want %q", "Hello, World!", got, want)
	}
}

func TestTruncateCyrillic(t *testing.T) {
	got := Truncate("Привет, мир!", 8)
	want := "Приве..."
	if got != want {
		t.Errorf("Truncate(%q, 8) = %q, want %q", "Привет, мир!", got, want)
	}
}

func TestTruncateExact(t *testing.T) {
	got := Truncate("abc", 3)
	want := "abc"
	if got != want {
		t.Errorf("Truncate(%q, 3) = %q, want %q", "abc", got, want)
	}
}
