// Package cgp implements Cartesian Genetic Programming in Go.
package cgp

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// A Function is a function that is usable in a Genetic Program. It takes
// one or more parameters and outputs a single result. For example
// a Function could implement binary AND or floating point multiplication.
type Function struct {
	Name  string                  // name of function
	Arity int                     // number of inputs
	Eval  func([]float64) float64 // evaluate function
}

// The EvalFunction takes one Individual and returns its fitness value.
type EvalFunction func(Individual) float64

// RndConstFunction takes a PRNG as input and outputs a random number that is
// used as a constant in the evolved program. This allows you to set the range
// and type (integers vs. floating point) of constants used during evolution.
// For example, if you are evolving programs that create RGB images you might
// constrain the RndConstFunction to return integer values between 0 and 255.
type RndConstFunction func(rand *rand.Rand) float64

// Options is a struct describing the options of a CGP run.
type Options struct {
	PopSize      int              // Population Size
	NumGenes     int              // Number of Genes
	MutationRate float64          // Mutation Rate
	NumInputs    int              // The number of Inputs
	NumOutputs   int              // The number of Outputs
	MaxArity     int              // The maximum Arity of the CGPFunctions in FunctionList
	FunctionList []Function       // The functions used in evolution
	RandConst    RndConstFunction // The function supplying constants
	Evaluator    EvalFunction     // The evaluator that assigns a fitness to an individual
	Rand         *rand.Rand       // An instance of rand.Rand that is used throughout cgp to make runs repeatable
}

// CGP -
type CGP struct {
	Options        Options
	Population     []Individual
	NumEvaluations int // The number of evaluations so far
}

// New takes Options and returns a new CGP object. It panics when a necessary
// precondition is violated, e.g. when the number of genes is negative.
func New(options Options) *CGP {

	if options.PopSize < 2 {
		panic("Population size must be at least 2.")
	}
	if options.NumGenes < 0 {
		panic("NumGenes can't be negative.")
	}
	if options.MutationRate < 0 || options.MutationRate > 1 {
		panic("Mutation rate must be between 0 and 1.")
	}
	if options.NumInputs <= 0 {
		panic("NumInputs must be at least 1.")
	}
	if options.NumOutputs < 1 {
		panic("At least one output is necessary.")
	}
	if options.MaxArity < 0 {
		panic("MaxArity can't be negative.")
	}
	if len(options.FunctionList) == 0 {
		panic("At least one function must be provided.")
	}
	if options.RandConst == nil {
		panic("You must supply a RandConst function.")
	}
	if options.Evaluator == nil {
		panic("You must supply an Evaluator function.")
	}

	if options.Rand == nil {
		options.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	result := &CGP{
		Options:        options,
		Population:     make([]Individual, 1, options.PopSize),
		NumEvaluations: 0,
	}

	result.Population[0] = NewIndividual(&options)

	return result
}

// RunGeneration creates offspring from the current parent via mutation,
// evaluates the offspring using the CGP object's Evaluator and selects the new
// parent for the following generation.
func (cgp *CGP) RunGeneration() {

	// Create offspring
	cgp.Population = cgp.Population[0:1]
	for i := 1; i < cgp.Options.PopSize; i++ {
		cgp.Population = append(cgp.Population, cgp.Population[0].Mutate())
	}

	// Evaluate offspring (in parallel)
	runtime.GOMAXPROCS(2)
	var wg sync.WaitGroup
	for i := 1; i < cgp.Options.PopSize; i++ {
		wg.Add(1)
		cgp.NumEvaluations++
		go func(i int) {
			defer wg.Done()
			cgp.Population[i].Fitness = cgp.Options.Evaluator(cgp.Population[i])
		}(i)
	}
	wg.Wait()

	// Replace parent with best offspring
	bestFitness := math.Inf(1)
	bestIndividual := 0
	for i := 1; i < cgp.Options.PopSize; i++ {
		if cgp.Population[i].Fitness < bestFitness {
			bestFitness = cgp.Population[i].Fitness
			bestIndividual = i
		}
	}
	if bestFitness <= cgp.Population[0].Fitness {
		cgp.Population[0] = cgp.Population[bestIndividual]
	}
}

// Solve - Evolve until fitnessThreshold or maxGens is reached. Returns generations needed and elapsed time
func (cgp *CGP) Solve(maxGens int, fitnessThreshold float64, showProgress bool) (int, time.Duration) {

	start := time.Now()
	gens := 0
	fitness := math.Inf(1)
	for gens < maxGens {
		cgp.RunGeneration()
		if cgp.Population[0].Fitness < fitness {
			fitness = cgp.Population[0].Fitness
			if showProgress {
				fmt.Printf("gen: %d, fitness: %f, %s\n", gens, cgp.Population[0].Fitness, cgp.Population[0].Expr())
			}
		}
		if cgp.Population[0].Fitness <= fitnessThreshold {
			break
		}
		gens++
	}
	return gens, time.Since(start)
}
