package utils

import (
	"fmt"
	"repo.abanicon.com/abantheter-microservices/websocket/app/configs"
	"runtime"
	"time"
)

type Profiling struct {
	config configs.Debug
}

func NewProfiling(config configs.Debug) *Profiling {
	return &Profiling{config: config}
}

func (p Profiling) Register() {
	if p.config.Active {
		p.run()
	}
}

func (p Profiling) run() {
	tiker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-tiker.C
			if p.config.GC {
				runtime.GC()
			}
			if p.config.Allocation {
				printStats()
			}
		}
	}()
}

func printStats() {
	// For memory
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	fmt.Printf("\tGoRoutines = %v\n", runtime.NumGoroutine())
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
