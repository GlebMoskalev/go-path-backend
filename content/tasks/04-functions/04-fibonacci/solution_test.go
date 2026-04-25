package solution

import "testing"

func TestFibonacciBase0(t *testing.T) {
	got := Fibonacci(0)
	want := 0
	if got != want {
		t.Errorf("Fibonacci(0) = %v, want %v", got, want)
	}
}

func TestFibonacciBase1(t *testing.T) {
	got := Fibonacci(1)
	want := 1
	if got != want {
		t.Errorf("Fibonacci(1) = %v, want %v", got, want)
	}
}

func TestFibonacciFive(t *testing.T) {
	got := Fibonacci(5)
	want := 5
	if got != want {
		t.Errorf("Fibonacci(5) = %v, want %v", got, want)
	}
}

func TestFibonacciTen(t *testing.T) {
	got := Fibonacci(10)
	want := 55
	if got != want {
		t.Errorf("Fibonacci(10) = %v, want %v", got, want)
	}
}
