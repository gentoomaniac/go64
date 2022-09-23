package mpu

import (
	"math/rand"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/gentoomaniac/go64/pkg/memory"
)

// TestSTackFunctions tests the stack functions
func TestStackFunctions(t *testing.T) {

	rand.Seed(time.Now().UTC().UnixNano())

	g := goblin.Goblin(t)
	g.Describe("Stack", func() {
		g.It("push decrements stackpointer and stores value", func() {
			var blankMemory memory.Memory

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			MOS6502.s = 0xff
			value := byte(rand.Intn(0xff))
			MOS6502.push(value, false)

			g.Assert(MOS6502.s).Equal(uint8(0xfe))
			g.Assert(MOS6502.Memory[StackOffset+uint16(MOS6502.s+1)]).Equal(value)
		})

		g.It("stack overflow behaves as expected", func() {
			var blankMemory memory.Memory

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			MOS6502.s = 0x00
			value := byte(rand.Intn(0xff))
			MOS6502.push(value, false)

			g.Assert(MOS6502.s).Equal(uint8(0xff))
			g.Assert(MOS6502.Memory[StackOffset+uint16(MOS6502.s+1)]).Equal(value)
		})

		g.It("pop increments stackpointer and retrieves value", func() {
			var blankMemory memory.Memory

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			MOS6502.s = 0xfe
			value := byte(rand.Intn(0xff))
			MOS6502.Memory[StackOffset+uint16(MOS6502.s+1)] = value

			g.Assert(MOS6502.pop(false)).Equal(value)
			g.Assert(MOS6502.s).Equal(uint8(0xff))
		})

		g.It("stack underflow behaves as expected", func() {
			var blankMemory memory.Memory

			MOS6502 := &MOS6502{}
			MOS6502.Memory = &blankMemory

			MOS6502.s = 0xff
			value := byte(rand.Intn(0xff))
			MOS6502.Memory[StackOffset+uint16(MOS6502.s+1)] = value

			g.Assert(MOS6502.pop(false)).Equal(value)
			g.Assert(MOS6502.s).Equal(uint8(0x00))
		})
	})
}
