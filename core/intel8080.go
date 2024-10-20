package core

import (
	"math/bits"

	"github.com/gaoliveira21/intel8080-space-invaders/debug"
)

type Intel8080Instruction struct {
	operation func()
	mnemonic  string
	cycles    uint16
}

type Intel8080 struct {
	debugger *debug.Debugger

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

func NewIntel8080(d *debug.Debugger) *Intel8080 {
	cpu := &Intel8080{
		debugger: d,
		flags:    &intel8080Flags{},
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
		0x12: {cpu._NI, "Not Impl", 0},
		0x13: {cpu._NI, "Not Impl", 0},
		0x14: {cpu._NI, "Not Impl", 0},
		0x15: {cpu._NI, "Not Impl", 0},
		0x16: {cpu._NI, "Not Impl", 0},
		0x17: {cpu._NI, "Not Impl", 0},
		0x18: {cpu._NI, "Not Impl", 0},
		0x19: {cpu._NI, "Not Impl", 0},
		0x1a: {cpu._NI, "Not Impl", 0},
		0x1b: {cpu._NI, "Not Impl", 0},
		0x1c: {cpu._NI, "Not Impl", 0},
		0x1d: {cpu._NI, "Not Impl", 0},
		0x1e: {cpu._NI, "Not Impl", 0},
		0x1f: {cpu._RAR, "RAR", 1},

		0x20: {cpu._NI, "Not Impl", 0},
		0x21: {cpu._NI, "Not Impl", 0},
		0x22: {cpu._NI, "Not Impl", 0},
		0x23: {cpu._NI, "Not Impl", 0},
		0x24: {cpu._NI, "Not Impl", 0},
		0x25: {cpu._NI, "Not Impl", 0},
		0x26: {cpu._NI, "Not Impl", 0},
		0x27: {cpu._NI, "Not Impl", 0},
		0x28: {cpu._NI, "Not Impl", 0},
		0x29: {cpu._NI, "Not Impl", 0},
		0x2a: {cpu._NI, "Not Impl", 0},
		0x2b: {cpu._NI, "Not Impl", 0},
		0x2c: {cpu._NI, "Not Impl", 0},
		0x2d: {cpu._NI, "Not Impl", 0},
		0x2e: {cpu._NI, "Not Impl", 0},
		0x2f: {cpu._CMA, "CMA", 1},

		0x30: {cpu._NI, "Not Impl", 0},
		0x31: {cpu._NI, "Not Impl", 0},
		0x32: {cpu._NI, "Not Impl", 0},
		0x33: {cpu._NI, "Not Impl", 0},
		0x34: {cpu._NI, "Not Impl", 0},
		0x35: {cpu._NI, "Not Impl", 0},
		0x36: {cpu._NI, "Not Impl", 0},
		0x37: {cpu._NI, "Not Impl", 0},
		0x38: {cpu._NI, "Not Impl", 0},
		0x39: {cpu._NI, "Not Impl", 0},
		0x3a: {cpu._NI, "Not Impl", 0},
		0x3b: {cpu._NI, "Not Impl", 0},
		0x3c: {cpu._NI, "Not Impl", 0},
		0x3d: {cpu._NI, "Not Impl", 0},
		0x3e: {cpu._NI, "Not Impl", 0},
		0x3f: {cpu._NI, "Not Impl", 0},

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
		0x76: {cpu._NI, "Not Impl", 0},
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
		0xd3: {cpu._NI, "Not Impl", 0},
		0xd4: {cpu._NI, "Not Impl", 0},
		0xd5: {cpu._NI, "Not Impl", 0},
		0xd6: {cpu._NI, "Not Impl", 0},
		0xd7: {cpu._NI, "Not Impl", 0},
		0xd8: {cpu._NI, "Not Impl", 0},
		0xd9: {cpu._NI, "Not Impl", 0},
		0xda: {cpu._NI, "Not Impl", 0},
		0xdb: {cpu._NI, "Not Impl", 0},
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
		0xf3: {cpu._NI, "Not Impl", 0},
		0xf4: {cpu._NI, "Not Impl", 0},
		0xf5: {cpu._NI, "Not Impl", 0},
		0xf6: {cpu._NI, "Not Impl", 0},
		0xf7: {cpu._NI, "Not Impl", 0},
		0xf8: {cpu._NI, "Not Impl", 0},
		0xf9: {cpu._NI, "Not Impl", 0},
		0xfa: {cpu._NI, "Not Impl", 0},
		0xfb: {cpu._NI, "Not Impl", 0},
		0xfc: {cpu._NI, "Not Impl", 0},
		0xfd: {cpu._NI, "Not Impl", 0},
		0xfe: {cpu._CPI, "CPI", 2},
		0xff: {cpu._NI, "Not Impl", 0},
	}

	return cpu
}

func (cpu *Intel8080) LoadProgram(program []byte) {
	copy(cpu.memory[:], program)

	if cpu.debugger != nil {
		cpu.debugger.DumpMemory(cpu.memory[:])
	}
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

func (cpu *Intel8080) _NI() {
	panic("Instruction not implemented")
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

func (cpu *Intel8080) _RLC() {
	highBit := cpu.a >> 7
	cpu.a = (cpu.a << 1) | highBit
	cpu.flags.Set(Carry, highBit == 1)
}

func (cpu *Intel8080) _DAD_B() {
	hl := uint32(cpu.h)<<8 | uint32(cpu.l)
	bc := uint32(cpu.b)<<8 | uint32(cpu.c)

	result := hl + bc

	cpu.h = byte(result >> 8)
	cpu.l = byte(result & 0xFF)

	cpu.flags.Set(Carry, result > 0xFFFF)
}

func (cpu *Intel8080) _LDAX_B() {
	bc := (uint16(cpu.b) << 8) | uint16(cpu.c)
	cpu.a = cpu.memory[bc]
}

func (cpu *Intel8080) _DCX_B() {
	result := (uint16(cpu.b) << 8) | uint16(cpu.c)
	result--
	cpu.b, cpu.c = uint8(result>>8), uint8(result&0xFF)
}

func (cpu *Intel8080) _INR_C() {
	cpu.c++

	cpu.flags.Set(Zero, cpu.c == 0)
	cpu.flags.Set(Sign, cpu.c&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.c&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.c))
}

func (cpu *Intel8080) _DCR_C() {
	cpu.c--

	cpu.flags.Set(Zero, cpu.c == 0)
	cpu.flags.Set(Sign, cpu.c&0x80 != 0)
	cpu.flags.Set(AuxCarry, cpu.c&0x0F == 0)
	cpu.flags.Set(Parity, hasParity(cpu.c))
}

func (cpu *Intel8080) _MVI_C() {
	cpu.c = cpu.memory[cpu.pc]
	cpu.pc++
}

func (cpu *Intel8080) _RRC() {
	lb := cpu.a & 0x1
	cpu.flags.Set(Carry, lb == 1)
	cpu.a = (cpu.a >> 1) | (lb << 7)
}

func (cpu *Intel8080) _LXI_D() {
	cpu.e = cpu.memory[cpu.pc]
	cpu.d = cpu.memory[cpu.pc+1]
	cpu.pc += 2
}

func (cpu *Intel8080) _RAR() {
	var c uint8 = 0
	if cpu.flags.Get(Carry) {
		c = 1
	}

	lb := cpu.a & 0x1
	cpu.flags.Set(Carry, lb == 1)
	cpu.a = (cpu.a >> 1) | (c << 7)
}

func (cpu *Intel8080) _CMA() {
	cpu.a = ^cpu.a
}

func (cpu *Intel8080) _JNZ() {
	if cpu.flags.Get(Zero) {
		lb := uint16(cpu.memory[cpu.pc])
		hb := uint16(cpu.memory[cpu.pc+1])

		cpu.pc = (hb << 8) | lb
	} else {
		cpu.pc += 2
	}
}

func (cpu *Intel8080) _JMP() {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) _CALL() {
	ret := cpu.pc + 2
	cpu.memory[cpu.sp-1] = uint8((ret >> 8) & 0xff)
	cpu.memory[cpu.sp-2] = uint8(ret & 0xff)
	cpu.sp -= 2

	lb, hb := uint16(cpu.memory[cpu.pc]), uint16(cpu.memory[cpu.pc+1])
	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) _RET() {
	lb, hb := uint16(cpu.memory[cpu.sp]), uint16(cpu.memory[cpu.sp+1])
	cpu.sp += 2
	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) _ANI() {
	x := cpu.a & cpu.memory[cpu.pc]

	cpu.flags.Set(Carry, false)
	cpu.flags.Set(AuxCarry, false)

	cpu.flags.Set(Zero, x == 0)
	cpu.flags.Set(Sign, x&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(x))

	cpu.a = x
	cpu.pc++
}

func (cpu *Intel8080) _CPI() {
	value := cpu.memory[cpu.pc]
	result := cpu.a - value

	cpu.flags.Set(Zero, result == 0)
	cpu.flags.Set(Sign, result&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(result))
	cpu.flags.Set(Carry, cpu.a < value)
	cpu.flags.Set(AuxCarry, (cpu.a&0xf) < (value&0xf))

	cpu.pc++
}
