package solution

import (
	"reflect"
	"testing"
)

func TestWordFrequencyRepeated(t *testing.T) {
	got := WordFrequency("go is go")
	want := map[string]int{"go": 2, "is": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency(%q) = %v, want %v", "go is go", got, want)
	}
}

func TestWordFrequencyUnique(t *testing.T) {
	got := WordFrequency("hello world")
	want := map[string]int{"hello": 1, "world": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency(%q) = %v, want %v", "hello world", got, want)
	}
}

func TestWordFrequencyEmpty(t *testing.T) {
	got := WordFrequency("")
	want := map[string]int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency(%q) = %v, want %v", "", got, want)
	}
}

func TestWordFrequencySingle(t *testing.T) {
	got := WordFrequency("hello")
	want := map[string]int{"hello": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency(%q) = %v, want %v", "hello", got, want)
	}
}
