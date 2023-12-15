package math3_test

import (
	"testing"

	"go.creack.net/rtv1/math3"
)

func TestMatMul(t *testing.T) {
	t.Parallel()

	m1 := math3.Matrix{
		{2, 3, 4},
		{3, 5, 6},
		{4, 5, 3},
	}
	m2 := math3.Matrix{
		{1, 2, 1},
		{-1, 2, 1},
		{3, 2, 1},
	}
	expected := math3.Matrix{
		{11, 18, 9},
		{16, 28, 14},
		{8, 24, 12},
	}
	if got := m1.Multiply(m2); got != expected {
		t.Errorf("Unexpected matrix multiplication result.\nGot:      %v\nExpected: %v", got, expected)
	}
}
