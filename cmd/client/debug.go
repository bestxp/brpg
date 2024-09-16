package main

import (
	"fmt"
	"runtime"
)

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	out := fmt.Sprintf("\n Alloc = %v MiB", bToMb(m.Alloc))
	out += fmt.Sprintf("\n Sys = %v MiB", bToMb(m.Sys))
	return out
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
