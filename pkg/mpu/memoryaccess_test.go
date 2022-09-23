package mpu

import (
	"math/rand"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/gentoomaniac/go64/pkg/cyclelock"
	"github.com/gentoomaniac/go64/pkg/memory"
)

func TestMemoryReads(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Test Memory Reads", func() {
		g.It("reads byte from memory", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))
				value := byte(rand.Intn(0xff))

				blankMemory[address] = value

				g.Assert(MOS6502.getByteFromMemory(address, true)).Equal(value)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(1)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("reads dword from memory given hi and lo address", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xfffe))
				value := uint16(rand.Intn(0xffff))

				hi := byte(value >> 8)
				lo := byte(value & 0xff)
				blankMemory[address] = lo
				blankMemory[address+1] = hi

				g.Assert(MOS6502.getDWordFromMemory(address+1, address)).Equal(value)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("reads dword within memory page", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			value := uint16(rand.Intn(0xffff))

			hi := byte(value >> 8)
			lo := byte(value & 0xff)
			blankMemory[0x01fe] = lo
			blankMemory[0x01ff] = hi

			g.Assert(MOS6502.getDWordFromMemoryByAddr(0x01fe, false)).Equal(value)
			g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)
			MOS6502.CycleLock.ResetCycleCount()
			g.Assert(MOS6502.getDWordFromMemoryByAddr(0x01fe, true)).Equal(value)
			g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)
			MOS6502.CycleLock.ResetCycleCount()
		})

		g.It("reads dword with 6502 bug when crossing page boundary", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			value := uint16(rand.Intn(0xffff))

			hi := byte(value >> 8)
			lo := byte(value & 0xff)
			blankMemory[0x01ff] = lo
			blankMemory[0x0100] = hi

			g.Assert(MOS6502.getDWordFromMemoryByAddr(0x01ff, true)).Equal(value)
			g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)

			MOS6502.CycleLock.ResetCycleCount()
		})

		g.It("reads dword from zeropage with 6502 bug", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			value := uint16(rand.Intn(0xffff))

			hi := byte(value >> 8)
			lo := byte(value & 0xff)
			blankMemory[0x00ff] = lo
			blankMemory[0x0000] = hi

			g.Assert(MOS6502.getDWordFromZeropage(uint8(0x00ff))).Equal(value)
			g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)

			MOS6502.CycleLock.ResetCycleCount()
		})

		g.It("getNextCodeByte increments PC register by one ", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			oldPC := MOS6502.pc

			MOS6502.getNextCodeByte()

			g.Assert(MOS6502.pc - oldPC).Equal(uint16(1))

			MOS6502.CycleLock.ResetCycleCount()
		})

		g.It("getNextCodeDWord increments PC register by two ", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			oldPC := MOS6502.pc

			MOS6502.getNextCodeDWord()

			g.Assert(MOS6502.pc - oldPC).Equal(uint16(2))

			MOS6502.CycleLock.ResetCycleCount()
		})
	})
}

func TestMemoryWrite(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Test Memory Writes", func() {
		g.It("write byte to memory", func() {
			var blankMemory memory.Memory

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))
				value := byte(rand.Intn(0xff))

				MOS6502.storeByteInMemory(address, value, true)

				g.Assert(MOS6502.Memory[address]).Equal(value)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(1)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})
	})
}
