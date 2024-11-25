package stats

import (
	"fmt"
	"time"
)

type accumulator interface {
	receive(*operation)
	report()
}

// --- Average Duration MS ---

type movingAverage struct {
	opCount int
	avgTime float64
}

func (acc *movingAverage) receive(op *operation) {
	dur := float64(op.duration())

	acc.opCount++
	a := float64(1 / acc.opCount)
	b := 1 - a
	acc.avgTime = a*dur + b*acc.avgTime
}

func (acc *movingAverage) report() {
	avgDurMs := acc.avgTime / float64(time.Millisecond)
	fmt.Printf("Average time of [%d] operations: [%fms]\n", acc.opCount, avgDurMs)
}

// --- Operation Sum ---

type durationSum struct {
	durSum time.Duration
}

func (acc *durationSum) receive(op *operation) {
	acc.durSum += op.duration()
}

func (acc *durationSum) report() {
	fmt.Printf("Total duration of all operations: [%v]\n", acc.durSum)
}
