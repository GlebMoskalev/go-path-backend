package solution

import "testing"

func TestDescribeInt(t *testing.T) {
	got := Describe(42)
	want := "int: 42"
	if got != want {
		t.Errorf("Describe(42) = %q, want %q", got, want)
	}
}

func TestDescribeString(t *testing.T) {
	got := Describe("hello")
	want := "string: hello"
	if got != want {
		t.Errorf("Describe(%q) = %q, want %q", "hello", got, want)
	}
}

func TestDescribeBool(t *testing.T) {
	got := Describe(true)
	want := "bool: true"
	if got != want {
		t.Errorf("Describe(true) = %q, want %q", got, want)
	}
}

func TestDescribeSlice(t *testing.T) {
	got := Describe([]int{1, 2, 3})
	want := "[]int с 3 элементами"
	if got != want {
		t.Errorf("Describe([]int{1,2,3}) = %q, want %q", got, want)
	}
}

func TestDescribeUnknown(t *testing.T) {
	got := Describe(3.14)
	want := "неизвестный тип"
	if got != want {
		t.Errorf("Describe(3.14) = %q, want %q", got, want)
	}
}
