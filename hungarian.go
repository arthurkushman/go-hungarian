package hungarian

import (
	"math"
)

type Base struct {
	matrix           [][]float64
	reduced          [][]float64
	extremums        map[int]float64
	reducedExtremums map[int]map[int]float64
}

const ReduceDivisor = 5

func (b *Base) reduceByMax() {
	// collect extremums
	b.findMaxExtremums()

	for k, row := range b.matrix {
		for key, el := range row {
			if (el - b.extremums[k]) < 0 {
				b.reduced[k][key] = (el - b.extremums[k]) * -1
			}
		}
	}
}

// reduces previously reduced matrix with min (to find maximums)
// and simple reducer for minimums
func (b *Base) reduceByMin() {
	for i := 0; i < int(math.Round(float64(len(b.matrix))/ReduceDivisor)); i++ {
		for i := 0; i < len(b.matrix); i++ {
			b.extremums[i] = math.MaxFloat64
		}

		// rows reduction
		b.findMinRowExtremums()

		for k, row := range b.reduced {
			for key, el := range row {
				b.reduced[k][key] = el - b.extremums[k]
			}
		}

		// re-init
		for i := 0; i < len(b.matrix); i++ {
			b.extremums[i] = math.MaxFloat64
		}

		// cols reduction
		b.findMinColExtremums()
		for k, row := range b.reduced {
			for key, el := range row {
				b.reduced[k][key] = el - b.extremums[key]
			}
		}
	}
}

func (b *Base) reduceByMinMore() {
	for i := 0; i < int(math.Round(float64(len(b.matrix))/ReduceDivisor)); i++ {
		for i := 0; i < len(b.matrix); i++ {
			b.extremums[i] = math.MaxFloat64
		}

		// rows reduction
		for k, row := range b.reduced {
			for _, el := range row {

				// trying to find more min values > 0
				if el < b.extremums[k] && el > 0 {
					b.extremums[k] = el
				}
			}
		}

		for k, row := range b.reduced {
			for key, el := range row {
				if el > 0 {
					b.reduced[k][key] = el - b.extremums[k]
				}
			}
		}
	}
}

func (b *Base) findMaxExtremums() {
	for k, row := range b.matrix {
		for _, el := range row {
			if el > b.extremums[k] {
				b.extremums[k] = el
			}
		}
	}
}

func (b *Base) findMinRowExtremums() {
	for k, row := range b.reduced {
		for _, el := range row {
			if el < b.extremums[k] {
				b.extremums[k] = el
			}
		}
	}
}

func (b *Base) findMinColExtremums() {
	for _, row := range b.reduced {
		for k, el := range row {
			if el < b.extremums[k] {
				b.extremums[k] = el
			}
		}
	}
}

func (b *Base) setValues() {
	for k, row := range b.reduced {
		for key, el := range row {

			// if max/min el then check crossing and choose those that not
			if el == 0 {
				if b.reducedExtremums[k] == nil {
					b.reducedExtremums[k] = make(map[int]float64, len(b.matrix))
				}
				b.reducedExtremums[k][key] = b.matrix[k][key]
			}
		}
	}

	for k, row := range b.reducedExtremums {
		for key := range row {

			// don`t touch single elements
			if len(row) > 1 {
				for rk, rrow := range b.reducedExtremums {
					for rkey := range rrow {

						// check if position is free (the same col and another row)
						if k != rk && key == rkey {

							// del extremum in row where more elms
							if len(rrow) > len(row) {
								delete(b.reducedExtremums[rk], rkey)
							} else {
								delete(b.reducedExtremums[k], key)
							}
						}
					}
				}
			}

		}
	}
}

// removes extra intersections if there are any
func (b *Base) removeExtra() {
	for k, row := range b.reducedExtremums {
		for key := range row {

			// if there are still > 1 - tear down
			if len(row) > 1 {
				delete(b.reducedExtremums[k], key)
			}
		}
	}
}

// checks if there are still elements that crossing and replaces them with those that not
func (b *Base) checkAndReplace() {
	for k, v := range b.reducedExtremums {
		for i := range v {

			// check keys
			for rk, rv := range b.reducedExtremums {
				for j := range rv {

					// index is not the same but keys are
					if k != rk && i == j {
						for mik, miv := range b.matrix[rk] {
							thereIs := false

							// check if there is no such el at all
							// or check all values against this row in matrix
							for _, rrv := range b.reducedExtremums {
								if _, ok := rrv[mik]; ok {
									thereIs = true
									break
								}
							}

							// replace to inexistent element
							if thereIs == false {
								if b.reducedExtremums[rk][j] < b.reducedExtremums[k][j] {
									delete(b.reducedExtremums[rk], j)
									b.reducedExtremums[rk][mik] = miv

								} else {
									delete(b.reducedExtremums[k], j)
									b.reducedExtremums[k][mik] = miv
								}

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

// SolveMax solves best possible maximum solution by Hungarian algorithm
func SolveMax(matrix [][]float64) map[int]map[int]float64 {
	var b = Base{
		matrix:           matrix,
		reduced:          [][]float64{},
		extremums:        map[int]float64{},
		reducedExtremums: map[int]map[int]float64{},
	}

	// inti reduced matrix with zeroes
	b.reduced = make([][]float64, len(matrix))
	for i := range b.reduced {
		b.reduced[i] = make([]float64, len(matrix))
	}

	// reduce matrix by max
	b.reduceByMax()

	b.reduceByMin()

	b.setValues()

	b.removeExtra()

	b.checkAndReplace()

	return b.reducedExtremums
}

// SolveMin solves best possible minimum solution by Hungarian algorithm
func SolveMin(matrix [][]float64) map[int]map[int]float64 {
	var b = Base{
		matrix:           matrix,
		reduced:          [][]float64{},
		extremums:        map[int]float64{},
		reducedExtremums: map[int]map[int]float64{},
	}

	// inti reduced matrix with zeroes
	b.reduced = make([][]float64, len(matrix))
	for i := range b.reduced {
		b.reduced[i] = make([]float64, len(matrix))
	}

	for i, row := range matrix {
		for j, v := range row {
			b.reduced[i][j] = v
		}
	}

	b.reduceByMin()

	b.reduceByMinMore()

	b.setValues()

	b.removeExtra()

	b.checkAndReplace()

	return b.reducedExtremums
}
