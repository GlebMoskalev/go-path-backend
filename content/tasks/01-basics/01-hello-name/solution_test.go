package solution

import "testing"

func TestHelloBasic(t *testing.T) {
	got := Hello("Иван")
	want := "Привет, Иван!"
	if got != want {
		t.Errorf("Hello(\"Иван\") = %q, want = %q", got, want)
	}
}

func TestHelloEnglish(t *testing.T) {
	got := Hello("World")
	want := "Привет, World!"
	if got != want {
		t.Errorf("Hello(\"World\") = %q, want %q", got, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	got := Hello("")
	want := "Привет, !"
	if got != want {
		t.Errorf("Hello(\"\") = %q, want %q", got, want)
	}
}
