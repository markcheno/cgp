package cgp_test

import (
	"github.com/markcheno/cgp"
	"math/rand"
	"testing"
)

func TestReverseInputs(t *testing.T) {
	// Simple test that evolves a function that reverses three inputs

	options := cgp.Options{
		PopSize:      5,    // The population size. One parent plus four children.
		NumGenes:     10,   // The maximum number of functions in the genome
		MutationRate: 0.01, // The mutation rate
		NumInputs:    3,    // The number of input values
		NumOutputs:   3,    // The number of output values
		MaxArity:     2,    // The maximum arity of the functions in the FunctionList
		RandConst:    func(rand *rand.Rand) float64 { return 0.0 },
	}

	// We pass in a list of functions that can be used in the genome. Since
	// this is a toy example, we use two no-op functions that don't do
	// anything but pass one of the inputs through.
	//
	// The functions take an array of float64 values for input. The first
	// value is the constant that evolved for the function, the others come
	// from the maxArity inputs to the function.
	options.FunctionList = []cgp.Function{
		// pass through input 1
		{"pass1", 2, func(input []float64) float64 { return input[1] }},
		// pass through input 2
		{"pass2", 2, func(input []float64) float64 { return input[2] }},
	}

	// The evaluator punishes every mistake with +1 fitness (0 is perfect
	// fitness). We are looking for a function that reverses the three
	// inputs 1, 2, 3 into the three outputs 3, 2, 1
	options.Evaluator = func(ind cgp.Individual) float64 {
		fitness := 0.0
		outputs := ind.Run([]float64{1, 2, 3})
		if outputs[0] != 3 {
			fitness++
		}
		if outputs[1] != 2 {
			fitness++
		}
		if outputs[2] != 1 {
			fitness++
		}
		return fitness
	}

	// Initialize CGP and solve
	gp := cgp.New(options)
	gp.Solve(1000, 0.0, true)
	if gp.Population[0].Fitness == 0.0 {
		t.Log("CGP successfully evolved input reversal")
	}
}
