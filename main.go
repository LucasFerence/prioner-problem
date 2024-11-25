package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/LucasFerence/prisoner-problem/stats"
)

type strategy = func() bool

// Execution flags
var numAttempts = flag.Int("n", 1000000, "Specify attempt count for each strategy")
var prisonerCount = flag.Int("p", 100, "Amount of prisoners to test against")
var boxAllowance = flag.Int("b", 50, "Amount of boxes a prisoner can check")

func main() {
	initializeArgs()

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

	stat.EndTracking()

	stat.PrintReport()
	fmt.Printf("Strategy [%s]: Success count [%d] over [%d] attempts. Average: [%f]\n\n", name, successCount, *numAttempts, avgSuccess)
}

// --- Strategies ---

// strategies return true if all prisoners succeeded in their attempt
// they will return false if any of the prisoners failed to find their number

func naive() bool {

	prisoners := createShuffled()
	for _, p := range prisoners {

		// trick here is to shuffle the boxes rather than selecting random numbers
		// same number of rand calls, but we don't need to worry if we've checked a box before
		boxes := createShuffled()
		foundBox := false
		for bi := 0; bi < *boxAllowance; bi++ {
			if boxes[bi] == p {
				foundBox = true
				break
			}
		}

		// we can exit early if we didn't find the box, since the rest of the test doesn't matter
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
		for bi := 0; bi < *boxAllowance; bi++ {
			if boxes[nextBox] == p {
				foundBox = true
				break
			} else {
				nextBox = boxes[nextBox]
			}
		}

		// exit early if no box foudn on attempt
		if !foundBox {
			return false
		}
	}

	return true
}

// --- Utility ---

// Create a shuffled list of numbers 1 to 100
func createShuffled() []int {
	numPrisoners := *prisonerCount

	list := make([]int, numPrisoners)
	for i := 0; i < numPrisoners; i++ {
		list[i] = i
	}

	// shuffle the list
	for i := range list {
		j := rand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}

	return list
}

func initializeArgs() {

	flag.Parse()

	// verify args
	if *boxAllowance > *prisonerCount {
		fmt.Printf("ERROR! Box allowance [%d] out of bounds!\n", *boxAllowance)
		os.Exit(1)
	}

	fmt.Println("----------------------------------------")
	fmt.Println("Beginning test...")
	fmt.Printf("Prisoner count: %d\n", *prisonerCount)
	fmt.Printf("Box allowance: %d\n", *boxAllowance)
	fmt.Printf("Test executions: %d\n", *numAttempts)
	fmt.Println("----------------------------------------")
}
