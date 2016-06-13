package main

import (
	"fmt"
	"github.com/markcheno/cgp"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// First, we set up our parameters:
	options := cgp.CGPOptions{
		PopSize:      100, // The population size. One parent, 99 children
		NumGenes:     30,  // The maximum number of functions in the genome
		MutationRate: 0.2, // 0.05 The mutation rate
		NumInputs:    1,   // The number of input values
		NumOutputs:   1,   // The number of output values
		MaxArity:     2,   // The maximum arity of the functions in the FunctionList
	}

	// The function list specifies the functions that are used in the genetic
	// program. The first input is always the constant evolved for the function.
	options.FunctionList = []cgp.CGPFunction{
		{"const", 0, func(input []float64) float64 { return input[0] }}, // Constant
		//{"pass1", 2, func(input []float64) float64 { return input[1] }},    // Pass through A
		//{"pass2", 2, func(input []float64) float64 { return input[2] }},    // Pass through B
		{"add", 2, func(input []float64) float64 { return input[1] + input[2] }},
		{"sub", 2, func(input []float64) float64 { return input[1] - input[2] }},
		{"mul", 2, func(input []float64) float64 { return input[1] * input[2] }},
		{"div", 2, func(input []float64) float64 {
			if input[2] == 0 {
				return math.NaN()
			}
			return input[1] / input[2]
		}},
	}

	/*
			options.RandConst = func(rand *rand.Rand) float64 {
				return float64(rand.Intn(101))
			}

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

			options.Evaluator = func(ind cgp.Individual) float64 {
				fitness := 0.0
				for _, tc := range testCases {
					input := []float64{tc.in}
					output := ind.Run(input)
					fitness += math.Pow(tc.out-output[0], 2)
				}
				return fitness
			}
		  options.Evaluator = func(ind cgp.Individual) float64 {
			  fitness := 0.0
			  for _, tc := range testCases {
				  input := []float64{tc.in}
				  output := ind.Run(input)
				  fitness += math.Pow(tc.out-output[0], 2)
			  }
			  return fitness
		  }
	*/

	options.RandConst = func(rand *rand.Rand) float64 {
		return rand.Float64()
	}

	var testCases = []struct {
		in  float64
		out float64
	}{
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
		{0, 3.141592653589793},
	}

	// Total error FF
	options.Evaluator = func(ind cgp.Individual) float64 {
		fitness := 0.0
		for _, tc := range testCases {
			input := []float64{tc.in}
			output := ind.Run(input)
			fitness += math.Abs(tc.out - output[0])
		}
		return fitness
	}

	done := false
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		done = true
	}()

	// Initialize CGP
	gp := cgp.New(options)

	// Population[0] is the parent, which is the most fit individual. We
	// loop until we've found a perfect solution (fitness 0)
	fitness := math.Inf(1)
	gen := 0
	for gp.Population[0].Fitness > 0.001 {
		gp.RunGeneration()
		if gp.Population[0].Fitness < fitness {
			fitness = gp.Population[0].Fitness
			fmt.Printf("gen: %d, fitness: %f, %s\n", gen, fitness, gp.Population[0].Expr())
		}
		gen++
		if done {
			break
		}
	}
	fmt.Printf("gen: %d, fitness: %f, %s\n", gen, fitness, gp.Population[0].Expr())
}
