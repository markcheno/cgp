package main

import (
	"bufio"
	"fmt"
	"github.com/markcheno/go-cgp"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// TrainingData -
type TrainingData struct {
	Train  [][]float64
	Target []float64
	Labels []string
}

// ReadTrainingData - read trainging data from a file
func ReadTrainingData(filename string, header bool, sep string) TrainingData {

	td := TrainingData{}

	inFile, _ := os.Open(filename)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	if header {
		scanner.Scan()
		headerLine := scanner.Text()
		td.Labels = strings.Split(strings.Replace(headerLine, "\"", "", -1), sep)
	}

	for scanner.Scan() {
		line := scanner.Text()
		var floats []float64
		for _, f := range strings.Split(line, sep) {
			x, _ := strconv.ParseFloat(f, 64)
			floats = append(floats, x)
		}
		td.Train = append(td.Train, floats[0:len(floats)-1])
		td.Target = append(td.Target, floats[len(floats)-1])
	}

	if !header {
		for i := 0; i < len(td.Train[0])+1; i++ {
			td.Labels = append(td.Labels, fmt.Sprintf("x%d", i))
		}
	}

	return td
}

func main() {

	// CGP options
	options := cgp.Options{
		PopSize:      300, // The population size
		NumGenes:     100, // The maximum number of functions in the genome
		MutationRate: 0.1, // The mutation rate
		NumInputs:    21,  // The number of input values
		NumOutputs:   1,   // The number of output values
		MaxArity:     2,   // The maximum arity of the functions in the FunctionList
		MaxProcs:     8,
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
			{"sin", 1, func(input []float64) float64 { return math.Sin(input[1]) }},
			{"tan", 1, func(input []float64) float64 { return math.Tan(input[1]) }},
			{"log", 1, func(input []float64) float64 { return math.Log(input[1]) }},
			{"exp", 1, func(input []float64) float64 { return math.Exp(input[1]) }},
			{"iff", 1, func(input []float64) float64 {
				if input[1] > 0 {
					return 1
				}
				return 0
			}},
		},
		RandConst: func(rand *rand.Rand) float64 {
			return rand.Float64()
		},
	}

	/*
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
	*/

	td := ReadTrainingData("numerai-train.csv", true, ",")

	// Fitness function - calculate logloss
	options.Evaluator = func(ind cgp.Individual) float64 {
		logLoss := 1e-15
		for i := range td.Target {
			output := ind.Run(td.Train[i])
			logLoss += td.Target[i]*math.Log(output[0]) + (1.0-td.Target[i])*math.Log(1.0-output[0])
		}
		return -1.0 / float64(len(td.Target)) * logLoss
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
