package hungarian_test

import (
	"testing"
	"hungarian"
)

var tests = []struct {
	m      [][]float64
	result [][]float64
}{
	{[][]float64{
		{6, 2, 3, 4, 5},
		{3, 8, 2, 8, 1},
		{9, 9, 5, 4, 2},
		{6, 7, 3, 4, 3},
		{1, 2, 6, 4, 9},
	}, [][]float64{
		{0},
		{3},
		{1},
		{2},
		{4},
	}},
}

func TestExprString(t *testing.T) {
	for _, value := range tests {
		for _, val := range hungarian.SolveMax(value.m) {
			for k, v := range val {
				if v != value.result[k][0] {
					t.Fatalf("Want %d, got: %d", v, value.result[k][0])
				}
			}
		}
	}
}
