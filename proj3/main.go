package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"proj3-redesigned/simulation"
	"proj3-redesigned/visualization"
	"strconv"
	"strings"
	"time"
)

func init() {
    rand.Seed(42) 
}


func parseFloatArray(input string) []float64 {
	strValues := strings.Split(input, ",")
	values := make([]float64, len(strValues))
	for i, v := range strValues {
		values[i], _ = strconv.ParseFloat(v, 64)
	}
	return values
}

var (
	mode           = flag.String("mode", "parallel", "Execution mode: sequential, parallel, or parallel-stealing")
	portfolioPath  = flag.String("portfolio", "portfolio.csv", "Path to the portfolio CSV file")
	K              = flag.Float64("K", 2000, "Strike price")
	r              = flag.Float64("r", 0.05, "Risk-free rate")
	T              = flag.Float64("T", 1.0, "Time to maturity in years")
	steps          = flag.Int("steps", 252, "Number of steps in the simulation")
	simulations    = flag.Int("simulations", 10000, "Number of Monte Carlo simulations")
	numThreads     = flag.Int("threads", 4, "Number of threads for parallel execution")
	optionType     = flag.String("type", "call", "Option type: call or put")
)


func main() {
	flag.Parse()

	// Load the portfolio data from the CSV file
	portfolio, err := simulation.LoadPortfolio(*portfolioPath)
	if err != nil {
		log.Fatalf("Failed to load portfolio: %v\n", err)
	}

	// Extract data from the portfolio
	var S0, weights, mu, sigma []float64
	for _, asset := range portfolio {
		S0 = append(S0, asset.InitialPrice)
		weights = append(weights, asset.Weight)
		mu = append(mu, asset.MeanReturn)
		sigma = append(sigma, asset.Volatility)
	}

	var result float64
	var sampledPayoffs []float64
	start := time.Now()
	switch *mode {
	case "sequential":
		result, sampledPayoffs = simulation.SequentialBasketOption(
			S0,
			weights,
			mu,
			sigma,
			*K,
			*r,
			*T,
			*steps,
			*simulations,
			*optionType, 
		)
	case "parallel":
		result, sampledPayoffs = simulation.ParallelBasketOption(
			S0,
			weights,
			mu,
			sigma,
			*K,
			*r,
			*T,
			*steps,
			*simulations,
			*numThreads,
			*optionType, 
		)
	case "parallel-stealing":
    result, sampledPayoffs = simulation.ParallelBasketOptionWithWorkStealing(
        S0,
        weights,
        mu,
        sigma,
        *K,
        *r,
        *T,
        *steps,
        *simulations,
        *numThreads,
        *optionType,
    )
	default:
		fmt.Println("Error: Invalid mode")
	}


	duration := time.Since(start)

	fmt.Printf("Basket Option Price: $%.2f\n", result)
	fmt.Printf("Execution Time (%s): %v\n", *mode, duration)
	visualization.ExportToCSV(sampledPayoffs, "results/payoff_distribution.csv")

}
