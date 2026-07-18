package encryption

import (
	"runtime"

	"github.com/pbnjay/memory"
)

const (
	MinMemoryMB = 64
	MaxMemoryMB = 4096
	MinThreads  = 1
	MaxThreads  = 8
)

// Capabilities of the current computer. Used for Argon2 parameters, i.e. use the parameters generated with
type DerivationCapabilities struct {
	Memory  uint32
	Time    uint32
	Threads uint8
}

// obtain reasonable capabilities of the computer running.
func getCurrentCapabilities() DerivationCapabilities {
	// get number of CPU cores available (i.e. number of threads we can use)
	threads := runtime.NumCPU()

	// clamp this to a reasonable range (here, 1..=8). we want at least 1
	// but at most 8, to avoid hogging all the resources
	threads = max(min(threads, 8), 1)

	// total amount of memory free, so we can use whatever is available. returned in bytes.
	memToUse := memory.FreeMemory()
	memMB := memToUse >> 20 // convert to mebibytes

	// clamp to a reasonable range, to 512MiB..=4GiB, so that we don't hog resources
	// that might not be available come decryption time
	memMB = max(min(memMB, 4096), 512)
	memToUse = memMB << 10 // convert to kibibytes

	// spend more or less time depending on available threads
	time := 8 / threads

	return DerivationCapabilities{
		Memory:  uint32(memToUse),
		Time:    uint32(time),
		Threads: uint8(threads),
	}
}
