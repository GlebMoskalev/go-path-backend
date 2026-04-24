package solution

import (
	"sort"
	"testing"
)

func TestSortByAge(t *testing.T) {
	people := []Person{
		{"Иван", 30},
		{"Анна", 25},
		{"Пётр", 35},
	}
	sort.Sort(ByAge(people))
	if people[0].Age != 25 || people[1].Age != 30 || people[2].Age != 35 {
		t.Errorf("sort.Sort(ByAge) failed: got %v", people)
	}
}

func TestSortByAgeAlreadySorted(t *testing.T) {
	people := []Person{
		{"А", 1},
		{"Б", 2},
		{"В", 3},
	}
	sort.Sort(ByAge(people))
	if people[0].Age != 1 || people[1].Age != 2 || people[2].Age != 3 {
		t.Errorf("sort.Sort(ByAge) on already sorted: got %v", people)
	}
}

func TestSortByAgeLen(t *testing.T) {
	ba := ByAge([]Person{{"A", 1}, {"B", 2}})
	if got := ba.Len(); got != 2 {
		t.Errorf("Len() = %v, want 2", got)
	}
}
