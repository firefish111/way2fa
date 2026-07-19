package config

import (
	"math/bits"
	"runtime"

	"github.com/pbnjay/memory"
)

const (
	MinThreads = 1
	MaxThreads = 4

	// 1 << 26 = 64MiB, the absolute minimum recommendation for Argon2
	MinMemoryBytesLog2 = 26
	// 1 << 32 = 4GiB, a reasonable limit to not hog resources
	MaxMemoryBytesLog2 = 32
)

// Capabilities of the current computer.
// Used for Argon2 parameters, i.e. use this configuration to generate a key.
// Must be kept the same between write and read, so is saved alongside plaintext to prevent upgraded hardware from messing things up.
// Or even downgraded hardware to an extent, so long as you don't go from a supercomputer to a potato
//
// These specific int widths are chosen because that's what argon2.IDKey expects
type DerivationCapabilities struct {
	MemKiB  uint32
	Time    uint32
	Threads uint8
}

// obtain reasonable capabilities of the computer running.
func GetCurrentCapabilities() (dc DerivationCapabilities) {
	// get number of CPU cores available (i.e. number of threads we can use)
	threads := runtime.NumCPU()

	// clamp this to a reasonable range. we want at least 1
	// but at most a sane value (like 4 or 8), to avoid hogging all the resources
	threads = max(min(threads, MaxThreads), MinThreads)

	dc.Threads = uint8(threads)

	// obtain total amount of memory free, so we can use whatever is available. returned in bytes.
	// have to use a library function cause google are too lazy to include this in the already overflowing stdlib
	memFree := memory.FreeMemory()

	// get log2. this rounds DOWN to the nearest power of 2.
	// the or with 1 prevents wrapping to math.MaxUint64 when memFree = 0
	memLog2 := bits.Len64(memFree|1) - 1

	// clamp to a reasonable range, so that we don't hog resources that might not be available come decryption time
	// see the constants' definition for the rational behind choosing them
	memLog2 = max(min(memLog2, MaxMemoryBytesLog2), MinMemoryBytesLog2)

	// subtract 10 to convert to KiB, then shift left 1 by that
	memKiB := 1 << (memLog2 - 10) // convert to kibibytes

	dc.MemKiB = uint32(memKiB)

	// spend more or less time depending on other variables
	// between 0 and 3 seconds, more time for less memory
	time := (MaxMemoryBytesLog2 - memLog2) / 2
	// an extra 0-2 depending on threads
	time += threads / 2

	// set to have a minimum of 1 second
	dc.Time = uint32(max(time, 1))

	return
}

// packed version of the struct for use in the header.
// has the following memory layout:
// - bits 0-7:   log2(MemKiB)
// - bits 8-11:  Time
// - bits 12-15: Threads
// This is because Time and Threads can never exceed 15, so can easily be pushed together.
type DerivationCapabilitiesPacked uint16

func (dc DerivationCapabilities) Pack() (packed DerivationCapabilitiesPacked) {
	// set all bits to 0
	packed = 0

	// set low eight bits to log2(MemKiB). see above for why we or with 1
	memKiBLog2 := bits.Len32(dc.MemKiB|1) - 1
	packed = (DerivationCapabilitiesPacked(memKiBLog2) & 0xff)

	packed |= (DerivationCapabilitiesPacked(dc.Time) & 0x0f) << 8
	packed |= (DerivationCapabilitiesPacked(dc.Threads) & 0x0f) << 12
	return
}

func (packed DerivationCapabilitiesPacked) Unpack() (dc DerivationCapabilities) {
	dc.MemKiB = 1 << uint32((packed & 0x00ff))
	dc.Time = uint32((packed & 0x0f00) >> 8)
	dc.Threads = uint8((packed & 0xff00) >> 12)
	return
}
