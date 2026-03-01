package solution

import (
	"math"
	"testing"
)

func TestRectangleArea(t *testing.T) {
	r := Rectangle{Width: 3, Height: 4}
	got := r.Area()
	if got != 12 {
		t.Errorf("Rectangle{3,4}.Area() = %v, want 12", got)
	}
}

func TestCircleArea(t *testing.T) {
	c := Circle{Radius: 5}
	got := c.Area()
	want := math.Pi * 25
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Circle{5}.Area() = %v, want %v", got, want)
	}
}

func TestTotalArea(t *testing.T) {
	shapes := []Shape{
		Rectangle{Width: 3, Height: 4},
		Circle{Radius: 5},
	}
	got := TotalArea(shapes)
	want := 12 + math.Pi*25
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

func TestRectangleZero(t *testing.T) {
	r := Rectangle{Width: 0, Height: 10}
	got := r.Area()
	if got != 0 {
		t.Errorf("Rectangle{0,10}.Area() = %v, want 0", got)
	}
}
