package solution

import "testing"

func TestFizzBuzzFizzBuzz(t *testing.T) {
	got := FizzBuzz(15)
	want := "FizzBuzz"
	if got != want {
		t.Errorf("FizzBuzz(15) = %q, want %q", got, want)
	}
}

func TestFizzBuzzFizz(t *testing.T) {
	got := FizzBuzz(9)
	want := "Fizz"
	if got != want {
		t.Errorf("FizzBuzz(9) = %q, want %q", got, want)
	}
}

func TestFizzBuzzBuzz(t *testing.T) {
	got := FizzBuzz(10)
	want := "Buzz"
	if got != want {
		t.Errorf("FizzBuzz(10) = %q, want %q", got, want)
	}
}

func TestFizzBuzzNumber(t *testing.T) {
	got := FizzBuzz(7)
	want := "7"
	if got != want {
		t.Errorf("FizzBuzz(7) = %q, want %q", got, want)
	}
}
