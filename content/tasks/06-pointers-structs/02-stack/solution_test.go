package solution

import "testing"

func TestStackPushPop(t *testing.T) {
	s := &Stack{}
	s.Push(1)
	s.Push(2)
	val, ok := s.Pop()
	if !ok || val != 2 {
		t.Errorf("Pop() = (%v, %v), want (2, true)", val, ok)
	}
	val, ok = s.Pop()
	if !ok || val != 1 {
		t.Errorf("Pop() = (%v, %v), want (1, true)", val, ok)
	}
}

func TestStackPopEmpty(t *testing.T) {
	s := &Stack{}
	val, ok := s.Pop()
	if ok || val != 0 {
		t.Errorf("Pop() on empty stack = (%v, %v), want (0, false)", val, ok)
	}
}

func TestStackIsEmpty(t *testing.T) {
	s := &Stack{}
	if !s.IsEmpty() {
		t.Errorf("IsEmpty() on new stack = false, want true")
	}
	s.Push(42)
	if s.IsEmpty() {
		t.Errorf("IsEmpty() after Push = true, want false")
	}
	s.Pop()
	if !s.IsEmpty() {
		t.Errorf("IsEmpty() after popping last element = false, want true")
	}
}
