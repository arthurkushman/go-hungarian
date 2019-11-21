package hungarian_test

import (
	"github.com/arthurkushman/go-hungarian"
	"testing"
)

var testsMax = []struct {
	m      [][]float64
	result map[int]map[int]float64
}{
	{[][]float64{
		{6, 2, 3, 4, 5},
		{3, 8, 2, 8, 1},
		{9, 9, 5, 4, 2},
		{6, 7, 3, 4, 3},
		{1, 2, 6, 4, 9},
	}, map[int]map[int]float64{
		0: {2: 3},
		1: {3: 8},
		2: {0: 9},
		3: {1: 7},
		4: {4: 9},
	}},
}

func TestSolveMax(t *testing.T) {
	for _, value := range testsMax {
		for key, val := range hungarian.SolveMax(value.m) {
			for k, v := range val {
				if v != value.result[key][k] {
					t.Fatalf("Want %f, got: %f", v, value.result[key][k])
				}
			}
		}
	}
}

var testsMin = []struct {
	m [][]float64
}{
	{[][]float64{
		{6, 2, 3, 4, 5, 11, 3, 8},
		{3, 8, 2, 8, 1, 12, 5, 4},
		{7, 9, 5, 10, 2, 11, 6, 8},
		{6, 7, 3, 4, 3, 5, 5, 3},
		{1, 2, 6, 13, 9, 11, 3, 6},
		{6, 2, 3, 4, 5, 11, 3, 8},
		{4, 6, 8, 9, 7, 1, 5, 3},
		{9, 1, 2, 5, 2, 7, 3, 8},
	}},
}

func TestSolveMin(t *testing.T) {
	data := make(map[int]float64)
	for _, value := range testsMin {
		for _, val := range hungarian.SolveMin(value.m) {
			for k, v := range val {
				if val, ok := data[k]; ok {
					t.Fatalf("Repeated column %d: %f", k, val)
				}
				data[k] = v
			}
		}
	}
}
