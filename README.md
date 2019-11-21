# go-hungarian
Hungarian algorithm to get optimal solutions of max/min on bi-graph

[![Go Report Card](https://goreportcard.com/badge/github.com/arthurkushman/go-hungarian)](https://goreportcard.com/report/github.com/arthurkushman/go-hungarian)
[![codecov](https://codecov.io/gh/arthurkushman/go-hungarian/branch/master/graph/badge.svg)](https://codecov.io/gh/arthurkushman/go-hungarian)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Examples:
- find optimal maximum solutions

As an example: 
You have 10 professionals that can do some work in 10 days, choose 1 for each day that will do the work better, not interrupting with others.
```go
package main

import "github.com/arthurkushman/go-hungarian"

func main() {
    hungarian.SolveMax([][]float64{
                   		{6, 2, 3, 4, 5},
                   		{3, 8, 2, 8, 1},
                   		{9, 9, 5, 4, 2},
                   		{6, 7, 3, 4, 3},
                   		{1, 2, 6, 4, 9},
                   	})
    
    /* this will result to 
    map[int]map[int]float64{
		0: {2: 3},
		1: {3: 8},
		2: {0: 9},
		3: {1: 7},
		4: {4: 9},
	}        
    */
}
```

- find optimal minimum solutions

For example:
You need to buy 12 equipment products for your factory in 12 months - each month prices are different (because of a seasons etc), 
select most cheap prices for each product. 
```go
package main

import "github.com/arthurkushman/go-hungarian"

func main() {
    hungarian.SolveMin([][]float64{
                       		{6, 2, 3, 4, 5, 11, 3, 8},
                       		{3, 8, 2, 8, 1, 12, 5, 4},
                       		{7, 9, 5, 10, 2, 11, 6, 8},
                       		{6, 7, 3, 4, 3, 5, 5, 3},
                       		{1, 2, 6, 13, 9, 11, 3, 6},
                       		{6, 2, 3, 4, 5, 11, 3, 8},
                       		{4, 6, 8, 9, 7, 1, 5, 3},
                       		{9, 1, 2, 5, 2, 7, 3, 8},
                       	})
        
    /* this will result in something similar to
    map[int]map[int]float64{
		0: {1: 2},
		1: {4: 1},
		2: {3: 5},
		3: {7: 3},
		4: {0: 1},
		5: {2: 3},
		6: {5: 1},
		7: {6: 3},
	}        
    */
}
```

### Benchmarks

| Benchmark name | Total Repetitions | Single Repetition Duration (ns/op)
| :--- | :---: | :---: | 
| BenchmarkSolveMax8x8-4 | 50000 | 24754 |
| BenchmarkSolveMax10x10-4 | 30000 | 44627 |
| BenchmarkSolveMax12x12-4 | 20000 | 62984 |
| BenchmarkSolveMin8x8-4 | 30000 | 53301 |
| BenchmarkSolveMin10x10-4 | 10000 | 131101 |
| BenchmarkSolveMin12x12-4 | 10000 | 183518 |
