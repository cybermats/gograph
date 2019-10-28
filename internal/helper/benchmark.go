package helper

import (
	"log"
	"time"
)

// Benchmark outputs the time from created to printed. Used for timing.
type Benchmark struct {
	start time.Time
	name  string
}

// NewBenchmark creates a new benchmark with the name used for later outputting.
func NewBenchmark(name string) Benchmark {
	return Benchmark{time.Now(), name}
}

// Println prints the name of the benchmark and how long it was since it was first created. Doesn't reset.
func (b Benchmark) Println() {
	log.Printf("%s took %v\n", b.name, time.Now().Sub(b.start))
}

func (b Benchmark) Duration() time.Duration {
	return time.Now().Sub(b.start)
}
