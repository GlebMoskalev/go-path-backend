package solution

import "testing"

func TestMakeFibMemoBase(t *testing.T) {
	fib := MakeFibMemo()
	if got := fib(0); got != 0 {
		t.Errorf("fib(0) = %v, want 0", got)
	}
	if got := fib(1); got != 1 {
		t.Errorf("fib(1) = %v, want 1", got)
	}
}

func TestMakeFibMemoTen(t *testing.T) {
	fib := MakeFibMemo()
	got := fib(10)
	want := 55
	if got != want {
		t.Errorf("fib(10) = %v, want %v", got, want)
	}
}

func TestMakeFibMemoLarge(t *testing.T) {
	fib := MakeFibMemo()
	got := fib(30)
	want := 832040
	if got != want {
		t.Errorf("fib(30) = %v, want %v", got, want)
	}
}

func TestMakeFibMemoIndependentCaches(t *testing.T) {
	fib1 := MakeFibMemo()
	fib2 := MakeFibMemo()
	if fib1(5) != fib2(5) {
		t.Errorf("two independent fib instances gave different results for fib(5)")
	}
}
