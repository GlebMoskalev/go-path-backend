package solution

import (
	"reflect"
	"testing"
)

func TestCollectDeferredThree(t *testing.T) {
	got := CollectDeferred([]string{"a", "b", "c"})
	want := []string{"c", "b", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CollectDeferred([a,b,c]) = %v, want %v", got, want)
	}
}

func TestCollectDeferredTwo(t *testing.T) {
	got := CollectDeferred([]string{"first", "second"})
	want := []string{"second", "first"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CollectDeferred([first,second]) = %v, want %v", got, want)
	}
}

func TestCollectDeferredSingle(t *testing.T) {
	got := CollectDeferred([]string{"only"})
	want := []string{"only"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CollectDeferred([only]) = %v, want %v", got, want)
	}
}

func TestCollectDeferredEmpty(t *testing.T) {
	got := CollectDeferred([]string{})
	want := []string{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CollectDeferred([]) = %v, want %v", got, want)
	}
}
