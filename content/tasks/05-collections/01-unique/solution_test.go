package solution

import (
	"reflect"
	"testing"
)

func TestUniqueDuplicates(t *testing.T) {
	got := Unique([]int{1, 2, 2, 3, 1})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique([1,2,2,3,1]) = %v, want %v", got, want)
	}
}

func TestUniqueAllSame(t *testing.T) {
	got := Unique([]int{5, 5, 5})
	want := []int{5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique([5,5,5]) = %v, want %v", got, want)
	}
}

func TestUniqueNoDuplicates(t *testing.T) {
	got := Unique([]int{1, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique([1,2,3]) = %v, want %v", got, want)
	}
}

func TestUniqueEmpty(t *testing.T) {
	got := Unique([]int{})
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Unique([]) = %v, want %v", got, want)
	}
}
