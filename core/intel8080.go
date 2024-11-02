package core

import (
	"log"
	"math/bits"
)

type Intel8080Instruction struct {
	operation func() uint
	mnemonic  string
	size      uint16
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

	sp     uint16
	pc     uint16
	cycles uint

	memory       [0x4000]byte
	instructions [256]*Intel8080Instruction
}

func NewIntel8080() *Intel8080 {
	cpu := &Intel8080{
		flags: &intel8080Flags{},
	}

	cpu.instructions = [256]*Intel8080Instruction{
		0x00: {cpu._NOP, "NOP", 1},
		0x01: {cpu._LXI_B, "LXI B", 3},
		0x02: {cpu._STAX_B, "STAX B", 1},
		0x03: {cpu._INX_B, "INX B", 1},
		0x04: {cpu._INR_B, "INR B", 1},
		0x05: {cpu._DCR_B, "DCR B", 1},
		0x06: {cpu._MVI_B, "MVI B", 2},
		0x07: {cpu._RLC, "RLC", 1},
		0x08: {cpu._NOP, "*NOP", 1},
		0x09: {cpu._DAD_B, "DAD B", 1},
		0x0a: {cpu._LDAX_B, "LDAX B", 1},
		0x0b: {cpu._DCX_B, "DCX B", 1},
		0x0c: {cpu._INR_C, "INR C", 1},
		0x0d: {cpu._DCR_C, "DCR C", 1},
		0x0e: {cpu._MVI_C, "MVI C", 2},
		0x0f: {cpu._RRC, "RRC", 1},

		0x10: {cpu._NOP, "*NOP", 1},
		0x11: {cpu._LXI_D, "LXI D", 3},
		0x12: {cpu._STAX_D, "STAX D", 1},
		0x13: {cpu._INX_D, "INX D", 1},
		0x14: {cpu._INR_D, "INR D", 1},
		0x15: {cpu._DCR_D, "DCR D", 1},
		0x16: {cpu._MVI_D, "MVI D", 2},
		0x17: {cpu._RAL, "RAL", 1},
		0x18: {cpu._NOP, "*NOP", 1},
		0x19: {cpu._DAD_D, "DAD D", 1},
		0x1a: {cpu._LDAX_D, "LDAX D", 1},
		0x1b: {cpu._DCX_D, "DCX D", 1},
		0x1c: {cpu._INR_E, "INR E", 1},
		0x1d: {cpu._DCR_E, "DCR E", 1},
		0x1e: {cpu._MVI_E, "MVI E", 2},
		0x1f: {cpu._RAR, "RAR", 1},

		0x20: {cpu._NOP, "*NOP", 1},
		0x21: {cpu._LXI_H, "LXI H", 3},
		0x22: {cpu._SHLD, "SHLD", 3},
		0x23: {cpu._INX_H, "INX H", 1},
		0x24: {cpu._INR_H, "INR H", 1},
		0x25: {cpu._DCR_H, "DCR H", 1},
		0x26: {cpu._MVI_H, "MVI H", 2},
		0x27: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0x28: {cpu._NOP, "*NOP", 1},
		0x29: {cpu._DAD_H, "DAD H", 1},
		0x2a: {cpu._LHLD, "LHLD", 3},
		0x2b: {cpu._DCX_H, "DCX H", 1},
		0x2c: {cpu._INR_L, "INR L", 1},
		0x2d: {cpu._DCR_L, "DCR L", 1},
		0x2e: {cpu._MVI_L, "MVI L", 2},
		0x2f: {cpu._CMA, "CMA", 1},

		0x30: {cpu._NOP, "*NOP", 1},
		0x31: {cpu._LXI_SP, "LXI SP", 3},
		0x32: {cpu._STA, "STA", 3},
		0x33: {cpu._INX_SP, "INX SP", 1},
		0x34: {cpu._INR_M, "INR M", 1},
		0x35: {cpu._DCR_M, "DCR M", 1},
		0x36: {cpu._MVI_M, "MVI M", 2},
		0x37: {cpu._STC, "STC", 1},
		0x38: {cpu._NOP, "*NOP", 1},
		0x39: {cpu._DAD_SP, "DAD SP", 1},
		0x3a: {cpu._LDA, "LDA", 3},
		0x3b: {cpu._DCX_SP, "DCX SP", 1},
		0x3c: {cpu._INR_A, "INR A", 1},
		0x3d: {cpu._DCR_A, "DCR A", 1},
		0x3e: {cpu._MVI_A, "MVI A", 2},
		0x3f: {cpu._CMC, "CMC", 1},

		0x40: {cpu._NI, "Not Impl", 0},
		0x41: {cpu._NI, "Not Impl", 0},
		0x42: {cpu._NI, "Not Impl", 0},
		0x43: {cpu._NI, "Not Impl", 0},
		0x44: {cpu._NI, "Not Impl", 0},
		0x45: {cpu._NI, "Not Impl", 0},
		0x46: {cpu._NI, "Not Impl", 0},
		0x47: {cpu._NI, "Not Impl", 0},
		0x48: {cpu._NI, "Not Impl", 0},
		0x49: {cpu._NI, "Not Impl", 0},
		0x4a: {cpu._NI, "Not Impl", 0},
		0x4b: {cpu._NI, "Not Impl", 0},
		0x4c: {cpu._NI, "Not Impl", 0},
		0x4d: {cpu._NI, "Not Impl", 0},
		0x4e: {cpu._NI, "Not Impl", 0},
		0x4f: {cpu._NI, "Not Impl", 0},

		0x50: {cpu._NI, "Not Impl", 0},
		0x51: {cpu._NI, "Not Impl", 0},
		0x52: {cpu._NI, "Not Impl", 0},
		0x53: {cpu._NI, "Not Impl", 0},
		0x54: {cpu._NI, "Not Impl", 0},
		0x55: {cpu._NI, "Not Impl", 0},
		0x56: {cpu._NI, "Not Impl", 0},
		0x57: {cpu._NI, "Not Impl", 0},
		0x58: {cpu._NI, "Not Impl", 0},
		0x59: {cpu._NI, "Not Impl", 0},
		0x5a: {cpu._NI, "Not Impl", 0},
		0x5b: {cpu._NI, "Not Impl", 0},
		0x5c: {cpu._NI, "Not Impl", 0},
		0x5d: {cpu._NI, "Not Impl", 0},
		0x5e: {cpu._NI, "Not Impl", 0},
		0x5f: {cpu._NI, "Not Impl", 0},

		0x60: {cpu._NI, "Not Impl", 0},
		0x61: {cpu._NI, "Not Impl", 0},
		0x62: {cpu._NI, "Not Impl", 0},
		0x63: {cpu._NI, "Not Impl", 0},
		0x64: {cpu._NI, "Not Impl", 0},
		0x65: {cpu._NI, "Not Impl", 0},
		0x66: {cpu._NI, "Not Impl", 0},
		0x67: {cpu._NI, "Not Impl", 0},
		0x68: {cpu._NI, "Not Impl", 0},
		0x69: {cpu._NI, "Not Impl", 0},
		0x6a: {cpu._NI, "Not Impl", 0},
		0x6b: {cpu._NI, "Not Impl", 0},
		0x6c: {cpu._NI, "Not Impl", 0},
		0x6d: {cpu._NI, "Not Impl", 0},
		0x6e: {cpu._NI, "Not Impl", 0},
		0x6f: {cpu._NI, "Not Impl", 0},

		0x70: {cpu._NI, "Not Impl", 0},
		0x71: {cpu._NI, "Not Impl", 0},
		0x72: {cpu._NI, "Not Impl", 0},
		0x73: {cpu._NI, "Not Impl", 0},
		0x74: {cpu._NI, "Not Impl", 0},
		0x75: {cpu._NI, "Not Impl", 0},
		0x76: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0x77: {cpu._NI, "Not Impl", 0},
		0x78: {cpu._NI, "Not Impl", 0},
		0x79: {cpu._NI, "Not Impl", 0},
		0x7a: {cpu._NI, "Not Impl", 0},
		0x7b: {cpu._NI, "Not Impl", 0},
		0x7c: {cpu._NI, "Not Impl", 0},
		0x7d: {cpu._NI, "Not Impl", 0},
		0x7e: {cpu._NI, "Not Impl", 0},
		0x7f: {cpu._NI, "Not Impl", 0},

		0x80: {cpu._NI, "Not Impl", 0},
		0x81: {cpu._NI, "Not Impl", 0},
		0x82: {cpu._NI, "Not Impl", 0},
		0x83: {cpu._NI, "Not Impl", 0},
		0x84: {cpu._NI, "Not Impl", 0},
		0x85: {cpu._NI, "Not Impl", 0},
		0x86: {cpu._NI, "Not Impl", 0},
		0x87: {cpu._NI, "Not Impl", 0},
		0x88: {cpu._NI, "Not Impl", 0},
		0x89: {cpu._NI, "Not Impl", 0},
		0x8a: {cpu._NI, "Not Impl", 0},
		0x8b: {cpu._NI, "Not Impl", 0},
		0x8c: {cpu._NI, "Not Impl", 0},
		0x8d: {cpu._NI, "Not Impl", 0},
		0x8e: {cpu._NI, "Not Impl", 0},
		0x8f: {cpu._NI, "Not Impl", 0},

		0x90: {cpu._NI, "Not Impl", 0},
		0x91: {cpu._NI, "Not Impl", 0},
		0x92: {cpu._NI, "Not Impl", 0},
		0x93: {cpu._NI, "Not Impl", 0},
		0x94: {cpu._NI, "Not Impl", 0},
		0x95: {cpu._NI, "Not Impl", 0},
		0x96: {cpu._NI, "Not Impl", 0},
		0x97: {cpu._NI, "Not Impl", 0},
		0x98: {cpu._NI, "Not Impl", 0},
		0x99: {cpu._NI, "Not Impl", 0},
		0x9a: {cpu._NI, "Not Impl", 0},
		0x9b: {cpu._NI, "Not Impl", 0},
		0x9c: {cpu._NI, "Not Impl", 0},
		0x9d: {cpu._NI, "Not Impl", 0},
		0x9e: {cpu._NI, "Not Impl", 0},
		0x9f: {cpu._NI, "Not Impl", 0},

		0xa0: {cpu._NI, "Not Impl", 0},
		0xa1: {cpu._NI, "Not Impl", 0},
		0xa2: {cpu._NI, "Not Impl", 0},
		0xa3: {cpu._NI, "Not Impl", 0},
		0xa4: {cpu._NI, "Not Impl", 0},
		0xa5: {cpu._NI, "Not Impl", 0},
		0xa6: {cpu._NI, "Not Impl", 0},
		0xa7: {cpu._NI, "Not Impl", 0},
		0xa8: {cpu._NI, "Not Impl", 0},
		0xa9: {cpu._NI, "Not Impl", 0},
		0xaa: {cpu._NI, "Not Impl", 0},
		0xab: {cpu._NI, "Not Impl", 0},
		0xac: {cpu._NI, "Not Impl", 0},
		0xad: {cpu._NI, "Not Impl", 0},
		0xae: {cpu._NI, "Not Impl", 0},
		0xaf: {cpu._NI, "Not Impl", 0},

		0xb0: {cpu._NI, "Not Impl", 0},
		0xb1: {cpu._NI, "Not Impl", 0},
		0xb2: {cpu._NI, "Not Impl", 0},
		0xb3: {cpu._NI, "Not Impl", 0},
		0xb4: {cpu._NI, "Not Impl", 0},
		0xb5: {cpu._NI, "Not Impl", 0},
		0xb6: {cpu._NI, "Not Impl", 0},
		0xb7: {cpu._NI, "Not Impl", 0},
		0xb8: {cpu._NI, "Not Impl", 0},
		0xb9: {cpu._NI, "Not Impl", 0},
		0xba: {cpu._NI, "Not Impl", 0},
		0xbb: {cpu._NI, "Not Impl", 0},
		0xbc: {cpu._NI, "Not Impl", 0},
		0xbd: {cpu._NI, "Not Impl", 0},
		0xbe: {cpu._NI, "Not Impl", 0},
		0xbf: {cpu._NI, "Not Impl", 0},

		0xc0: {cpu._NI, "Not Impl", 0},
		0xc1: {cpu._NI, "Not Impl", 0},
		0xc2: {cpu._JNZ, "JNZ addr", 3},
		0xc3: {cpu._JMP, "JMP addr", 3},
		0xc4: {cpu._NI, "Not Impl", 0},
		0xc5: {cpu._NI, "Not Impl", 0},
		0xc6: {cpu._NI, "Not Impl", 0},
		0xc7: {cpu._NI, "Not Impl", 0},
		0xc8: {cpu._NI, "Not Impl", 0},
		0xc9: {cpu._RET, "RET", 1},
		0xca: {cpu._NI, "Not Impl", 0},
		0xcb: {cpu._NI, "Not Impl", 0},
		0xcc: {cpu._NI, "Not Impl", 0},
		0xcd: {cpu._CALL, "CALL addr", 3},
		0xce: {cpu._NI, "Not Impl", 0},
		0xcf: {cpu._NI, "Not Impl", 0},

		0xd0: {cpu._NI, "Not Impl", 0},
		0xd1: {cpu._NI, "Not Impl", 0},
		0xd2: {cpu._NI, "Not Impl", 0},
		0xd3: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0xd4: {cpu._NI, "Not Impl", 0},
		0xd5: {cpu._NI, "Not Impl", 0},
		0xd6: {cpu._NI, "Not Impl", 0},
		0xd7: {cpu._NI, "Not Impl", 0},
		0xd8: {cpu._NI, "Not Impl", 0},
		0xd9: {cpu._NI, "Not Impl", 0},
		0xda: {cpu._NI, "Not Impl", 0},
		0xdb: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0xdc: {cpu._NI, "Not Impl", 0},
		0xdd: {cpu._NI, "Not Impl", 0},
		0xde: {cpu._NI, "Not Impl", 0},
		0xdf: {cpu._NI, "Not Impl", 0},

		0xe0: {cpu._NI, "Not Impl", 0},
		0xe1: {cpu._NI, "Not Impl", 0},
		0xe2: {cpu._NI, "Not Impl", 0},
		0xe3: {cpu._NI, "Not Impl", 0},
		0xe4: {cpu._NI, "Not Impl", 0},
		0xe5: {cpu._NI, "Not Impl", 0},
		0xe6: {cpu._ANI, "ANI", 2},
		0xe7: {cpu._NI, "Not Impl", 0},
		0xe8: {cpu._NI, "Not Impl", 0},
		0xe9: {cpu._NI, "Not Impl", 0},
		0xea: {cpu._NI, "Not Impl", 0},
		0xeb: {cpu._NI, "Not Impl", 0},
		0xec: {cpu._NI, "Not Impl", 0},
		0xed: {cpu._NI, "Not Impl", 0},
		0xee: {cpu._NI, "Not Impl", 0},
		0xef: {cpu._NI, "Not Impl", 0},

		0xf0: {cpu._NI, "Not Impl", 0},
		0xf1: {cpu._NI, "Not Impl", 0},
		0xf2: {cpu._NI, "Not Impl", 0},
		0xf3: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0xf4: {cpu._NI, "Not Impl", 0},
		0xf5: {cpu._NI, "Not Impl", 0},
		0xf6: {cpu._NI, "Not Impl", 0},
		0xf7: {cpu._NI, "Not Impl", 0},
		0xf8: {cpu._NI, "Not Impl", 0},
		0xf9: {cpu._NI, "Not Impl", 0},
		0xfa: {cpu._NI, "Not Impl", 0},
		0xfb: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0xfc: {cpu._NI, "Not Impl", 0},
		0xfd: {cpu._NI, "Not Impl", 0},
		0xfe: {cpu._CPI, "CPI", 2},
		0xff: {cpu._NI, "Not Impl", 0},
	}

	return cpu
}

func (cpu *Intel8080) GetMemory() []byte {
	return cpu.memory[:]
}

func (cpu *Intel8080) LoadProgram(program []byte) {
	copy(cpu.memory[:], program)
}

func (cpu *Intel8080) Run() {
	opcode := cpu.memory[cpu.pc]
	cpu.pc++

	instruction := cpu.instructions[opcode]
	log.Printf("%s Size=%d", instruction.mnemonic, instruction.size)
	cycles := instruction.operation()

	cpu.cycles += cycles
}

func hasParity(b byte) bool {
	return bits.OnesCount8(b)%2 == 0
}

func (cpu *Intel8080) _NI() uint {
	panic("Instruction not implemented")
}

func (cpu *Intel8080) _NOP() uint {
	// No operation
	return 4
}

func (cpu *Intel8080) _LXI_B() uint {
	cpu.c = cpu.memory[cpu.pc]
	cpu.b = cpu.memory[cpu.pc+1]
	cpu.pc += 2

	return 10
}

func (cpu *Intel8080) _STAX_B() uint {
	addr := uint16(cpu.b)<<8 | uint16(cpu.c)
	cpu.memory[addr] = cpu.a

	return 7
}
func (cpu *Intel8080) _INX_B() uint {
	bc := uint16(cpu.b)<<8 | uint16(cpu.c)
	bc++

	cpu.b = byte(bc >> 8)
	cpu.c = byte(bc & 0x00FF)

	return 5
}

func (cpu *Intel8080) _INR_B() uint {
	cpu.b++

	cpu.flags.Set(Zero, cpu.b == 0)
	cpu.flags.Set(Sign, cpu.b&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.b&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.b))

	return 5
}

func (cpu *Intel8080) _DCR_B() uint {
	cpu.b--

	cpu.flags.Set(Zero, cpu.b == 0)
	cpu.flags.Set(Sign, cpu.b&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.b&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.b))

	return 5
}

func (cpu *Intel8080) _MVI_B() uint {
	cpu.b = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _RLC() uint {
	highBit := cpu.a >> 7
	cpu.a = (cpu.a << 1) | highBit
	cpu.flags.Set(Carry, highBit == 1)

	return 4
}

func (cpu *Intel8080) _DAD_B() uint {
	hl := uint32(cpu.h)<<8 | uint32(cpu.l)
	bc := uint32(cpu.b)<<8 | uint32(cpu.c)

	result := hl + bc

	cpu.h = byte(result >> 8)
	cpu.l = byte(result & 0xFF)

	cpu.flags.Set(Carry, result > 0xFFFF)

	return 10
}

func (cpu *Intel8080) _LDAX_B() uint {
	bc := (uint16(cpu.b) << 8) | uint16(cpu.c)
	cpu.a = cpu.memory[bc]

	return 7
}

func (cpu *Intel8080) _DCX_B() uint {
	result := (uint16(cpu.b) << 8) | uint16(cpu.c)
	result--
	cpu.b, cpu.c = uint8(result>>8), uint8(result&0xFF)

	return 5
}

func (cpu *Intel8080) _INR_C() uint {
	cpu.c++

	cpu.flags.Set(Zero, cpu.c == 0)
	cpu.flags.Set(Sign, cpu.c&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.c&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.c))

	return 5
}

func (cpu *Intel8080) _DCR_C() uint {
	cpu.c--

	cpu.flags.Set(Zero, cpu.c == 0)
	cpu.flags.Set(Sign, cpu.c&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.c&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.c))

	return 5
}

func (cpu *Intel8080) _MVI_C() uint {
	cpu.c = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _RRC() uint {
	lb := cpu.a & 0x1
	cpu.flags.Set(Carry, lb == 1)
	cpu.a = (cpu.a >> 1) | (lb << 7)

	return 4
}

func (cpu *Intel8080) _LXI_D() uint {
	cpu.e = cpu.memory[cpu.pc]
	cpu.d = cpu.memory[cpu.pc+1]
	cpu.pc += 2

	return 10
}

func (cpu *Intel8080) _STAX_D() uint {
	addr := uint16(cpu.d)<<8 | uint16(cpu.e)
	cpu.memory[addr] = cpu.a

	return 7
}

func (cpu *Intel8080) _INX_D() uint {
	de := uint16(cpu.d)<<8 | uint16(cpu.e)
	de++

	cpu.d = byte(de >> 8)
	cpu.e = byte(de & 0x00FF)

	return 5
}

func (cpu *Intel8080) _INR_D() uint {
	cpu.d++

	cpu.flags.Set(Zero, cpu.d == 0)
	cpu.flags.Set(Sign, cpu.d&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.d&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.d))

	return 5
}

func (cpu *Intel8080) _DCR_D() uint {
	cpu.d--

	cpu.flags.Set(Zero, cpu.d == 0)
	cpu.flags.Set(Sign, cpu.d&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.d&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.d))

	return 5
}

func (cpu *Intel8080) _MVI_D() uint {
	cpu.d = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _RAL() uint {
	var c uint8 = 0
	if cpu.flags.Get(Carry) {
		c = 1
	}

	hb := cpu.a >> 7
	cpu.flags.Set(Carry, hb == 1)
	cpu.a = (cpu.a << 1) | c

	return 4
}

func (cpu *Intel8080) _DAD_D() uint {
	hl := uint32(cpu.h)<<8 | uint32(cpu.l)
	de := uint32(cpu.d)<<8 | uint32(cpu.e)

	result := hl + de

	cpu.h = byte(result >> 8)
	cpu.l = byte(result & 0xFF)

	cpu.flags.Set(Carry, result > 0xFFFF)

	return 10
}

func (cpu *Intel8080) _LDAX_D() uint {
	de := uint16(cpu.d)<<8 | uint16(cpu.e)
	cpu.a = cpu.memory[de]

	return 7
}

func (cpu *Intel8080) _DCX_D() uint {
	result := uint16(cpu.d)<<8 | uint16(cpu.e)
	result--
	cpu.d, cpu.e = uint8(result>>8), uint8(result&0xFF)

	return 5
}

func (cpu *Intel8080) _INR_E() uint {
	cpu.e++

	cpu.flags.Set(Zero, cpu.e == 0)
	cpu.flags.Set(Sign, cpu.e&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.e&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.e))

	return 5
}

func (cpu *Intel8080) _DCR_E() uint {
	cpu.e--

	cpu.flags.Set(Zero, cpu.e == 0)
	cpu.flags.Set(Sign, cpu.e&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.e&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.e))

	return 5
}

func (cpu *Intel8080) _MVI_E() uint {
	cpu.e = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _RAR() uint {
	var c uint8 = 0
	if cpu.flags.Get(Carry) {
		c = 1
	}

	lb := cpu.a & 0x1
	cpu.flags.Set(Carry, lb == 1)
	cpu.a = (cpu.a >> 1) | (c << 7)

	return 4
}

func (cpu *Intel8080) _LXI_H() uint {
	cpu.l = cpu.memory[cpu.pc]
	cpu.h = cpu.memory[cpu.pc+1]
	cpu.pc += 2

	return 10
}

func (cpu *Intel8080) _SHLD() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	addr := (hb << 8) | lb
	cpu.memory[addr] = cpu.l
	cpu.memory[addr+1] = cpu.h

	cpu.pc += 2

	return 16
}

func (cpu *Intel8080) _INX_H() uint {
	hl := uint16(cpu.h)<<8 | uint16(cpu.l)
	hl++

	cpu.h = byte(hl >> 8)
	cpu.l = byte(hl & 0x00FF)

	return 5
}

func (cpu *Intel8080) _INR_H() uint {
	cpu.h++

	cpu.flags.Set(Zero, cpu.h == 0)
	cpu.flags.Set(Sign, cpu.h&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.h&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.h))

	return 5
}

func (cpu *Intel8080) _DCR_H() uint {
	cpu.h--

	cpu.flags.Set(Zero, cpu.h == 0)
	cpu.flags.Set(Sign, cpu.h&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.h&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.h))

	return 5
}

func (cpu *Intel8080) _MVI_H() uint {
	cpu.h = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _DAD_H() uint {
	hl := uint32(cpu.h)<<8 | uint32(cpu.l)

	result := hl + hl

	cpu.h = byte(result >> 8)
	cpu.l = byte(result & 0xFF)

	cpu.flags.Set(Carry, result > 0xFFFF)

	return 10
}

func (cpu *Intel8080) _LHLD() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	addr := (hb << 8) | lb

	cpu.l = cpu.memory[addr]
	cpu.h = cpu.memory[addr+1]

	cpu.pc += 2

	return 16
}

func (cpu *Intel8080) _DCX_H() uint {
	result := uint16(cpu.h)<<8 | uint16(cpu.l)
	result--
	cpu.h, cpu.l = uint8(result>>8), uint8(result&0xFF)

	return 5
}

func (cpu *Intel8080) _INR_L() uint {
	cpu.l++

	cpu.flags.Set(Zero, cpu.l == 0)
	cpu.flags.Set(Sign, cpu.l&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.l&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.l))

	return 5
}

func (cpu *Intel8080) _DCR_L() uint {
	cpu.l--

	cpu.flags.Set(Zero, cpu.l == 0)
	cpu.flags.Set(Sign, cpu.l&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.l&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.l))

	return 5
}

func (cpu *Intel8080) _MVI_L() uint {
	cpu.l = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _CMA() uint {
	cpu.a = ^cpu.a

	return 4
}

func (cpu *Intel8080) _LXI_SP() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	cpu.sp = (hb << 8) | lb

	cpu.pc += 2

	return 10
}

func (cpu *Intel8080) _STA() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	addr := (hb << 8) | lb
	cpu.memory[addr] = cpu.a

	cpu.pc += 2

	return 13
}

func (cpu *Intel8080) _INX_SP() uint {
	cpu.sp++

	return 5
}

func (cpu *Intel8080) _INR_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	value := cpu.memory[addr]
	value++
	cpu.memory[addr] = value

	cpu.flags.Set(Zero, value == 0)
	cpu.flags.Set(Sign, value&0x80 != 0)
	cpu.flags.Set(AuxCarry, value&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(value))

	return 10
}

func (cpu *Intel8080) _DCR_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	value := cpu.memory[addr]
	value--
	cpu.memory[addr] = value

	cpu.flags.Set(Zero, value == 0)
	cpu.flags.Set(Sign, value&0x80 != 0)
	cpu.flags.Set(AuxCarry, value&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(value))

	return 10
}

func (cpu *Intel8080) _MVI_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	value := cpu.memory[cpu.pc]
	cpu.memory[addr] = value

	cpu.pc++

	return 10
}

func (cpu *Intel8080) _STC() uint {
	cpu.flags.Set(Carry, true)

	return 4
}

func (cpu *Intel8080) _DAD_SP() uint {
	hl := uint32(cpu.h)<<8 | uint32(cpu.l)

	result := hl + uint32(cpu.sp)

	cpu.h = byte(result >> 8)
	cpu.l = byte(result & 0xFF)

	cpu.flags.Set(Carry, result > 0xFFFF)

	return 10
}

func (cpu *Intel8080) _LDA() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	addr := hb<<8 | lb
	cpu.a = cpu.memory[addr]

	cpu.pc += 2

	return 13
}

func (cpu *Intel8080) _DCX_SP() uint {
	cpu.sp--

	return 5
}

func (cpu *Intel8080) _INR_A() uint {
	cpu.a++

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.a&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))

	return 5
}

func (cpu *Intel8080) _DCR_A() uint {
	cpu.a--

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.a&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))

	return 5
}

func (cpu *Intel8080) _MVI_A() uint {
	cpu.a = cpu.memory[cpu.pc]
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _CMC() uint {
	if cpu.flags.Get(Carry) {
		cpu.flags.Set(Carry, false)
	} else {
		cpu.flags.Set(Carry, true)
	}

	return 4
}

func (cpu *Intel8080) _JNZ() uint {
	if cpu.flags.Get(Zero) {
		lb := uint16(cpu.memory[cpu.pc])
		hb := uint16(cpu.memory[cpu.pc+1])

		cpu.pc = (hb << 8) | lb
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _JMP() uint {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	cpu.pc = (hb << 8) | lb

	return 10
}

func (cpu *Intel8080) _RET() uint {
	lb, hb := uint16(cpu.memory[cpu.sp]), uint16(cpu.memory[cpu.sp+1])
	cpu.sp += 2
	cpu.pc = (hb << 8) | lb

	return 10
}

func (cpu *Intel8080) _CALL() uint {
	ret := cpu.pc + 2
	cpu.memory[cpu.sp-1] = uint8((ret >> 8) & 0xff)
	cpu.memory[cpu.sp-2] = uint8(ret & 0xff)
	cpu.sp -= 2

	lb, hb := uint16(cpu.memory[cpu.pc]), uint16(cpu.memory[cpu.pc+1])
	cpu.pc = (hb << 8) | lb

	return 17
}

func (cpu *Intel8080) _ANI() uint {
	x := cpu.a & cpu.memory[cpu.pc]

	cpu.flags.Set(Carry, false)
	cpu.flags.Set(AuxCarry, false)

	cpu.flags.Set(Zero, x == 0)
	cpu.flags.Set(Sign, x&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(x))

	cpu.a = x
	cpu.pc++

	return 7
}

func (cpu *Intel8080) _CPI() uint {
	value := cpu.memory[cpu.pc]
	result := cpu.a - value

	cpu.flags.Set(Zero, result == 0)
	cpu.flags.Set(Sign, result&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(result))
	cpu.flags.Set(Carry, cpu.a < value)
	cpu.flags.Set(AuxCarry, (cpu.a&0xf) < (value&0xf))

	cpu.pc++

	return 7
}
