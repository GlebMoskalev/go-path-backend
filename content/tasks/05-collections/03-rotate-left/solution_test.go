package solution

import (
	"reflect"
	"testing"
)

func TestRotateLeftBasic(t *testing.T) {
	got := RotateLeft([]int{1, 2, 3, 4, 5}, 2)
	want := []int{3, 4, 5, 1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RotateLeft([1..5], 2) = %v, want %v", got, want)
	}
}

func TestRotateLeftOne(t *testing.T) {
	got := RotateLeft([]int{1, 2, 3}, 1)
	want := []int{2, 3, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RotateLeft([1,2,3], 1) = %v, want %v", got, want)
	}
}

func TestRotateLeftFullRotation(t *testing.T) {
	got := RotateLeft([]int{1, 2, 3}, 3)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RotateLeft([1,2,3], 3) = %v, want %v", got, want)
	}
}

func TestRotateLeftOverLength(t *testing.T) {
	got := RotateLeft([]int{1, 2, 3}, 4)
	want := []int{2, 3, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RotateLeft([1,2,3], 4) = %v, want %v", got, want)
	}
}

func TestRotateLeftEmpty(t *testing.T) {
	got := RotateLeft([]int{}, 2)
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("RotateLeft([], 2) = %v, want %v", got, want)
	}
}
