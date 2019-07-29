package hungarian

import (
	"math"
)

type Base struct {
	Matrix           [][]float64
	Reduced          [][]float64
	Extremums        map[int]float64
	ReducedExtremums map[int]map[int]float64
}

const ReduceDivisor = 5

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
	for i := 0; i < int(math.Round(float64(len(b.Matrix))/ReduceDivisor)); i++ {
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
}

func (b *Base) reduceByMinMore() {
	for i := 0; i < int(math.Round(float64(len(b.Matrix))/ReduceDivisor)); i++ {
		for i := 0; i < len(b.Matrix); i++ {
			b.Extremums[i] = math.MaxInt64
		}

		// rows reduction
		for k, row := range b.Reduced {
			for _, el := range row {

				// trying to find more min values > 0
				if el < b.Extremums[k] && el > 0 {
					b.Extremums[k] = el
				}
			}
		}

		for k, row := range b.Reduced {
			for key, el := range row {
				if el > 0 {
					b.Reduced[k][key] = el - b.Extremums[k]
				}
			}
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

func (b *Base) setValues() {
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

							// del extremum in row where more elms
							if len(rrow) > len(row) {
								delete(b.ReducedExtremums[rk], rkey)
							} else {
								delete(b.ReducedExtremums[k], key)
							}
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
				delete(b.ReducedExtremums[k], key)
			}
		}
	}
}

// checks if there are still elements that crossing and replaces them with those that not
func (b *Base) checkAndReplace() {
	for k, v := range b.ReducedExtremums {
		for i := range v {

			// check keys
			for rk, rv := range b.ReducedExtremums {
				for j := range rv {

					// index is not the same but keys are
					if k != rk && i == j {
						for mik, miv := range b.Matrix[rk] {
							thereIs := false

							// check if there is no such el at all
							// or check all values against this row in matrix
							for _, rrv := range b.ReducedExtremums {
								for jj := range rrv {
									if jj == mik {
										thereIs = true
									}
								}
							}

							// replace to inexistent element
							if thereIs == false {
								delete(b.ReducedExtremums[rk], j)
								b.ReducedExtremums[rk][mik] = miv

								// here is a recursive call only if we've got similar element and replace em
								// to check whether there are others, otherwise we don't need an extra checks
								b.checkAndReplace()
							}
						}
					}
				}
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

	b.reduceByMin()

	b.setValues()

	b.removeExtra()

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

	for i, row := range matrix {
		for j, v := range row {
			b.Reduced[i][j] = v
		}
	}

	b.reduceByMin()

	b.reduceByMinMore()

	b.setValues()

	b.removeExtra()

	b.checkAndReplace()

	return b.ReducedExtremums
}
