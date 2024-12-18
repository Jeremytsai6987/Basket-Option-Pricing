package simulation

import (
	"math"
	"math/rand"
)

func GeneratePricePath(S0, mu, sigma, T float64, steps int) []float64 {
	dt := T / float64(steps)
	prices := make([]float64, steps+1)
	prices[0] = S0

	for i := 1; i <= steps; i++ {
		normal := rand.NormFloat64()
		prices[i] = prices[i-1] * math.Exp((mu-0.5*sigma*sigma)*dt + sigma*math.Sqrt(dt)*normal)
	}

	return prices
}

func GeneratePricePathWithRand(S0, mu, sigma, T float64, steps int, randGen *rand.Rand) []float64 {
    dt := T / float64(steps)
    prices := make([]float64, steps+1)
    prices[0] = S0

    for i := 1; i <= steps; i++ {
        z := randGen.NormFloat64() 
        prices[i] = prices[i-1] * math.Exp((mu-0.5*sigma*sigma)*dt + sigma*math.Sqrt(dt)*z)
    }

    return prices
}

