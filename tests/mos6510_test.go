package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/gentoomaniac/go64/cyclelock"
	"github.com/gentoomaniac/go64/mpu"

	"github.com/franela/goblin"
)

const (
	RandomTestCount int = 100
)

//TestMOS6510 tests the basic registers, and helpers
func TestMOS6510Registers(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Regsisters", func() {
		g.It("PC", func() {
			var blankMemory [0x10000]byte
			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory

			value := uint16(0x0000)
			mos6510.SetPC(value)
			g.Assert(mos6510.PC()).Equal(value)

			value = uint16(0xffff)
			mos6510.SetPC(value)
			g.Assert(mos6510.PC()).Equal(value)

			for i := 0; i < RandomTestCount; i++ {
				value = uint16(rand.Intn(0xffff))
				mos6510.SetPC(value)
				g.Assert(mos6510.PC()).Equal(value)
			}
		})

		g.It("PCL", func() {
			var blankMemory [0x10000]byte
			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory

			value := uint8(0x00)
			mos6510.SetPCL(value)
			g.Assert(mos6510.PCL()).Equal(value)
			g.Assert(mos6510.PC()).Equal(uint16(value))

			value = uint8(0xff)
			mos6510.SetPCL(value)
			g.Assert(mos6510.PCL()).Equal(value)
			g.Assert(mos6510.PC()).Equal(uint16(value))

			for i := 0; i < RandomTestCount; i++ {
				value = uint8(rand.Intn(0xff))
				mos6510.SetPCL(value)
				g.Assert(mos6510.PCL()).Equal(value)
				g.Assert(mos6510.PC()).Equal(uint16(value))
			}
		})

		g.It("PCH", func() {
			var blankMemory [0x10000]byte
			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory

			value := uint8(0x00)
			pcValue := uint16(value) << 8
			mos6510.SetPCH(value)
			g.Assert(mos6510.PCH()).Equal(value)
			g.Assert(mos6510.PC()).Equal(uint16(pcValue))

			value = uint8(0xff)
			pcValue = uint16(value) << 8
			mos6510.SetPCH(value)
			g.Assert(mos6510.PCH()).Equal(value)
			g.Assert(mos6510.PC()).Equal(uint16(pcValue))

			for i := 0; i < RandomTestCount; i++ {
				value = uint8(rand.Intn(0xff))
				pcValue = uint16(value) << 8
				mos6510.SetPCH(value)
				g.Assert(mos6510.PCH()).Equal(value)
				g.Assert(mos6510.PC()).Equal(uint16(pcValue))
			}
		})
	})
}

func TestMemoryReads(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	memory := goblin.Goblin(t)
	memory.Describe("Test Memory Reads", func() {
		memory.It("reads byte from memory", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory
			mos6510.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))
				value := byte(rand.Intn(0xff))

				blankMemory[address] = value

				memory.Assert(mos6510.GetByteFromMemory(address)).Equal(value)
				memory.Assert(mos6510.CyckleLock.CycleCount()).Equal(1)

				mos6510.CyckleLock.ResetCycleCount()
			}
		})

		memory.It("reads dword from memory given hi and lo address", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory
			mos6510.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xfffe))
				value := uint16(rand.Intn(0xffff))

				hi := byte(value >> 8)
				lo := byte(value & 0xff)
				blankMemory[address] = lo
				blankMemory[address+1] = hi

				memory.Assert(mos6510.GetDWordFromMemory(address+1, address)).Equal(value)
				memory.Assert(mos6510.CyckleLock.CycleCount()).Equal(2)

				mos6510.CyckleLock.ResetCycleCount()
			}
		})

		memory.It("reads dword within memory page", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory
			mos6510.Init(lock)

			value := uint16(rand.Intn(0xffff))

			hi := byte(value >> 8)
			lo := byte(value & 0xff)
			blankMemory[0x01fe] = lo
			blankMemory[0x01ff] = hi

			memory.Assert(mos6510.GetDWordFromMemoryByAddr(0x01fe, false)).Equal(value)
			memory.Assert(mos6510.GetDWordFromMemoryByAddr(0x01fe, true)).Equal(value)

			mos6510.CyckleLock.ResetCycleCount()
		})

		memory.It("reads dword with 6502 bug when crossing page boundary", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			mos6510 := &mpu.MOS6510{}
			mos6510.Memory = &blankMemory
			mos6510.Init(lock)

			value := uint16(rand.Intn(0xffff))

			hi := byte(value >> 8)
			lo := byte(value & 0xff)
			blankMemory[0x01ff] = lo
			blankMemory[0x0100] = hi

			memory.Assert(mos6510.GetDWordFromMemoryByAddr(0x01ff, true)).Equal(value)

			mos6510.CyckleLock.ResetCycleCount()
		})
	})
}
