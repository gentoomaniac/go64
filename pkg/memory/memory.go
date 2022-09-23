package memory

import (
	"fmt"
	"regexp"
	"strings"
)

type Memory [0x10000]byte

func (m *Memory) Set(addr uint16, value byte) {
	m[addr] = value
}

func (m Memory) Get(addr uint16) byte {
	return m[addr]
}

func (m *Memory) CopyTo(offset uint16, array []byte) {
	copy(m[offset:offset+uint16(len(array))], array)
}

// DumpMemory debug prints the memory in the given address range
func (m Memory) DumpMemory(start uint16, end uint16) string {
	dump := ""
	bytesPerRow := 16

	startAddr := int(start) - (int(start) % bytesPerRow)

	for index := int(startAddr); index < int(end); index += bytesPerRow {
		hexStrings := make([]string, 0)
		for _, value := range m[index : index+bytesPerRow] {
			hexStrings = append(hexStrings, fmt.Sprintf("%02x", value))
		}

		asText := string(m[index : index+bytesPerRow])
		reNonPrintabel := regexp.MustCompile("[^[:graph:] ]")
		asText = reNonPrintabel.ReplaceAllString(asText, ".")

		dump += fmt.Sprintf("0x%04x %s %s", index, strings.Join(hexStrings, " "), asText) + "\n"
	}

	return dump
}
