package mpu

import (
	"fmt"
	"log"

	"github.com/gentoomaniac/go64/cyclelock"
)

type ProcessorStatus int

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

	Memory *[0x10000]byte

	CyckleLock cyclelock.CycleLock
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

//DumpRegisters returns a string with the curremnt register states
func (m MOS6502) DumpRegisters() string {
	buffer := ""
	buffer += fmt.Sprintf("PC: 0x%04x\tPCL: 0x%02x\tPCH: 0x%02x\n", m.pc, m.PCL(), m.PCH())
	buffer += fmt.Sprintf("S: 0x%02x\nP: 0x%02x\nA: 0x%02x\nX: 0x%02x\nY: 0x%02x\n", m.s, m.p, m.a, m.x, m.y)

	return buffer
}

func (m *MOS6502) setProcessorStatusBit(status ProcessorStatus, isSet bool) {
	if isSet {
		m.p = m.p | m.s
	} else {
		m.p = m.p & (0xff ^ m.s)
	}
}

func (m MOS6502) getByteFromMemory(addr uint16) byte {
	m.CyckleLock.EnterCycle()
	b := m.Memory[addr]
	m.CyckleLock.ExitCycle()
	fmt.Printf("byte loaded from address 0x%04x: 0x%02x\n", addr, b)

	return b
}

func (m *MOS6502) storeByteInMemory(addr uint16, value byte) {
	m.CyckleLock.EnterCycle()
	m.Memory[addr] = value
	m.CyckleLock.ExitCycle()
	//log.Printf("Byte loaded from address 0x%04x: 0x%02x", addr, value)
}

func (m *MOS6502) getNextCodeByte() byte {
	b := m.getByteFromMemory(m.pc)
	m.pc++
	return b
}

func (m MOS6502) getDWordFromMemory(hi uint16, lo uint16) uint16 {
	word := uint16(m.getByteFromMemory(hi)) << 8
	result := word | uint16(m.getByteFromMemory(lo))

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

func (m *MOS6502) getDWordFromZeropage(addr byte) uint16 {
	return m.getDWordFromMemoryByAddr(uint16(addr), true)
}

func (m *MOS6502) getNextCodeDWord() uint16 {
	word := m.getDWordFromMemory(m.pc, m.pc+1)
	m.pc += 2
	return word
}

// Init initialises the MPU
func (m *MOS6502) Init(cyclelock cyclelock.CycleLock) {
	log.SetFlags(log.Lmicroseconds)
	log.SetFlags(log.Lshortfile)
	m.s = 0xff
	m.CyckleLock = cyclelock
}

// Run starts the execution of the MPU
func (m *MOS6502) Run() {
	// https://www.pagetable.com/?p=410
	// load reset vector
	//m.pc = getWordFromMemory(0xfffc);
	//log.Debug(string.Format("loading reset vector took {0} cycles", cycleLock.getCycleCount()));

	for i := 0; i <= 0xffff; i++ {
		m.getNextCodeByte()
		//log.Printf("++ CycleCount: %d\n", m.cyckleLock.CycleCount())
		m.CyckleLock.ResetCycleCount()
	}
}
