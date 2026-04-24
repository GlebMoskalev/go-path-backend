package solution

import (
	"reflect"
	"testing"
)

func TestFilterEvenInts(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
	want := []int{2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter evens = %v, want %v", got, want)
	}
}

func TestFilterStrings(t *testing.T) {
	got := Filter([]string{"go", "rust", "java"}, func(s string) bool { return len(s) == 2 })
	want := []string{"go"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter short strings = %v, want %v", got, want)
	}
}

func TestFilterEmpty(t *testing.T) {
	got := Filter([]int{}, func(n int) bool { return true })
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter([]) = %v, want %v", got, want)
	}
}

func TestFilterNoneMatch(t *testing.T) {
	got := Filter([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 })
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filter no matches = %v, want %v", got, want)
	}
}
