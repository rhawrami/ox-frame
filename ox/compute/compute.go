package compute

import (
	"runtime"
)

// NumWorkers defines the number of worker goroutines at any given moment. This value
// defaults to `runtime.NumCPU()`.
var NumWorkers int = runtime.NumCPU()
