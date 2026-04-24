package solution

import "testing"

func TestTableRowBasic(t *testing.T) {
	got := TableRow("Ivan", 42)
	want := "Ivan                |  42"
	if got != want {
		t.Errorf("TableRow(%q, 42) = %q, want %q", "Ivan", got, want)
	}
}

func TestTableRowLongName(t *testing.T) {
	got := TableRow("Alex Ivanov", 100)
	want := "Alex Ivanov         | 100"
	if got != want {
		t.Errorf("TableRow(%q, 100) = %q, want %q", "Alex Ivanov", got, want)
	}
}

func TestTableRowEmpty(t *testing.T) {
	got := TableRow("", 0)
	want := "                    |   0"
	if got != want {
		t.Errorf("TableRow(%q, 0) = %q, want %q", "", got, want)
	}
}
