package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {

	// **** execute_consumer command: go tool pprof cpu.prof *****
	// > pprof web

	/*
		// profiling
		var cpuprofile = "cpu.prof"
		//var memprofile = "mem.prof"

		fCPU, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal("Could not create CPU profile: ", err)
		}
		defer fCPU.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(fCPU); err != nil {
			log.Fatal("Could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	*/
	// Print our starting memory usage (should be around 0mb)
	PrintMemUsage()

	var overall [][]int
	for i := 0; i < 4; i++ {

		// Allocate memory using make() and append to overall (so it doesn't get
		// garbage collected). This is to create an ever increasing memory usage
		// which we can track. We're just using []int as an example.
		a := make([]int, 0, 999999)
		overall = append(overall, a)

		// Print our memory usage at each interval
		PrintMemUsage()
		time.Sleep(5 * time.Second)
	}

	// Clear our memory and print usage, unless the GC has run 'Alloc' will remain the same
	overall = nil

	PrintMemUsage()

	// Force GC to clear up, should see a memory drop
	runtime.GC()
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	//fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("Alloc = %v Bytes", m.Alloc)
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	//fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tSys = %v Bytes", m.Sys)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
