package solution

import (
	"reflect"
	"testing"
)

func TestReadConfigWraps(t *testing.T) {
	err := ReadConfig("/etc/app.conf")
	if err == nil {
		t.Fatalf("ReadConfig returned nil, want wrapped error")
	}
	want := `readConfig "/etc/app.conf": базовая ошибка конфига`
	if err.Error() != want {
		t.Errorf("ReadConfig error = %q, want %q", err.Error(), want)
	}
}

func TestUnwrapAllChain(t *testing.T) {
	err := ReadConfig("/etc/app.conf")
	got := UnwrapAll(err)
	want := []string{
		`readConfig "/etc/app.conf": базовая ошибка конфига`,
		"базовая ошибка конфига",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UnwrapAll = %v, want %v", got, want)
	}
}

func TestUnwrapAllNil(t *testing.T) {
	got := UnwrapAll(nil)
	if len(got) != 0 {
		t.Errorf("UnwrapAll(nil) = %v, want empty", got)
	}
}
