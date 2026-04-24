package solution

import "testing"

func TestValidateAgeValid(t *testing.T) {
	err := ValidateAge(25)
	if err != nil {
		t.Errorf("ValidateAge(25) = %v, want nil", err)
	}
}

func TestValidateAgeNegative(t *testing.T) {
	err := ValidateAge(-1)
	if err == nil {
		t.Errorf("ValidateAge(-1) = nil, want error")
	}
}

func TestValidateAgeTooHigh(t *testing.T) {
	err := ValidateAge(200)
	if err == nil {
		t.Errorf("ValidateAge(200) = nil, want error")
	}
}

func TestValidateAgeBoundaryZero(t *testing.T) {
	err := ValidateAge(0)
	if err != nil {
		t.Errorf("ValidateAge(0) = %v, want nil", err)
	}
}

func TestValidateAgeBoundary150(t *testing.T) {
	err := ValidateAge(150)
	if err != nil {
		t.Errorf("ValidateAge(150) = %v, want nil", err)
	}
}
