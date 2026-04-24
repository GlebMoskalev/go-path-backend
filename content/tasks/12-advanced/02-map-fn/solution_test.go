package solution

import (
	"reflect"
	"testing"
)

func TestMapFnSquare(t *testing.T) {
	got := MapFn([]int{1, 2, 3}, func(n int) int { return n * n })
	want := []int{1, 4, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapFn squares = %v, want %v", got, want)
	}
}

func TestMapFnStringLen(t *testing.T) {
	got := MapFn([]string{"hi", "bye"}, func(s string) int { return len(s) })
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapFn string len = %v, want %v", got, want)
	}
}

func TestMapFnEmpty(t *testing.T) {
	got := MapFn([]int{}, func(n int) int { return n })
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapFn([]) = %v, want %v", got, want)
	}
}

func TestMapFnIntToString(t *testing.T) {
	got := MapFn([]int{1, 2}, func(n int) bool { return n > 1 })
	want := []bool{false, true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MapFn int to bool = %v, want %v", got, want)
	}
}
