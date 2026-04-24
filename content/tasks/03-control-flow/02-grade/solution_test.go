package solution

import "testing"

func TestGradeA(t *testing.T) {
	got := Grade(95)
	want := "A"
	if got != want {
		t.Errorf("Grade(95) = %q, want %q", got, want)
	}
}

func TestGradeB(t *testing.T) {
	got := Grade(85)
	want := "B"
	if got != want {
		t.Errorf("Grade(85) = %q, want %q", got, want)
	}
}

func TestGradeC(t *testing.T) {
	got := Grade(75)
	want := "C"
	if got != want {
		t.Errorf("Grade(75) = %q, want %q", got, want)
	}
}

func TestGradeF(t *testing.T) {
	got := Grade(45)
	want := "F"
	if got != want {
		t.Errorf("Grade(45) = %q, want %q", got, want)
	}
}

func TestGradeZero(t *testing.T) {
	got := Grade(0)
	want := "F"
	if got != want {
		t.Errorf("Grade(0) = %q, want %q", got, want)
	}
}
