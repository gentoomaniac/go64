package c64

import (
	"fmt"
	"regexp"
	"strings"

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
func (c C64) DumpMemory(start uint16, end uint16) string {
	dump := ""
	bytesPerRow := 16

	startAddr := int(start) - (int(start) % bytesPerRow)

	for index := int(startAddr); index < int(end); index += bytesPerRow {
		hexStrings := make([]string, 0)
		for _, value := range c.Memory[index : index+bytesPerRow-1] {
			hexStrings = append(hexStrings, fmt.Sprintf("%02x", value))
		}

		asText := string(c.Memory[index : index+bytesPerRow-1])
		reNonPrintabel := regexp.MustCompile("[^[:graph:] ]")
		asText = reNonPrintabel.ReplaceAllString(asText, ".")

		dump += fmt.Sprintf("0x%04x %s %s", index, strings.Join(hexStrings, " "), asText) + "\n"
	}

	return dump
}
