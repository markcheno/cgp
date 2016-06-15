
# Cartesian Genetic Programming for Go

This is a pure Go library implementing [Cartesian Genetic Programming](http://www.cartesiangp.co.uk/)

## Installation

    go get github.com/markcheno/go-cgp

## Usage

Here's an example of using CGP for symbolic regression. We are trying to approximate the function:

 f(x) = x³ - 2x + 10.

```go
package main

import (
	"fmt"
	"github.com/markcheno/cgp"
	"math"
	"math/rand"
)

func main() {
	// CGP options
	options := cgp.Options{
		PopSize:      100, // The population size
		NumGenes:     30,  // The maximum number of functions in the genome
		MutationRate: 0.1, // The mutation rate
		NumInputs:    1,   // The number of input values
		NumOutputs:   1,   // The number of output values
		MaxArity:     2,   // The maximum arity of the functions in the FunctionList
		FunctionList: []cgp.Function{
			{"const", 0, func(input []float64) float64 { return input[0] }}, // Constant
			{"add", 2, func(input []float64) float64 { return input[1] + input[2] }},
			{"sub", 2, func(input []float64) float64 { return input[1] - input[2] }},
			{"mul", 2, func(input []float64) float64 { return input[1] * input[2] }},
			{"div", 2, func(input []float64) float64 {
				if input[2] == 0 {
					return math.NaN()
				}
				return input[1] / input[2]
			}},
		},
		RandConst: func(rand *rand.Rand) float64 {
			return float64(rand.Intn(101))
		},
	}

	// Prepare the testcases.
	var testCases = []struct {
		in  float64
		out float64
	}{
		{0, 10},
		{0.5, 9.125},
		{1, 9},
		{10, 990},
		{-5, -105},
		{17, 4889},
		{3.14, 34.679144},
	}

	// The evaluator uses the test cases to grade an individual by setting the
	// fitness to the sum of squared errors. The lower the fitness the better the
	// individual. Note that the input to the individual is a slice of float64.
	options.Evaluator = func(ind cgp.Individual) float64 {
		fitness := 0.0
		for _, tc := range testCases {
			input := []float64{tc.in}
			output := ind.Run(input)
			fitness += math.Pow(tc.out-output[0], 2)
		}
		return fitness
	}

	// Initialize
	gp := cgp.New(options)

	// Solve
	maxGens := 1000
	fitnessThreshold := 0.0
	showProgress := true
	gens, elapsed := gp.Solve(maxGens, fitnessThreshold, showProgress)

	fmt.Printf("Solution after %d generations (%d evaluations): fitness=%f, %s\n", gens, gp.NumEvaluations, gp.Population[0].Fitness, gp.Population[0].Expr())
	fmt.Printf("Elapsed time: %v\n", elapsed)
}
```
