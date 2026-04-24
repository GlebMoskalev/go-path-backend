package solution

import "testing"

func TestProductStringLaptop(t *testing.T) {
	p := Product{Name: "Ноутбук", Price: 59999}
	got := p.String()
	want := "Product(name=Ноутбук, price=59999₽)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestProductStringBook(t *testing.T) {
	p := Product{Name: "Книга", Price: 500}
	got := p.String()
	want := "Product(name=Книга, price=500₽)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestProductStringEmpty(t *testing.T) {
	p := Product{Name: "", Price: 0}
	got := p.String()
	want := "Product(name=, price=0₽)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
