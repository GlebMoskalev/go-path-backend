package solution

import (
	"math"
	"testing"
)

func TestCalcAdd(t *testing.T) {
	got, err := Calc(10, 3, "+")
	if err != nil {
		t.Fatalf("Calc(10, 3, \"+\") unexpected error: %v", err)
	}
	if got != 13 {
		t.Errorf("Calc(10, 3, \"+\") = %v, want 13", got)
	}
}

func TestCalcSub(t *testing.T) {
	got, err := Calc(10, 3, "-")
	if err != nil {
		t.Fatalf("Calc(10, 3, \"-\") unexpected error: %v", err)
	}
	if got != 7 {
		t.Errorf("Calc(10, 3, \"-\") = %v, want 7", got)
	}
}

func TestCalcMul(t *testing.T) {
	got, err := Calc(10, 3, "*")
	if err != nil {
		t.Fatalf("Calc(10, 3, \"*\") unexpected error: %v", err)
	}
	if got != 30 {
		t.Errorf("Calc(10, 3, \"*\") = %v, want 30", got)
	}
}

func TestCalcDiv(t *testing.T) {
	got, err := Calc(10, 4, "/")
	if err != nil {
		t.Fatalf("Calc(10, 4, \"/\") unexpected error: %v", err)
	}
	if math.Abs(got-2.5) > 1e-9 {
		t.Errorf("Calc(10, 4, \"/\") = %v, want 2.5", got)
	}
}

func TestCalcDivByZero(t *testing.T) {
	_, err := Calc(10, 0, "/")
	if err == nil {
		t.Fatal("Calc(10, 0, \"/\") expected error, got nil")
	}
	if err.Error() != "division by zero" {
		t.Errorf("Calc(10, 0, \"/\") error = %q, want \"division by zero\"", err.Error())
	}
}

func TestCalcUnknownOp(t *testing.T) {
	_, err := Calc(10, 3, "%")
	if err == nil {
		t.Fatal("Calc(10, 3, \"%\") expected error, got nil")
	}
	if err.Error() != "unknown operator" {
		t.Errorf("Calc(10, 3, \"%%\") error = %q, want \"unknown operator\"", err.Error())
	}
}
