package solution

import (
	"math"
	"testing"
)

func TestDistancePythagorean(t *testing.T) {
	a := Point{0, 0}
	b := Point{3, 4}
	got := a.Distance(b)
	want := 5.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Distance((0,0), (3,4)) = %v, want %v", got, want)
	}
}

func TestDistanceSamePoint(t *testing.T) {
	a := Point{1, 1}
	got := a.Distance(a)
	want := 0.0
	if got != want {
		t.Errorf("Distance to same point = %v, want %v", got, want)
	}
}

func TestDistanceHorizontal(t *testing.T) {
	a := Point{0, 0}
	b := Point{1, 0}
	got := a.Distance(b)
	want := 1.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Distance((0,0), (1,0)) = %v, want %v", got, want)
	}
}

func TestDistanceNegativeCoords(t *testing.T) {
	a := Point{-1, -1}
	b := Point{2, 3}
	got := a.Distance(b)
	want := 5.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Distance((-1,-1), (2,3)) = %v, want %v", got, want)
	}
}
