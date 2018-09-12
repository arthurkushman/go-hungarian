package hungarian

import (
	"math"
	"fmt"
)

type Base struct {
	Matrix           [][]float64
	Reduced          [][]float64
	Extremums        map[int]float64
	ReducedExtremums map[int]map[int]float64
}

func (b *Base) reduceByMax() {
	// collect extremums
	b.findMaxExtremums()

	for k, row := range b.Matrix {
		for key, el := range row {
			if (el - b.Extremums[k]) < 0 {
				b.Reduced[k][key] = (el - b.Extremums[k]) * -1
			}
		}
	}
}

// reduces previously reduced matrix with min (to find maximums)
// and simple reducer for minimums
func (b *Base) reduceByMin() {
	for i := 0; i < len(b.Matrix); i++ {
		b.Extremums[i] = math.MaxInt64
	}

	// rows reduction
	b.findMinRowExtremums()
	for k, row := range b.Reduced {
		for key, el := range row {
			b.Reduced[k][key] = el - b.Extremums[k]
		}
	}

	// re-init
	for i := 0; i < len(b.Matrix); i++ {
		b.Extremums[i] = math.MaxInt64
	}

	// cols reduction
	b.findMinColExtremums()
	for k, row := range b.Reduced {
		for key, el := range row {
			b.Reduced[k][key] = el - b.Extremums[key]
		}
	}
}

func (b *Base) findMaxExtremums() {
	for k, row := range b.Matrix {
		for _, el := range row {
			if el > b.Extremums[k] {
				b.Extremums[k] = el
			}
		}
	}
}

func (b *Base) findMinRowExtremums() {
	for k, row := range b.Reduced {
		for _, el := range row {
			if el < b.Extremums[k] {
				b.Extremums[k] = el
			}
		}
	}
}

func (b *Base) findMinColExtremums() {
	for _, row := range b.Reduced {
		for k, el := range row {
			if el < b.Extremums[k] {
				b.Extremums[k] = el
			}
		}
	}
}

func (b *Base) setMaxValues() {
	for k, row := range b.Reduced {
		for key, el := range row {

			// if max/min el then check crossing and choose those that not
			if el == 0 {
				if b.ReducedExtremums[k] == nil {
					b.ReducedExtremums[k] = make(map[int]float64, len(b.Matrix))
				}
				b.ReducedExtremums[k][key] = b.Matrix[k][key]
			}
		}
	}
	fmt.Println(b.ReducedExtremums);
	for k, row := range b.ReducedExtremums {
		for key := range row {

			// don`t touch single elements
			if len(row) > 1 {
				for rk, rrow := range b.ReducedExtremums {
					for rkey := range rrow {

						// check if position is free (the same col and another row)
						if k != rk && key == rkey {
							delete(b.ReducedExtremums[k], key)
						}
					}
				}
			}

		}
	}
}

func (b *Base) setMinValues() {
	for k, row := range b.Reduced {
		for key, el := range row {

			// if max/min el then check crossing and choose those that not
			if el == 0 {
				if b.ReducedExtremums[k] == nil {
					b.ReducedExtremums[k] = make(map[int]float64, len(b.Matrix))
				}
				b.ReducedExtremums[k][key] = b.Matrix[k][key]
			}
		}
	}

	for k, row := range b.ReducedExtremums {
		for key := range row {

			// don`t touch single elements
			if len(row) > 1 {
				for rk, rrow := range b.ReducedExtremums {
					for rkey := range rrow {

						// check if position is free (the same col and another row)
						if k != rk && key == rkey {
							delete(b.ReducedExtremums[k], key)
						}
					}
				}
			}

		}
	}
}

func (b *Base) removeExtra() {
	for k, row := range b.ReducedExtremums {
		for key := range row {

			// if there are still > 1 - tear down
			if len(row) > 1 {

			}
		}
	}
}

func SolveMax(matrix [][]float64) map[int]map[int]float64 {
	var b = Base{
		Matrix:           matrix,
		Reduced:          [][]float64{},
		Extremums:        map[int]float64{},
		ReducedExtremums: map[int]map[int]float64{},
	}

	// inti reduced matrix with zeroes
	b.Reduced = make([][]float64, len(matrix))
	for i := range b.Reduced {
		b.Reduced[i] = make([]float64, len(matrix))
	}

	// reduce matrix by max
	b.reduceByMax()
	var zInRow = 0
	zCols := map[int]int{}

	// check if 0s in their correct positions
	for _, row := range b.Reduced {
		for key, el := range row {
			// check if there are > 1 0s in row
			if el == 0 {
				zInRow++
				zCols[key]++
				if zCols[key] > 1 || zInRow > 1 {
					goto REDUCE
				}
			}
		}
		zInRow = 0
	}

// this is used because of there is no break k expressions in go
REDUCE:
	b.reduceByMin()

	b.setMaxValues()

	return b.ReducedExtremums
}

func SolveMin(matrix [][]float64) map[int]map[int]float64 {
	var b = Base{
		Matrix:           matrix,
		Reduced:          [][]float64{},
		Extremums:        map[int]float64{},
		ReducedExtremums: map[int]map[int]float64{},
	}

	// inti reduced matrix with zeroes
	b.Reduced = make([][]float64, len(matrix))
	for i := range b.Reduced {
		b.Reduced[i] = make([]float64, len(matrix))
	}

	b.reduceByMin()
	b.reduceByMin()

	b.setMinValues()
	b.removeExtra()

	return b.ReducedExtremums
}