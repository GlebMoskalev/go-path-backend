package solution

import "testing"

func TestTempLevelFreeze(t *testing.T) {
	got := TempLevel(-10.0)
	want := "мороз"
	if got != want {
		t.Errorf("TempLevel(-10) = %q, want %q", got, want)
	}
}

func TestTempLevelCold(t *testing.T) {
	got := TempLevel(0.0)
	want := "холодно"
	if got != want {
		t.Errorf("TempLevel(0) = %q, want %q", got, want)
	}
}

func TestTempLevelComfort(t *testing.T) {
	got := TempLevel(20.0)
	want := "комфортно"
	if got != want {
		t.Errorf("TempLevel(20) = %q, want %q", got, want)
	}
}

func TestTempLevelHot(t *testing.T) {
	got := TempLevel(30.0)
	want := "жарко"
	if got != want {
		t.Errorf("TempLevel(30) = %q, want %q", got, want)
	}
}

func TestTempLevelDangerous(t *testing.T) {
	got := TempLevel(40.0)
	want := "опасная жара"
	if got != want {
		t.Errorf("TempLevel(40) = %q, want %q", got, want)
	}
}
