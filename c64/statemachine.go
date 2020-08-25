package c64

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gentoomaniac/go64/cyclelock"

	"github.com/gentoomaniac/go64/mpu"
)

const (
	// MaxMemoryAddress is the highest adressable Memory position
	MaxMemoryAddress uint16 = 0xffff
)

const (
	LORAM  byte = 0x01
	HIRAM       = 0x02
	CHAREN      = 0x04
)

// C64 represents the internal state of the system
type C64 struct {
	KernalRom    []byte
	BasicRom     []byte
	CharacterRom []byte

	// Memory represents the 64kB memory of the C64
	Memory [int(MaxMemoryAddress) + 1]byte

	// Mpu represents the MOS6502 of the C64
	Mpu     mpu.MOS6502
	mpuLock cyclelock.CycleLock
}

// DumpMemory debug prints the memory in the given address range
func (c C64) DumpMemory(start uint16, end uint16) string {
	dump := ""
	bytesPerRow := 16

	startAddr := int(start) - (int(start) % bytesPerRow)

	for index := int(startAddr); index < int(end); index += bytesPerRow {
		hexStrings := make([]string, 0)
		for _, value := range c.Memory[index : index+bytesPerRow] {
			hexStrings = append(hexStrings, fmt.Sprintf("%02x", value))
		}

		asText := string(c.Memory[index : index+bytesPerRow])
		reNonPrintabel := regexp.MustCompile("[^[:graph:] ]")
		asText = reNonPrintabel.ReplaceAllString(asText, ".")

		dump += fmt.Sprintf("0x%04x %s %s", index, strings.Join(hexStrings, " "), asText) + "\n"
	}

	return dump
}

func (c *C64) updateMemoryBanks() {
	// http://www.zimmers.net/anonftp/pub/cbm/maps/C64.MemoryMap
	if c.Memory[0x01]&LORAM == 1 {
		copy(c.Memory[0xa000:0xa000+len(c.BasicRom)], c.BasicRom)
	}

	if c.Memory[0x01]&HIRAM == 1 {
		copy(c.Memory[0xe000:0xe000+len(c.KernalRom)], c.KernalRom)
	}
	if c.Memory[0x01]&CHAREN == 1 {
		fmt.Println("ToDo: CHAREN is set, I/O should be mapped")
	} else {
		copy(c.Memory[0xd000:0xd000+len(c.CharacterRom)], c.CharacterRom)
	}
}

// Init initialises all components (loading roms, setting specific memory values etc)
func (c *C64) Init(basicRom string, kernalRom string, characterRom string) {
	log.SetFlags(log.Lmicroseconds)
	log.SetFlags(log.Lshortfile)

	var err error
	c.BasicRom, err = ioutil.ReadFile(basicRom)
	if err != nil {
		log.Panic(err)
	}

	c.KernalRom, err = ioutil.ReadFile(kernalRom)
	if err != nil {
		log.Panic(err)
	}

	c.CharacterRom, err = ioutil.ReadFile(characterRom)
	if err != nil {
		log.Panic(err)
	}

	c.Memory[0x00] = 0xff
	c.Memory[0x01] = 0x07

	c.updateMemoryBanks()

	c.Mpu.Memory = &c.Memory

	channelLock := &cyclelock.ChannelLock{}
	channelLock.Init()
	c.mpuLock = channelLock
	c.Mpu.Init(c.mpuLock)

	fmt.Println(c.Mpu.DumpRegisters())
}

// Run starts the simulation
func (c *C64) Run() {
	go c.Mpu.Run()

	time.Sleep(100 * time.Millisecond)
	cycle := 0
	for true {
		fmt.Printf("-- Cycle #%02d\n", cycle)
		c.mpuLock.Unlock()
		c.mpuLock.WaitForLock()
		cycle++
		time.Sleep(977 * time.Nanosecond)
	}
}
