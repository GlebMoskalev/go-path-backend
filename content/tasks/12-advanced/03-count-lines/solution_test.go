package solution

import (
	"strings"
	"testing"
)

func TestCountLinesMultiple(t *testing.T) {
	r := strings.NewReader("line1\nline2\nline3")
	got, err := CountLines(r)
	if err != nil || got != 3 {
		t.Errorf("CountLines(3 lines) = (%v, %v), want (3, nil)", got, err)
	}
}

func TestCountLinesSingle(t *testing.T) {
	r := strings.NewReader("single")
	got, err := CountLines(r)
	if err != nil || got != 1 {
		t.Errorf("CountLines(single) = (%v, %v), want (1, nil)", got, err)
	}
}

func TestCountLinesEmpty(t *testing.T) {
	r := strings.NewReader("")
	got, err := CountLines(r)
	if err != nil || got != 0 {
		t.Errorf("CountLines('') = (%v, %v), want (0, nil)", got, err)
	}
}

func TestCountLinesTrailingNewline(t *testing.T) {
	r := strings.NewReader("a\nb\n")
	got, err := CountLines(r)
	if err != nil || got != 2 {
		t.Errorf("CountLines('a\\nb\\n') = (%v, %v), want (2, nil)", got, err)
	}
}
