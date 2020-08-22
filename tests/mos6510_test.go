package tests

import (
	"math/rand"
	"testing"

	"github.com/gentoomaniac/go64/mpu"

	"github.com/franela/goblin"
)

const (
	RandomTestCount int = 100
)

//TestMOS6510 tests the basic registers, and helpers
func TestMOS6510(t *testing.T) {
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
