package config

import (
	"runtime"

	"github.com/pbnjay/memory"
)

const (
	MinMemoryMB = 64
	MaxMemoryMB = 4096
	MinThreads  = 1
	MaxThreads  = 4
)

// Capabilities of the current computer.
// Used for Argon2 parameters, i.e. use this configuration to generate a key.
// Must be kept the same between write and read, so is saved alongside plaintext to prevent upgraded hardware from messing things up.
// Or even downgraded hardware to an extent, so long as you don't go from a supercomputer to a potato
//
// These specific int widths are chosen because that's what argon2.IDKey expects
type DerivationCapabilities struct {
	Time    uint32
	Memory  uint32
	Threads uint8
}

// obtain reasonable capabilities of the computer running.
func GetCurrentCapabilities() DerivationCapabilities {
	// get number of CPU cores available (i.e. number of threads we can use)
	threads := runtime.NumCPU()

	// clamp this to a reasonable range. we want at least 1
	// but at most a sane value (like 4 or 8), to avoid hogging all the resources
	threads = max(min(threads, MaxThreads), MinThreads)

	// obtain total amount of memory free, so we can use whatever is available. returned in bytes.
	// have to use a library function cause google are too lazy to include this in the already overflowing stdlib
	memToUse := memory.FreeMemory()
	memMB := memToUse >> 20 // convert to mebibytes

	// clamp to a reasonable range, so that we don't hog resources that might not be available come decryption time
	memMB = max(min(memMB, MaxMemoryMB), MinMemoryMB)
	memToUse = memMB << 10 // convert to kibibytes

	// spend more or less time depending on available threads
	time := 8 / threads

	return DerivationCapabilities{
		Time:    uint32(time),
		Memory:  uint32(memToUse),
		Threads: uint8(threads),
	}
}
