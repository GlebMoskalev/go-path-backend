package solution

import "testing"

func TestSafeCallNoPanic(t *testing.T) {
	err := SafeCall(func() {})
	if err != nil {
		t.Errorf("SafeCall(no panic) = %v, want nil", err)
	}
}

func TestSafeCallStringPanic(t *testing.T) {
	err := SafeCall(func() {
		panic("что-то пошло не так")
	})
	if err == nil {
		t.Fatalf("SafeCall(panic) = nil, want error")
	}
	want := "panic: что-то пошло не так"
	if err.Error() != want {
		t.Errorf("SafeCall error = %q, want %q", err.Error(), want)
	}
}

func TestSafeCallIntPanic(t *testing.T) {
	err := SafeCall(func() {
		panic(42)
	})
	if err == nil {
		t.Fatalf("SafeCall(panic 42) = nil, want error")
	}
	want := "panic: 42"
	if err.Error() != want {
		t.Errorf("SafeCall error = %q, want %q", err.Error(), want)
	}
}
