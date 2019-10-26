package helper

import (
	"log"
	"time"
)

type Benchmark struct {
	start time.Time
	name  string
}

func NewBenchmark(name string) Benchmark {
	return Benchmark{time.Now(), name}
}

func (b Benchmark) Println() {
	log.Printf("%s took %v\n", b.name, time.Now().Sub(b.start))
}
