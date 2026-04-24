package solution

import (
	"reflect"
	"testing"
)

func TestGroupByFirstLetterBasic(t *testing.T) {
	got := GroupByFirstLetter([]string{"apple", "ant", "banana", "ax"})
	want := map[string][]string{
		"a": {"apple", "ant", "ax"},
		"b": {"banana"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupByFirstLetter = %v, want %v", got, want)
	}
}

func TestGroupByFirstLetterSkipEmpty(t *testing.T) {
	got := GroupByFirstLetter([]string{"go", ""})
	want := map[string][]string{"g": {"go"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupByFirstLetter = %v, want %v", got, want)
	}
}

func TestGroupByFirstLetterEmpty(t *testing.T) {
	got := GroupByFirstLetter([]string{})
	want := map[string][]string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupByFirstLetter([]) = %v, want %v", got, want)
	}
}

func TestGroupByFirstLetterSingleGroup(t *testing.T) {
	got := GroupByFirstLetter([]string{"cat", "car", "cup"})
	want := map[string][]string{"c": {"cat", "car", "cup"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupByFirstLetter = %v, want %v", got, want)
	}
}
