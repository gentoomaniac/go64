package mpu

import (
	"math/rand"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/gentoomaniac/go64/pkg/cyclelock"
)

func TestAdressingModes(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Test Adressing Modes", func() {
		g.It("resolves implied Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))

				g.Assert(MOS6502.impliedAdressing(address)).Equal(uint16(0))
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves accumulator Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				MOS6502.a = uint8(rand.Intn(0xff))

				g.Assert(MOS6502.accumulatorAdressing()).Equal(MOS6502.a)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves absolute Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))

				g.Assert(MOS6502.absoluteAdressing(address)).Equal(address)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves indexed Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))
				offset := uint8(rand.Intn(0xff))

				g.Assert(MOS6502.indexedAdressing(address, offset)).Equal(address + uint16(offset))
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves zeropage Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint8(rand.Intn(0xff))

				g.Assert(MOS6502.zeropageAdressing(address)).Equal(uint16(address))
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves zeropage indexed Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint8(rand.Intn(0xff))
				offset := uint8(rand.Intn(0xff))

				g.Assert(MOS6502.zeropageIndexedAdressing(address, offset)).Equal(uint16(address + offset))
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves relative Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xffff))
				offset := uint8(rand.Intn(0xff))

				g.Assert(MOS6502.relativeAdressing(address, offset)).Equal(address + uint16(offset))
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(0)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves absolute indirect Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xfffe))
				targetAddress := uint16(rand.Intn(0xfffe))

				blankMemory[address+1] = byte(targetAddress >> 8)
				blankMemory[address] = byte(targetAddress & 0xff)

				g.Assert(MOS6502.absoluteIndirectAdressing(address)).Equal(targetAddress)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves indexed indirect Adressing correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint8(rand.Intn(0xfe))
				targetAddress := uint16(rand.Intn(0xffff))

				MOS6502.x = uint8(rand.Intn(0xff))

				blankMemory[address+MOS6502.x+1] = byte(targetAddress >> 8)
				blankMemory[address+MOS6502.x] = byte(targetAddress & 0xff)

				g.Assert(MOS6502.indexedIndirectAdressing(address)).Equal(targetAddress)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(2)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})

		g.It("resolves indirect indexed Adressing without page boundary correctly", func() {
			var blankMemory [0x10000]byte

			var lock cyclelock.CycleLock
			lock = &cyclelock.AlwaysOpenLock{}

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory
			MOS6502.Init(lock)

			for i := 0; i < RandomTestCount; i++ {
				address := uint16(rand.Intn(0xfffe))
				intermediateAddress := uint16(rand.Intn(0xfffe))
				targetAddress := uint16(rand.Intn(0xffff))

				MOS6502.y = uint8(rand.Intn(0xff))

				blankMemory[address+1] = byte(intermediateAddress >> 8)
				blankMemory[address] = byte(intermediateAddress & 0xff)

				blankMemory[intermediateAddress+uint16(MOS6502.y)+1] = byte(targetAddress >> 8)
				blankMemory[intermediateAddress+uint16(MOS6502.y)] = byte(targetAddress & 0xff)

				g.Assert(MOS6502.indirectIndexedAdressing(address, false)).Equal(targetAddress)
				g.Assert(MOS6502.CycleLock.CycleCount()).Equal(4)

				MOS6502.CycleLock.ResetCycleCount()
			}
		})
	})
}
