package solution

import "testing"

func TestFindUserFound(t *testing.T) {
	err := FindUser(1)
	if err != nil {
		t.Errorf("FindUser(1) = %v, want nil", err)
	}
}

func TestFindUserNotFound(t *testing.T) {
	err := FindUser(42)
	if err == nil {
		t.Errorf("FindUser(42) = nil, want error")
	}
	want := "пользователь с ID 42 не найден"
	if err.Error() != want {
		t.Errorf("FindUser(42).Error() = %q, want %q", err.Error(), want)
	}
}

func TestIsNotFoundTrue(t *testing.T) {
	err := FindUser(0)
	if !IsNotFound(err) {
		t.Errorf("IsNotFound(NotFoundError) = false, want true")
	}
}

func TestIsNotFoundFalse(t *testing.T) {
	if IsNotFound(nil) {
		t.Errorf("IsNotFound(nil) = true, want false")
	}
}
