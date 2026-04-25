package solution

import (
	"reflect"
	"testing"
)

func TestParallelMapSquare(t *testing.T) {
	got := ParallelMap([]int{1, 2, 3, 4}, func(n int) int { return n * n })
	want := []int{1, 4, 9, 16}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap squares = %v, want %v", got, want)
	}
}

func TestParallelMapDouble(t *testing.T) {
	got := ParallelMap([]int{1, 2, 3}, func(n int) int { return n * 2 })
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap double = %v, want %v", got, want)
	}
}

func TestParallelMapEmpty(t *testing.T) {
	got := ParallelMap([]int{}, func(n int) int { return n })
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap([]) = %v, want %v", got, want)
	}
}

func TestParallelMapSingle(t *testing.T) {
	got := ParallelMap([]int{5}, func(n int) int { return n + 1 })
	want := []int{6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelMap([5]) = %v, want %v", got, want)
	}
}
