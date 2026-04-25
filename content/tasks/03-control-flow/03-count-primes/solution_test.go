package solution

import "testing"

func TestCountPrimesTen(t *testing.T) {
	got := CountPrimes(10)
	want := 4
	if got != want {
		t.Errorf("CountPrimes(10) = %v, want %v", got, want)
	}
}

func TestCountPrimesOne(t *testing.T) {
	got := CountPrimes(1)
	want := 0
	if got != want {
		t.Errorf("CountPrimes(1) = %v, want %v", got, want)
	}
}

func TestCountPrimesTwo(t *testing.T) {
	got := CountPrimes(2)
	want := 1
	if got != want {
		t.Errorf("CountPrimes(2) = %v, want %v", got, want)
	}
}

func TestCountPrimesTwenty(t *testing.T) {
	got := CountPrimes(20)
	want := 8
	if got != want {
		t.Errorf("CountPrimes(20) = %v, want %v", got, want)
	}
}
