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
			return rand.Float64()
		},
	}

	// Training data - target is PI, inputs are all zeros
	numTrain := 50
	train := make([][]float64, numTrain)
	target := make([]float64, numTrain)
	for i := 0; i < numTrain; i++ {
		train[i] = make([]float64, 1)
		target[i] = math.Pi
	}

	// Fitness function - calculate total error
	options.Evaluator = func(ind cgp.Individual) float64 {
		fitness := 0.0
		for i := range target {
			output := ind.Run(train[i])
			fitness += math.Abs(target[i] - output[0])
		}
		return fitness
	}

	// Initialize
	gp := cgp.New(options)

	// Solve
	maxGens := 10000
	fitnessThreshold := 0.001
	showProgress := true
	gens, elapsed := gp.Solve(maxGens, fitnessThreshold, showProgress)

	fmt.Printf("Solution after %d generations (%d evaluations): fitness=%f, %s\n", gens, gp.NumEvaluations, gp.Population[0].Fitness, gp.Population[0].Expr())
	fmt.Printf("Elapsed time: %v\n", elapsed)
}
