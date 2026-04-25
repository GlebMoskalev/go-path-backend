package solution

import (
	"reflect"
	"testing"
)

func TestIntersectBasic(t *testing.T) {
	got := Intersect([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
	want := []int{3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect([1,2,3,4], [3,4,5,6]) = %v, want %v", got, want)
	}
}

func TestIntersectDuplicatesInInput(t *testing.T) {
	got := Intersect([]int{1, 2, 2, 3}, []int{2, 3, 3})
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect([1,2,2,3], [2,3,3]) = %v, want %v", got, want)
	}
}

func TestIntersectNoCommon(t *testing.T) {
	got := Intersect([]int{1, 2, 3}, []int{4, 5, 6})
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect([1,2,3], [4,5,6]) = %v, want %v", got, want)
	}
}

func TestIntersectEmptyFirst(t *testing.T) {
	got := Intersect([]int{}, []int{1, 2})
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Intersect([], [1,2]) = %v, want %v", got, want)
	}
}
