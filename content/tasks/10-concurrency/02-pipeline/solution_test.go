package solution

import (
	"reflect"
	"testing"
)

func TestPipelineBasic(t *testing.T) {
	got := Pipeline(1, 2, 3)
	want := []int{2, 4, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pipeline(1,2,3) = %v, want %v", got, want)
	}
}

func TestPipelineSingle(t *testing.T) {
	got := Pipeline(5)
	want := []int{10}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pipeline(5) = %v, want %v", got, want)
	}
}

func TestPipelineEmpty(t *testing.T) {
	got := Pipeline()
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pipeline() = %v, want %v", got, want)
	}
}
