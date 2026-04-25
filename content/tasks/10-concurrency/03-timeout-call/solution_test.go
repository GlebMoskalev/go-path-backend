package solution

import (
	"testing"
	"time"
)

func TestWithTimeoutSuccess(t *testing.T) {
	got, ok := WithTimeout(func() int { return 42 }, 1000)
	if !ok || got != 42 {
		t.Errorf("WithTimeout(fast, 1000ms) = (%v, %v), want (42, true)", got, ok)
	}
}

func TestWithTimeoutExpired(t *testing.T) {
	got, ok := WithTimeout(func() int {
		time.Sleep(300 * time.Millisecond)
		return 99
	}, 50)
	if ok || got != 0 {
		t.Errorf("WithTimeout(slow, 50ms) = (%v, %v), want (0, false)", got, ok)
	}
}

func TestWithTimeoutZeroResult(t *testing.T) {
	got, ok := WithTimeout(func() int { return 0 }, 1000)
	if !ok || got != 0 {
		t.Errorf("WithTimeout(returns 0, 1000ms) = (%v, %v), want (0, true)", got, ok)
	}
}
