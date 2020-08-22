package c64

import (
	"github.com/gentoomaniac/go64/mpu"
)

const (
	// MaxMemoryAddress is the highest adressable Memory position
	MaxMemoryAddress uint16 = 0xffff
)

// C64 represents the internal state of the system
type C64 struct {

	// Memory represents the 64kB memory of the C64
	Memory [int(MaxMemoryAddress) + 1]byte

	// Mpu represents the MOS6510 of the C64
	Mpu mpu.MOS6510
}

// DumpMemory debug prints the memory in the given address range
func (c C64) DumpMemory(min uint16, max uint16) {

}
