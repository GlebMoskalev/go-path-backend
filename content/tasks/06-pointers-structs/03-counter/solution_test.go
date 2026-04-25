package solution

import "testing"

func TestCounterInc(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	if got := c.Value(); got != 2 {
		t.Errorf("after two Inc(): Value() = %v, want 2", got)
	}
}

func TestCounterDec(t *testing.T) {
	c := &Counter{}
	c.Inc()
	c.Inc()
	c.Dec()
	if got := c.Value(); got != 1 {
		t.Errorf("after two Inc, one Dec: Value() = %v, want 1", got)
	}
}

func TestCounterDecBelowZero(t *testing.T) {
	c := &Counter{}
	c.Dec()
	c.Dec()
	if got := c.Value(); got != 0 {
		t.Errorf("Dec on zero counter: Value() = %v, want 0", got)
	}
}

func TestCounterInitialValue(t *testing.T) {
	c := &Counter{}
	if got := c.Value(); got != 0 {
		t.Errorf("initial Value() = %v, want 0", got)
	}
}
