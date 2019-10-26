package helper

import (
	"fmt"
	"log"
	"runtime"
)

// PrintMemUsage outputs the current, total and OS memory being used.
// As well as the number of garage collection cycles completed.
func PrintMemUsage() {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	str := fmt.Sprintf("Alloc = %v MiB", bToMb(m.Alloc))
	str = str + fmt.Sprintf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	str = str + fmt.Sprintf("\tSys = %v MiB", bToMb(m.Sys))
	str = str + fmt.Sprintf("\tNumGC = %v\n", m.NumGC)
	log.Print(str)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
