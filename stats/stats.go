package stats

import (
	"fmt"
	"time"
)

type stats struct {
	name         string

	trackingStart time.Time
	trackingDuration time.Duration

	accumulators []accumulator
	opChan       chan *operation
	done         chan bool
}

type operation struct {
	start time.Time
	stop  time.Time
}

func (op *operation) duration() time.Duration {
	return op.stop.Sub(op.start)
}

func defaultAccumulators() []accumulator {
	accumulators := []accumulator{}

	accumulators = append(accumulators, &movingAverage{0, 0})
	accumulators = append(accumulators, &durationSum{0})

	return accumulators
}

/*
Begin tracking stats under the defined name.
This will use default accumulators to track specific stats
*/
func Track(name string) *stats {

	stats := stats{
		name,

		time.Now(),
		0,

		defaultAccumulators(),
		make(chan *operation, 100),
		make(chan bool),
	}

	// create a thread constantly reading from the op channel
	go func() {
		for op := range stats.opChan {
			for _, col := range stats.accumulators {
				col.receive(op)
			}
		}

		// signal to the done channel that it has read all of the operations
		stats.done <- true
	}()

	return &stats
}

/*
Complete the tracking for this instance of stats
*/
func (s *stats) EndTracking() {
	s.trackingDuration = time.Now().Sub(s.trackingStart)

	close(s.opChan)
	<-s.done
}

/*
Output the report of the stats to stdout
*/
func (s *stats) PrintReport() {
	// generate the report for all accumulators
	fmt.Printf("Completed tracking for [%s] over time [%v]\n", s.name, s.trackingDuration)
	for _, col := range s.accumulators {
		col.report()
	}
}

/*
Begin an operation for the instance of stats
*/
func (s *stats) BeginOperation() *operation {
	return &operation{
		time.Now(),
		time.Time{},
	}
}

/*
Complete an operation for the instance of stats
*/
func (s *stats) StopOperation(op *operation) {
	op.stop = time.Now()
	s.opChan <- op
}
