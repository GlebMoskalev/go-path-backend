package solution

import "testing"

func TestRunLengthMixed(t *testing.T) {
	got := RunLength("aaabbc")
	want := "3a2bc"
	if got != want {
		t.Errorf("RunLength(%q) = %q, want %q", "aaabbc", got, want)
	}
}

func TestRunLengthNoDuplicates(t *testing.T) {
	got := RunLength("abc")
	want := "abc"
	if got != want {
		t.Errorf("RunLength(%q) = %q, want %q", "abc", got, want)
	}
}

func TestRunLengthEmpty(t *testing.T) {
	got := RunLength("")
	want := ""
	if got != want {
		t.Errorf("RunLength(%q) = %q, want %q", "", got, want)
	}
}

func TestRunLengthAllSame(t *testing.T) {
	got := RunLength("aaaa")
	want := "4a"
	if got != want {
		t.Errorf("RunLength(%q) = %q, want %q", "aaaa", got, want)
	}
}
