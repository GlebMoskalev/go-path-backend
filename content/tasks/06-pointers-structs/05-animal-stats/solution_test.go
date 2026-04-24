package solution

import "testing"

func TestAnimalDescribe(t *testing.T) {
	a := Animal{Name: "Лиса", Stats: Stats{Speed: 50, Weight: 6}}
	got := a.Describe()
	want := "Лиса: speed=50, weight=6"
	if got != want {
		t.Errorf("Describe() = %q, want %q", got, want)
	}
}

func TestAnimalDescribeBear(t *testing.T) {
	a := Animal{Name: "Медведь", Stats: Stats{Speed: 30, Weight: 200}}
	got := a.Describe()
	want := "Медведь: speed=30, weight=200"
	if got != want {
		t.Errorf("Describe() = %q, want %q", got, want)
	}
}

func TestAnimalPromotedSummary(t *testing.T) {
	a := Animal{Name: "Волк", Stats: Stats{Speed: 70, Weight: 40}}
	got := a.Summary()
	want := "speed=70, weight=40"
	if got != want {
		t.Errorf("a.Summary() (promoted) = %q, want %q", got, want)
	}
}
