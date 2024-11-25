package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/LucasFerence/prisoner-problem/stats"
)

type strategy = func() bool

// Execution flags
var numAttempts = flag.Int("n", 1000000, "Specify attempt count for each strategy")

func main() {
	flag.Parse()

	tryStrategy("Naive", naive)
	tryStrategy("Smart", smart)
}

func tryStrategy(name string, strat strategy) {
	stat := stats.Track(name)

	successCount := 0
	for i := 0; i < *numAttempts; i++ {
		op := stat.BeginOperation()
		if strat() {
			successCount++
		}
		stat.StopOperation(op)
	}

	avgSuccess := float64(successCount) / float64(*numAttempts)

	fmt.Println("------------------------------------------------------")
	fmt.Printf("Strategy [%s]: Success count [%d] over [%d] attempts. Average: [%f]\n", name, successCount, *numAttempts, avgSuccess)

	stat.EndTracking()

	fmt.Println("------------------------------------------------------")
}

// --- Strategies ---

func naive() bool {

	prisoners := createShuffled()
	for _, p := range prisoners {

		boxes := createShuffled()
		foundBox := false
		for bi := 0; bi < 50; bi++ {
			if boxes[bi] == p {
				foundBox = true
				break
			}
		}

		if !foundBox {
			return false
		}
	}

	return true
}

func smart() bool {

	prisoners := createShuffled()
	boxes := createShuffled()
	for _, p := range prisoners {

		foundBox := false
		nextBox := p
		for bi := 0; bi < 50; bi++ {
			if boxes[nextBox] == p {
				foundBox = true
				break
			} else {
				nextBox = boxes[nextBox]
			}
		}

		if !foundBox {
			return false
		}
	}

	return true
}

// --- Utility ---

// Create a shuffled list of numbers 1 to 100
func createShuffled() []int {
	var list = make([]int, 100)
	for i := 0; i < 100; i++ {
		list[i] = i
	}

	// shuffle the list
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}
