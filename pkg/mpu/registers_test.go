package mpu

import (
	"math/rand"
	"testing"
	"time"

	"github.com/franela/goblin"
)

const (
	RandomTestCount int = 100
)

// TestMOS6502 tests the basic registers, and helpers
func TestMOS6502Registers(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Regsisters", func() {
		g.It("PC", func() {
			var blankMemory [0x10000]byte

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			value := uint16(0x0000)
			MOS6502.SetPC(value)
			g.Assert(MOS6502.PC()).Equal(value)

			value = uint16(0xffff)
			MOS6502.SetPC(value)
			g.Assert(MOS6502.PC()).Equal(value)

			for i := 0; i < RandomTestCount; i++ {
				value = uint16(rand.Intn(0xffff))
				MOS6502.SetPC(value)
				g.Assert(MOS6502.PC()).Equal(value)
			}
		})

		g.It("PCL", func() {
			var blankMemory [0x10000]byte
			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			value := uint8(0x00)
			MOS6502.SetPCL(value)
			g.Assert(MOS6502.PCL()).Equal(value)
			g.Assert(MOS6502.PC()).Equal(uint16(value))

			value = uint8(0xff)
			MOS6502.SetPCL(value)
			g.Assert(MOS6502.PCL()).Equal(value)
			g.Assert(MOS6502.PC()).Equal(uint16(value))

			for i := 0; i < RandomTestCount; i++ {
				value = uint8(rand.Intn(0xff))
				MOS6502.SetPCL(value)
				g.Assert(MOS6502.PCL()).Equal(value)
				g.Assert(MOS6502.PC()).Equal(uint16(value))
			}
		})

		g.It("PCH", func() {
			var blankMemory [0x10000]byte
			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			value := uint8(0x00)
			pcValue := uint16(value) << 8
			MOS6502.SetPCH(value)
			g.Assert(MOS6502.PCH()).Equal(value)
			g.Assert(MOS6502.PC()).Equal(uint16(pcValue))

			value = uint8(0xff)
			pcValue = uint16(value) << 8
			MOS6502.SetPCH(value)
			g.Assert(MOS6502.PCH()).Equal(value)
			g.Assert(MOS6502.PC()).Equal(uint16(pcValue))

			for i := 0; i < RandomTestCount; i++ {
				value = uint8(rand.Intn(0xff))
				pcValue = uint16(value) << 8
				MOS6502.SetPCH(value)
				g.Assert(MOS6502.PCH()).Equal(value)
				g.Assert(MOS6502.PC()).Equal(uint16(pcValue))
			}
		})

		g.It("setting processor status bits", func() {
			var blankMemory [0x10000]byte

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			MOS6502.p = 0

			MOS6502.setProcessorStatusBit(N, true)
			g.Assert(MOS6502.p).Equal(uint8(N))
			MOS6502.setProcessorStatusBit(N, false)
			g.Assert(MOS6502.p).Equal(uint8(0))

			MOS6502.setProcessorStatusBit(B, true)
			g.Assert(MOS6502.p).Equal(uint8(B))
			MOS6502.setProcessorStatusBit(B, false)
			g.Assert(MOS6502.p).Equal(uint8(0))

			MOS6502.setProcessorStatusBit(C, true)
			g.Assert(MOS6502.p).Equal(uint8(C))
			MOS6502.setProcessorStatusBit(C, false)
			g.Assert(MOS6502.p).Equal(uint8(0))

			MOS6502.setProcessorStatusBit(N, true)
			MOS6502.setProcessorStatusBit(V, true)
			MOS6502.setProcessorStatusBit(X, true)
			MOS6502.setProcessorStatusBit(B, true)
			MOS6502.setProcessorStatusBit(D, true)
			MOS6502.setProcessorStatusBit(I, true)
			MOS6502.setProcessorStatusBit(Z, true)
			MOS6502.setProcessorStatusBit(C, true)
			g.Assert(MOS6502.p).Equal(uint8(0xff))
		})

		g.It("setting bit masks", func() {
			g.Assert(setBits(byte(0x00), byte(0xff), true)).Equal(byte(0xff))
			g.Assert(setBits(byte(0xff), byte(0xaa), false)).Equal(byte(0x55))
		})
	})
}
