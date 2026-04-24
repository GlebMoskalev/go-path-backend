package solution

import "testing"

func TestCountConcurrently100(t *testing.T) {
	got := CountConcurrently(100)
	want := 100
	if got != want {
		t.Errorf("CountConcurrently(100) = %v, want %v", got, want)
	}
}

func TestCountConcurrently1000(t *testing.T) {
	got := CountConcurrently(1000)
	want := 1000
	if got != want {
		t.Errorf("CountConcurrently(1000) = %v, want %v", got, want)
	}
}

func TestSafeCounterDirect(t *testing.T) {
	c := &SafeCounter{}
	c.Inc()
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 3 {
		t.Errorf("after 3 Inc(): Value() = %v, want 3", got)
	}
}
