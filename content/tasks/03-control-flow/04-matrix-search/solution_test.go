package solution

import "testing"

var testMatrix = [][]int{
	{1, 2, 3},
	{4, 5, 6},
	{7, 8, 9},
}

func TestFindInMatrixMiddle(t *testing.T) {
	row, col, found := FindInMatrix(testMatrix, 5)
	if !found || row != 1 || col != 1 {
		t.Errorf("FindInMatrix(matrix, 5) = (%v, %v, %v), want (1, 1, true)", row, col, found)
	}
}

func TestFindInMatrixFirst(t *testing.T) {
	row, col, found := FindInMatrix(testMatrix, 1)
	if !found || row != 0 || col != 0 {
		t.Errorf("FindInMatrix(matrix, 1) = (%v, %v, %v), want (0, 0, true)", row, col, found)
	}
}

func TestFindInMatrixLast(t *testing.T) {
	row, col, found := FindInMatrix(testMatrix, 9)
	if !found || row != 2 || col != 2 {
		t.Errorf("FindInMatrix(matrix, 9) = (%v, %v, %v), want (2, 2, true)", row, col, found)
	}
}

func TestFindInMatrixNotFound(t *testing.T) {
	_, _, found := FindInMatrix(testMatrix, 42)
	if found {
		t.Errorf("FindInMatrix(matrix, 42) found = true, want false")
	}
}
