package core

import "math/bits"

type Intel8080Instruction struct {
	operation func()
	mnemonic  string
	cycles    uint16
}

type Intel8080 struct {
	// Registers
	a     byte
	b     byte
	c     byte
	d     byte
	e     byte
	flags *intel8080Flags
	h     byte
	l     byte

	sp uint16
	pc uint16

	memory       [0x4000]byte
	instructions [256]*Intel8080Instruction
}

func NewIntel8080() *Intel8080 {
	cpu := &Intel8080{
		flags: &intel8080Flags{},
	}

	cpu.instructions = [256]*Intel8080Instruction{
		{cpu._NOP, "NOP", 1}, {cpu._LXI_B, "LXI B", 3}, {cpu._STAX_B, "STAX B", 1}, {cpu._INX_B, "INX B", 1}, {cpu._INR_B, "INR B", 1}, {cpu._DCR_B, "DCR B", 1}, {cpu._MVI_B, "MVI B", 2},
	}

	return cpu
}

func (cpu *Intel8080) LoadProgram(program []byte) {
	copy(cpu.memory[:], program)
}

func (cpu *Intel8080) Run() {
	opcode := cpu.memory[cpu.pc]
	cpu.pc++

	instruction := cpu.instructions[opcode]
	instruction.operation()
}

func hasParity(b byte) bool {
	return bits.OnesCount8(b)%2 == 0
}

func (cpu *Intel8080) _NOP() {
	// No operation
}

func (cpu *Intel8080) _LXI_B() {
	cpu.c = cpu.memory[cpu.pc]
	cpu.b = cpu.memory[cpu.pc+1]
	cpu.pc += 2
}

func (cpu *Intel8080) _STAX_B() {
	addr := uint16(cpu.b)<<8 | uint16(cpu.c)

	cpu.memory[addr] = cpu.a
}
func (cpu *Intel8080) _INX_B() {
	bc := uint16(cpu.b)<<8 | uint16(cpu.c)
	bc++

	cpu.b = byte(bc >> 8)
	cpu.c = byte(bc & 0x00FF)
}

func (cpu *Intel8080) _INR_B() {
	cpu.b++

	cpu.flags.Set(Zero, cpu.b == 0)
	cpu.flags.Set(Sign, cpu.b&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.b&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.b))
}

func (cpu *Intel8080) _DCR_B() {
	cpu.b--

	cpu.flags.Set(Zero, cpu.b == 0)
	cpu.flags.Set(Sign, cpu.b&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.b&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.b))
}

func (cpu *Intel8080) _MVI_B() {
	cpu.b = cpu.memory[cpu.pc]
	cpu.pc++
}
