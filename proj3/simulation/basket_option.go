package simulation

import (
	"math"
	"math/rand"
	"sync"
)
type Deque struct {
	mu    sync.Mutex
	items []int
}

func (d *Deque) Push(item int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.items = append(d.items, item)
}

func (d *Deque) Pop() (int, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.items) == 0 {
		return 0, false
	}
	item := d.items[len(d.items)-1]
	d.items = d.items[:len(d.items)-1]
	return item, true
}

func (d *Deque) Steal() (int, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.items) == 0 {
		return 0, false
	}
	item := d.items[0]
	d.items = d.items[1:]
	return item, true
}


func NewRand(seed int64) *rand.Rand {
    return rand.New(rand.NewSource(seed))
}


func SampleData(data []float64, sampleSize int) []float64 {
    if len(data) <= sampleSize {
        return data 
    }

    sampled := make([]float64, sampleSize)
    for i := 0; i < sampleSize; i++ {
        sampled[i] = data[rand.Intn(len(data))]
    }
    return sampled
}



func ParallelBasketOption(S0, weights, mu, sigma []float64, K, r, T float64, steps, simulations, numThreads int, optionType string) (float64, []float64) {
    payoffs := make([]float64, simulations)
    tasks := make(chan int, simulations)
    results := make(chan struct {
        index  int
        payoff float64
    }, simulations)

    var wg sync.WaitGroup

    go func() {
        for i := 0; i < simulations; i++ {
            tasks <- i
        }
        close(tasks)
    }()

    for i := 0; i < numThreads; i++ {
        wg.Add(1)

        randGen := NewRand(42 + int64(i))

        go func(randGen *rand.Rand) {
            defer wg.Done()
            for task := range tasks {
                portfolioValue := 0.0

                for j := 0; j < len(S0); j++ {
                    pricePath := GeneratePricePathWithRand(S0[j], mu[j], sigma[j], T, steps, randGen)
                    finalPrice := pricePath[len(pricePath)-1]
                    portfolioValue += weights[j] * finalPrice
                }

                var payoff float64
                if optionType == "call" {
                    payoff = math.Max(portfolioValue-K, 0)
                } else if optionType == "put" {
                    payoff = math.Max(K-portfolioValue, 0)
                }

                results <- struct {
                    index  int
                    payoff float64
                }{index: task, payoff: payoff}
            }
        }(randGen)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        payoffs[result.index] = result.payoff
    }

    discountFactor := math.Exp(-r * T)
    totalPayoff := 0.0
    for _, payoff := range payoffs {
        totalPayoff += payoff
    }
    sampledPayoffs := SampleData(payoffs, 1000)      

    return discountFactor * (totalPayoff / float64(simulations)), sampledPayoffs
}

func SequentialBasketOption(S0, weights, mu, sigma []float64, K, r, T float64, steps, simulations int, optionType string) (float64, []float64) {
    payoffs := make([]float64, simulations)

    for i := 0; i < simulations; i++ {
        portfolioValue := 0.0

        for j := 0; j < len(S0); j++ {
            pricePath := GeneratePricePath(S0[j], mu[j], sigma[j], T, steps)
            finalPrice := pricePath[len(pricePath)-1]
            portfolioValue += weights[j] * finalPrice
        }

        if optionType == "call" {
            payoffs[i] = math.Max(portfolioValue-K, 0)
        } else if optionType == "put" {
            payoffs[i] = math.Max(K-portfolioValue, 0)
        }
    }

    discountFactor := math.Exp(-r * T)
    totalPayoff := 0.0
    for _, payoff := range payoffs {
        totalPayoff += payoff
    }
    sampledPayoffs := SampleData(payoffs, 1000)      


    return discountFactor * (totalPayoff / float64(simulations)), sampledPayoffs
}

func ParallelBasketOptionWithWorkStealing(S0, weights, mu, sigma []float64, K, r, T float64, steps, simulations, numThreads int, optionType string) (float64, []float64) {
    payoffs := make([]float64, simulations)
    deques := make([]*Deque, numThreads)
    results := make(chan struct {
        index  int
        payoff float64
    }, simulations)

    for i := 0; i < numThreads; i++ {
        deques[i] = &Deque{}
    }

    for i := 0; i < simulations; i++ {
        deques[i%numThreads].Push(i)
    }

    var wg sync.WaitGroup

    for i := 0; i < numThreads; i++ {
        wg.Add(1)
        randGen := NewRand(42 + int64(i))
        go func(threadID int, randGen *rand.Rand) {
            defer wg.Done()
            for {
                task, ok := deques[threadID].Pop()
                if !ok {
                    success := false
                    for j := 0; j < numThreads; j++ {
                        if j != threadID {
                            task, ok = deques[j].Steal()
                            if ok {
                                success = true
                                break
                            }
                        }
                    }
                    if !success {
                        return
                    }
                }

                portfolioValue := 0.0
                for j := 0; j < len(S0); j++ {
                    pricePath := GeneratePricePathWithRand(S0[j], mu[j], sigma[j], T, steps, randGen)
                    finalPrice := pricePath[len(pricePath)-1]
                    portfolioValue += weights[j] * finalPrice
                }

                var payoff float64
                if optionType == "call" {
                    payoff = math.Max(portfolioValue-K, 0)
                } else if optionType == "put" {
                    payoff = math.Max(K-portfolioValue, 0)
                }

                results <- struct {
                    index  int
                    payoff float64
                }{index: task, payoff: payoff}
            }
        }(i, randGen)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    for result := range results {
        payoffs[result.index] = result.payoff
    }

    discountFactor := math.Exp(-r * T)
    totalPayoff := 0.0
    for _, payoff := range payoffs {
        totalPayoff += payoff
    }

    sampledPayoffs := SampleData(payoffs, 1000)      
    return discountFactor * (totalPayoff / float64(simulations)), sampledPayoffs
}


