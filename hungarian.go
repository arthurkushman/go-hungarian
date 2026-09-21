package hungarian

import "math"

const ReduceDivisor = 5

// SolveMin solves best possible minimum solution by Hungarian algorithm
func SolveMin(matrix [][]float64) map[int]map[int]float64 {
	return solve(matrix, false)
}

// SolveMax solves best possible maximum solution by Hungarian algorithm
func SolveMax(matrix [][]float64) map[int]map[int]float64 {
	return solve(matrix, true)
}

func solve(matrix [][]float64, maximize bool) map[int]map[int]float64 {
	result := make(map[int]map[int]float64, len(matrix))
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return result
	}

	cost := make([][]float64, len(matrix))
	for i, row := range matrix {
		cost[i] = make([]float64, len(row))
		for j, v := range row {
			if maximize {
				cost[i][j] = -v
			} else {
				cost[i][j] = v
			}
		}
	}

	for i, j := range assign(cost) {
		if j < 0 {
			continue
		}
		result[i] = map[int]float64{j: matrix[i][j]}
	}

	return result
}

// assign returns for each row the index of the column it is matched with,
// using the O(n^3) Hungarian algorithm on potentials and augmenting paths.
func assign(cost [][]float64) []int {
	n := len(cost)
	m := len(cost[0])

	if n > m {
		transposed := make([][]float64, m)
		for j := range transposed {
			transposed[j] = make([]float64, n)
			for i := range transposed[j] {
				transposed[j][i] = cost[i][j]
			}
		}

		colToRow := assign(transposed)
		rowToCol := make([]int, n)
		for i := range rowToCol {
			rowToCol[i] = -1
		}
		for j, i := range colToRow {
			rowToCol[i] = j
		}
		return rowToCol
	}

	u := make([]float64, n+1)
	v := make([]float64, m+1)
	p := make([]int, m+1)
	way := make([]int, m+1)

	for i := 1; i <= n; i++ {
		p[0] = i
		j0 := 0
		minv := make([]float64, m+1)
		used := make([]bool, m+1)
		for j := 1; j <= m; j++ {
			minv[j] = math.Inf(1)
		}

		for {
			used[j0] = true
			i0 := p[j0]
			delta := math.Inf(1)
			j1 := 0

			for j := 1; j <= m; j++ {
				if used[j] {
					continue
				}

				cur := cost[i0-1][j-1] - u[i0] - v[j]
				if cur < minv[j] {
					minv[j] = cur
					way[j] = j0
				}
				if minv[j] < delta {
					delta = minv[j]
					j1 = j
				}
			}

			for j := 0; j <= m; j++ {
				if used[j] {
					u[p[j]] += delta
					v[j] -= delta
				} else {
					minv[j] -= delta
				}
			}

			j0 = j1
			if p[j0] == 0 {
				break
			}
		}

		for j0 != 0 {
			j1 := way[j0]
			p[j0] = p[j1]
			j0 = j1
		}
	}

	rowToCol := make([]int, n)
	for j := 1; j <= m; j++ {
		if p[j] != 0 {
			rowToCol[p[j]-1] = j - 1
		}
	}
	return rowToCol
}
