package solution

import (
	"math"
	"testing"
)

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 5}
	got := c.Area()
	want := math.Pi * 25
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Circle{5}.Area() = %v, want %v", got, want)
	}
}

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	got := r.Area()
	want := 12.0
	if got != want {
		t.Errorf("Rectangle{3,4}.Area() = %v, want %v", got, want)
	}
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{
		Circle{Radius: 1},
		Rectangle{Width: 2, Height: 3},
	}
	got := TotalArea(shapes)
	want := math.Pi + 6.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("TotalArea = %v, want %v", got, want)
	}
}

func TestTotalAreaEmpty(t *testing.T) {
	got := TotalArea([]Shape{})
	if got != 0 {
		t.Errorf("TotalArea([]) = %v, want 0", got)
	}
}
