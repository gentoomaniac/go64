package mpu

import (
	"fmt"

	"github.com/gentoomaniac/go64/pkg/cyclelock"
	"github.com/gentoomaniac/go64/pkg/memory"

	"github.com/rs/zerolog/log"
)

type ProcessorStatus uint8

const (
	N ProcessorStatus = 0x80 // Negative
	V                 = 0x40 // Overflow
	X                 = 0x20 // unused
	B                 = 0x10 // Break
	D                 = 0x08 // Decimal mode
	I                 = 0x04 // Interrupt disabled
	Z                 = 0x02 // Zero
	C                 = 0x01 // Carry

	PageSize    uint16 = 256
	StackOffset uint16 = 0x0100

	NMIVector   uint16 = 0xfffa
	ResetVector uint16 = 0xfffc
	IRQVector   uint16 = 0xfffe
)

// MOS6502 is a struct representing the internal state of the MOS 6510 MPU
type MOS6502 struct {

	/*  Program Counter

	This register points the address from which the next
	instruction byte (opcode or parameter) will be fetched.
	Unlike other registers, this one is 16 bits in length. The
	low and high 8-bit halves of the register are called PCL
	and PCH, respectively.

	The Program Counter may be read by pushing its value on
	the stack. This can be done either by jumping to a
	subroutine or by causing an interrupt. */
	pc uint16

	/*  Stack pointer

	The NMOS 65xx processors have 256 bytes of stack memory,
	ranging from $0100 to $01FF. The S register is a 8-bit
	offset to the stack page. In other words, whenever
	anything is being pushed on the stack, it will be stored
	to the address $0100+S.

	The Stack pointer can be read and written by transfering
	its value to or from the index register X (see below) with
	the TSX and TXS instructions. */
	s uint8

	/*  Processor status

	This 8-bit register stores the state of the processor. The
	bits in this register are called flags. Most of the flags
	have something to do with arithmetic operations.

	The P register can be read by pushing it on the stack
	(with PHP or by causing an interrupt). If you only need to
	read one flag, you can use the branch instructions.
	Setting the flags is possible by pulling the P register
	from stack or by using the flag set or clear instructions. */
	p uint8

	/*  Accumulator

	The accumulator is the main register for arithmetic and
	logic operations. Unlike the index registers X and Y, it
	has a direct connection to the Arithmetic and Logic Unit
	(ALU). This is why many operations are only available for
	the accumulator, not the index registers. */
	a uint8

	/*  Index register X

	This is the main register for addressing data with
	indices. It has a special addressing mode, indexed
	indirect, which lets you to have a vector table on the
	zero page. */
	x uint8

	/*  Index register Y

	The Y register has the least operations available. On the
	other hand, only it has the indirect indexed addressing
	mode that enables access to any memory place without
	having to use self-modifying code. */
	y uint8

	Memory *memory.Memory

	CycleLock cyclelock.CycleLock
}

// PC returns the value of the PC register
func (m MOS6502) PC() uint16 {
	return m.pc
}

// SetPC sets the value of the PC register
func (m *MOS6502) SetPC(value uint16) {
	m.pc = value
}

// PCH returns the value of the PCH register
func (m MOS6502) PCH() uint8 {
	return uint8(m.pc >> 8)
}

// SetPCH sets the value of the PCH register
func (m *MOS6502) SetPCH(value uint8) {
	m.pc = (uint16(value) << 8) | uint16(m.PCL())
}

// PCL returns the value of the PC register
func (m MOS6502) PCL() uint8 {
	return uint8(m.pc & 0x00ff)
}

// SetPCL sets the value of the PCL register
func (m *MOS6502) SetPCL(value uint8) {
	m.pc = (m.pc & 0xff00) | uint16(value)
}

// S returns the value of the S register
func (m MOS6502) S() uint8 {
	return m.s
}

// SetS sets the value of the S register
func (m *MOS6502) SetS(value uint8) {
	m.s = value
}

// P returns the value of the P register
func (m MOS6502) P() uint8 {
	return m.p
}

// SetP sets the value of the P register
func (m *MOS6502) SetP(value uint8) {
	m.p = value
}

// A returns the value of the A register
func (m MOS6502) A() uint8 {
	return m.a
}

// SetA sets the value of the A register
func (m *MOS6502) SetA(value uint8) {
	m.a = value
}

// X returns the value of the X register
func (m MOS6502) X() uint8 {
	return m.x
}

// SetX sets the value of the X register
func (m *MOS6502) SetX(value uint8) {
	m.x = value
}

// Y returns the value of the Y register
func (m MOS6502) Y() uint8 {
	return m.y
}

// SetY sets the value of the Y register
func (m *MOS6502) SetY(value uint8) {
	m.y = value
}

// DumpRegisters returns a string with the curremnt register states
func (m MOS6502) DumpRegisters() string {
	buffer := ""
	buffer += fmt.Sprintf("PC: 0x%04x\tPCL: 0x%02x\tPCH: 0x%02x\n", m.pc, m.PCL(), m.PCH())
	buffer += fmt.Sprintf("S: 0x%02x\nP: 0x%02x\nA: 0x%02x\nX: 0x%02x\nY: 0x%02x\n", m.s, m.p, m.a, m.x, m.y)

	return buffer
}

/* Memory Access */

func (m *MOS6502) setProcessorStatusBit(s ProcessorStatus, isSet bool) {
	if isSet {
		m.p = m.p | uint8(s)
	} else {
		m.p = m.p & uint8(0xff^s)
	}
}

func setBits(value byte, mask byte, set bool) byte {
	if set {
		return value | mask
	}
	return value & (0xff ^ mask)
}

func (m MOS6502) getByteFromMemory(addr uint16, lockToCycle bool) byte {
	if lockToCycle {
		m.CycleLock.EnterCycle()
	}
	b := m.Memory[addr]
	if lockToCycle {
		m.CycleLock.ExitCycle()
	}
	//fmt.Printf("byte loaded from address 0x%04x: 0x%02x\n", addr, b)

	return b
}

func (m *MOS6502) storeByteInMemory(addr uint16, value byte, lockToCycle bool) {
	if lockToCycle {
		m.CycleLock.EnterCycle()
	}
	m.Memory[addr] = value
	if lockToCycle {
		m.CycleLock.ExitCycle()
	}
	//log.Printf("Byte loaded from address 0x%04x: 0x%02x", addr, value)
}

func (m *MOS6502) getNextCodeByte() byte {
	b := m.getByteFromMemory(m.pc, true)
	m.pc++
	return b
}

func (m MOS6502) getDWordFromMemory(hi uint16, lo uint16) uint16 {
	word := uint16(m.getByteFromMemory(hi, true)) << 8
	result := word | uint16(m.getByteFromMemory(lo, true))

	//fmt.Printf("dword loaded from address 0x%02x: 0x%04x\n", lo, result)
	return result
}

func (m *MOS6502) getDWordFromMemoryByAddr(addr uint16, pageBoundry bool) uint16 {
	// if we ignore page boundries or the second byte is still on the same page
	if !pageBoundry || (addr%PageSize) < PageSize-1 {
		return m.getDWordFromMemory(addr+1, addr)
		// otherwise take the pages 0x??00 address for the high byte (see http://www.oxyron.de/html/opcodes02.html "The 6502 bugs")
	}

	return m.getDWordFromMemory((addr/(PageSize))<<8, addr)
}

func (m *MOS6502) getDWordFromZeropage(addr uint8) uint16 {
	return m.getDWordFromMemoryByAddr(uint16(addr), true)
}

func (m *MOS6502) getNextCodeDWord() uint16 {
	word := m.getDWordFromMemory(m.pc, m.pc+1)
	m.pc += 2
	return word
}

/* Methods to resolve different adressing modes */
// http://www.emulator101.com/6502-addressing-modes.html
// http://www.obelisk.me.uk/6502/addressing.html#:~:text=Indexed%20indirect%20addressing%20is%20normally,byte%20of%20the%20target%20address.

// ToDo: The 6502 bugs
// addressing, which is rather a "no addressing mode at all"-option: Instructions which do not address an arbitrary memory location only supports this mode.
func (m MOS6502) impliedAdressing(addr uint16) uint16 { return 0 }

// addressing, supported by bit-shifting instructions, turns the "action" of the operation towards the accumulator.
// ToDo:
func (m MOS6502) accumulatorAdressing() uint8 { return m.a }

func (m MOS6502) absoluteAdressing(addr uint16) uint16 { return addr }

// absolute addressing, indexed by either the X and Y index registers: These adds the index register to a base address, forming the final "destination" for the operation.
func (m MOS6502) indexedAdressing(addr uint16, offset uint8) uint16 { return addr + uint16(offset) }

// addressing, which is similar to absolute addressing, but only works on addresses within the zeropage.
func (m MOS6502) zeropageAdressing(addr uint8) uint16 { return uint16(addr) }

// Effective address is zero page address plus the contents of the given register (X, or Y).
func (m MOS6502) zeropageIndexedAdressing(addr uint8, offset uint8) uint16 {
	return uint16(addr + offset)
}

// addressing, which uses a single byte to specify the destination of conditional branches ("jumps") within 128 bytes of where the branching instruction resides.
func (m MOS6502) relativeAdressing(addr uint16, offset byte) uint16 { return addr + uint16(offset) }

// addressing, which takes the content of a vector as its destination address.
func (m *MOS6502) absoluteIndirectAdressing(addr uint16) uint16 {
	return m.getDWordFromMemoryByAddr(addr, false)
}

// addressing, which uses the X index register to select one of a range of vectors in zeropage and takes the address from that pointer. Extremely rarely used!
func (m *MOS6502) indexedIndirectAdressing(addr uint8) uint16 {
	return m.getDWordFromZeropage(uint8(addr + m.x))
}

// addressing, which adds the Y index register to the contents of a pointer to obtain the address. Very flexible instruction found in anything but the most trivial machine language routines!
func (m *MOS6502) indirectIndexedAdressing(addr uint16, pageBoundry bool) uint16 {
	return m.getDWordFromMemoryByAddr(
		m.getDWordFromMemoryByAddr(addr, pageBoundry)+uint16(m.y),
		pageBoundry)
}

func (m *MOS6502) indirectIndexedZeropageAdressing(addr uint8) uint8 {
	return (byte)(m.indirectIndexedAdressing(uint16(addr), true))
}

/* HELPERS */
// Save byte to stack
func (m *MOS6502) push(value byte, lockToCycle bool) {
	if lockToCycle {
		m.CycleLock.EnterCycle()
	}
	m.storeByteInMemory(StackOffset+uint16(m.s), value, false)
	m.s--
	if lockToCycle {
		m.CycleLock.ExitCycle()
	}
}

func (m *MOS6502) pop(lockToCycle bool) byte {
	if lockToCycle {
		m.CycleLock.EnterCycle()
	}
	m.s++
	value := m.getByteFromMemory(StackOffset+uint16(m.s), false)
	if lockToCycle {
		m.CycleLock.ExitCycle()
	}

	return value
}

func checkForOverflow(vOld byte, vNew byte) bool {
	if (vOld&0x80) == 0 && (vNew&0x80) != 0 {
		return true
	} else if (vOld&0x80) != 0 && (vNew&0x80) == 0 {
		return true
	}
	return false
}

// Init initialises the MPU
func (m *MOS6502) Init(cyclelock cyclelock.CycleLock) {
	m.s = 0xff
	m.CycleLock = cyclelock
}

// Run starts the execution of the MPU
func (m *MOS6502) Run() {
	// https://www.pagetable.com/?p=410
	// load reset vector
	//m.pc = getWordFromMemory(0xfffc);
	//log.Debug(string.Format("loading reset vector took {0} cycles", cycleLock.getCycleCount()));

	for i := 0; i <= 0xffff; i++ {
		m.getNextCodeByte()
		log.Debug().Int("cycleCount", m.CycleLock.CycleCount()).Msg("")
		m.CycleLock.ResetCycleCount()
	}
}
