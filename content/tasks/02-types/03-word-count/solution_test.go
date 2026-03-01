package solution

import (
	"reflect"
	"testing"
)

func TestWordCountBasic(t *testing.T) {
	got := WordCount("go go go")
	want := map[string]int{"go": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount(\"go go go\") = %v, want %v", got, want)
	}
}

func TestWordCountMixedCase(t *testing.T) {
	got := WordCount("Hello hello World")
	want := map[string]int{"hello": 2, "world": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount(\"Hello hello World\") = %v, want %v", got, want)
	}
}

func TestWordCountEmpty(t *testing.T) {
	got := WordCount("")
	if len(got) != 0 {
		t.Errorf("WordCount(\"\") = %v, want empty map", got)
	}
}

func TestWordCountSingleWord(t *testing.T) {
	got := WordCount("test")
	want := map[string]int{"test": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount(\"test\") = %v, want %v", got, want)
	}
}

func TestWordCountMultipleSpaces(t *testing.T) {
	got := WordCount("a  b  a")
	want := map[string]int{"a": 2, "b": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount(\"a  b  a\") = %v, want %v", got, want)
	}
}
