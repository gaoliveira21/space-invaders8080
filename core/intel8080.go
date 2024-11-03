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

		0x40: {cpu._MOV_BB, "MOV B,B", 1},
		0x41: {cpu._MOV_BC, "MOV B,C", 1},
		0x42: {cpu._MOV_BD, "MOV B,D", 1},
		0x43: {cpu._MOV_BE, "MOV B,E", 1},
		0x44: {cpu._MOV_BH, "MOV B,H", 1},
		0x45: {cpu._MOV_BL, "MOV B,L", 1},
		0x46: {cpu._MOV_BM, "MOV B,M", 1},
		0x47: {cpu._MOV_BA, "MOV B,A", 1},
		0x48: {cpu._MOV_CB, "MOV C,B", 1},
		0x49: {cpu._MOV_CC, "MOV C,C", 1},
		0x4a: {cpu._MOV_CD, "MOV C,D", 1},
		0x4b: {cpu._MOV_CE, "MOV C,E", 1},
		0x4c: {cpu._MOV_CH, "MOV C,H", 1},
		0x4d: {cpu._MOV_CL, "MOV C,L", 1},
		0x4e: {cpu._MOV_CM, "MOV C,M", 1},
		0x4f: {cpu._MOV_CA, "MOV C,A", 1},

		0x50: {cpu._MOV_DB, "MOV D,B", 1},
		0x51: {cpu._MOV_DC, "MOV D,C", 1},
		0x52: {cpu._MOV_DD, "MOV D,D", 1},
		0x53: {cpu._MOV_DE, "MOV D,E", 1},
		0x54: {cpu._MOV_DH, "MOV D,H", 1},
		0x55: {cpu._MOV_DL, "MOV D,L", 1},
		0x56: {cpu._MOV_DM, "MOV D,M", 1},
		0x57: {cpu._MOV_DA, "MOV D,A", 1},
		0x58: {cpu._MOV_EB, "MOV E,B", 1},
		0x59: {cpu._MOV_EC, "MOV E,C", 1},
		0x5a: {cpu._MOV_ED, "MOV E,D", 1},
		0x5b: {cpu._MOV_EE, "MOV E,E", 1},
		0x5c: {cpu._MOV_EH, "MOV E,H", 1},
		0x5d: {cpu._MOV_EL, "MOV E,L", 1},
		0x5e: {cpu._MOV_EM, "MOV E,M", 1},
		0x5f: {cpu._MOV_EA, "MOV E,A", 1},

		0x60: {cpu._MOV_HB, "MOV H,B", 1},
		0x61: {cpu._MOV_HC, "MOV H,C", 1},
		0x62: {cpu._MOV_HD, "MOV H,D", 1},
		0x63: {cpu._MOV_HE, "MOV H,E", 1},
		0x64: {cpu._MOV_HH, "MOV H,H", 1},
		0x65: {cpu._MOV_HL, "MOV H,L", 1},
		0x66: {cpu._MOV_HM, "MOV H,M", 1},
		0x67: {cpu._MOV_HA, "MOV H,A", 1},
		0x68: {cpu._MOV_LB, "MOV L,B", 1},
		0x69: {cpu._MOV_LC, "MOV L,C", 1},
		0x6a: {cpu._MOV_LD, "MOV L,D", 1},
		0x6b: {cpu._MOV_LE, "MOV L,E", 1},
		0x6c: {cpu._MOV_LH, "MOV L,H", 1},
		0x6d: {cpu._MOV_LL, "MOV L,L", 1},
		0x6e: {cpu._MOV_LM, "MOV L,M", 1},
		0x6f: {cpu._MOV_LA, "MOV L,A", 1},

		0x70: {cpu._MOV_MB, "MOV M,B", 1},
		0x71: {cpu._MOV_MC, "MOV M,C", 1},
		0x72: {cpu._MOV_MD, "MOV M,D", 1},
		0x73: {cpu._MOV_ME, "MOV M,E", 1},
		0x74: {cpu._MOV_MH, "MOV M,H", 1},
		0x75: {cpu._MOV_ML, "MOV M,L", 1},
		0x76: {cpu._NI, "Not Impl", 0}, // IO and Special Group
		0x77: {cpu._MOV_MA, "MOV M,A", 1},
		0x78: {cpu._MOV_AB, "MOV A,B", 1},
		0x79: {cpu._MOV_AC, "MOV A,C", 1},
		0x7a: {cpu._MOV_AD, "MOV A,D", 1},
		0x7b: {cpu._MOV_AE, "MOV A,E", 1},
		0x7c: {cpu._MOV_AH, "MOV A,H", 1},
		0x7d: {cpu._MOV_AL, "MOV A,L", 1},
		0x7e: {cpu._MOV_AM, "MOV A,M", 1},
		0x7f: {cpu._MOV_AA, "MOV A,A", 1},

		0x80: {cpu._ADD_B, "ADD B", 1},
		0x81: {cpu._ADD_C, "ADD C", 1},
		0x82: {cpu._ADD_D, "ADD D", 1},
		0x83: {cpu._ADD_E, "ADD E", 1},
		0x84: {cpu._ADD_H, "ADD H", 1},
		0x85: {cpu._ADD_L, "ADD L", 1},
		0x86: {cpu._ADD_M, "ADD M", 1},
		0x87: {cpu._ADD_A, "ADD A", 1},
		0x88: {cpu._ADC_B, "ADC B", 1},
		0x89: {cpu._ADC_C, "ADC C", 1},
		0x8a: {cpu._ADC_D, "ADC D", 1},
		0x8b: {cpu._ADC_E, "ADC E", 1},
		0x8c: {cpu._ADC_H, "ADC H", 1},
		0x8d: {cpu._ADC_L, "ADC L", 1},
		0x8e: {cpu._ADC_M, "ADC M", 1},
		0x8f: {cpu._ADC_A, "ADC A", 1},

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

func (cpu *Intel8080) add(value byte) {
	result := uint16(cpu.a) + uint16(value)

	cpu.flags.Set(Carry, result > 0xFF)
	cpu.flags.Set(AuxCarry, ((cpu.a^uint8(result)^value)&0x10) > 0)

	cpu.a = uint8(result & 0xFF)

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
}

func (cpu *Intel8080) adc(value byte) {
	carryVal := uint16(0)
	if cpu.flags.Get(Carry) {
		carryVal = 1
	}

	result := uint16(cpu.a) + uint16(value) + carryVal

	cpu.flags.Set(Carry, result > 0xFF)
	cpu.flags.Set(AuxCarry, ((cpu.a^uint8(result)^value)&0x10) > 0)

	cpu.a = uint8(result & 0xFF)

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
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

func (cpu *Intel8080) _MOV_BB() uint {
	// B <- B (no operation needed since we're moving B to itself)
	return 5
}

func (cpu *Intel8080) _MOV_BC() uint {
	cpu.b = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_BD() uint {
	cpu.b = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_BE() uint {
	cpu.b = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_BH() uint {
	cpu.b = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_BL() uint {
	cpu.b = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_BM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.b = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_BA() uint {
	cpu.b = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_CB() uint {
	cpu.c = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_CC() uint {
	// C <- C (no operation needed since we're moving C to itself)
	return 5
}

func (cpu *Intel8080) _MOV_CD() uint {
	cpu.c = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_CE() uint {
	cpu.c = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_CH() uint {
	cpu.c = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_CL() uint {
	cpu.c = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_CM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.c = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_CA() uint {
	cpu.c = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_DB() uint {
	cpu.d = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_DC() uint {
	cpu.d = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_DD() uint {
	// D <- D (no operation needed since we're moving D to itself)
	return 5
}

func (cpu *Intel8080) _MOV_DE() uint {
	cpu.d = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_DH() uint {
	cpu.d = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_DL() uint {
	cpu.d = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_DM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.d = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_DA() uint {
	cpu.d = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_EB() uint {
	cpu.e = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_EC() uint {
	cpu.e = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_ED() uint {
	cpu.e = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_EE() uint {
	// E <- E (no operation needed since we're moving E to itself)
	return 5
}

func (cpu *Intel8080) _MOV_EH() uint {
	cpu.e = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_EL() uint {
	cpu.e = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_EM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.e = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_EA() uint {
	cpu.e = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_HB() uint {
	cpu.h = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_HC() uint {
	cpu.h = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_HD() uint {
	cpu.h = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_HE() uint {
	cpu.h = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_HH() uint {
	// H <- H (no operation needed since we're moving H to itself)
	return 5
}

func (cpu *Intel8080) _MOV_HL() uint {
	cpu.h = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_HM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.h = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_HA() uint {
	cpu.h = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_LB() uint {
	cpu.l = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_LC() uint {
	cpu.l = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_LD() uint {
	cpu.l = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_LE() uint {
	cpu.l = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_LH() uint {
	cpu.l = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_LL() uint {
	// L <- L (no operation needed since we're moving L to itself)
	return 5
}

func (cpu *Intel8080) _MOV_LM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.l = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_LA() uint {
	cpu.l = cpu.a
	return 5
}

func (cpu *Intel8080) _MOV_MB() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.b
	return 7
}

func (cpu *Intel8080) _MOV_MC() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.c
	return 7
}

func (cpu *Intel8080) _MOV_MD() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.d
	return 7
}

func (cpu *Intel8080) _MOV_ME() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.e
	return 7
}

func (cpu *Intel8080) _MOV_MH() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.h
	return 7
}

func (cpu *Intel8080) _MOV_ML() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.l
	return 7
}

func (cpu *Intel8080) _MOV_MA() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.memory[addr] = cpu.a
	return 7
}

func (cpu *Intel8080) _MOV_AB() uint {
	cpu.a = cpu.b
	return 5
}

func (cpu *Intel8080) _MOV_AC() uint {
	cpu.a = cpu.c
	return 5
}

func (cpu *Intel8080) _MOV_AD() uint {
	cpu.a = cpu.d
	return 5
}

func (cpu *Intel8080) _MOV_AE() uint {
	cpu.a = cpu.e
	return 5
}

func (cpu *Intel8080) _MOV_AH() uint {
	cpu.a = cpu.h
	return 5
}

func (cpu *Intel8080) _MOV_AL() uint {
	cpu.a = cpu.l
	return 5
}

func (cpu *Intel8080) _MOV_AM() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.a = cpu.memory[addr]
	return 7
}

func (cpu *Intel8080) _MOV_AA() uint {
	// A <- A (no operation needed since we're moving A to itself)
	return 5
}

func (cpu *Intel8080) _ADD_B() uint {
	cpu.add(cpu.b)
	return 4
}

func (cpu *Intel8080) _ADD_C() uint {
	cpu.add(cpu.c)
	return 4
}

func (cpu *Intel8080) _ADD_D() uint {
	cpu.add(cpu.d)
	return 4
}

func (cpu *Intel8080) _ADD_E() uint {
	cpu.add(cpu.e)
	return 4
}

func (cpu *Intel8080) _ADD_H() uint {
	cpu.add(cpu.h)
	return 4
}

func (cpu *Intel8080) _ADD_L() uint {
	cpu.add(cpu.l)
	return 4
}

func (cpu *Intel8080) _ADD_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.add(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _ADD_A() uint {
	cpu.add(cpu.a)
	return 4
}

func (cpu *Intel8080) _ADC_B() uint {
	cpu.adc(cpu.b)
	return 4
}

func (cpu *Intel8080) _ADC_C() uint {
	cpu.adc(cpu.c)
	return 4
}

func (cpu *Intel8080) _ADC_D() uint {
	cpu.adc(cpu.d)
	return 4
}

func (cpu *Intel8080) _ADC_E() uint {
	cpu.adc(cpu.e)
	return 4
}

func (cpu *Intel8080) _ADC_H() uint {
	cpu.adc(cpu.h)
	return 4
}

func (cpu *Intel8080) _ADC_L() uint {
	cpu.adc(cpu.l)
	return 4
}

func (cpu *Intel8080) _ADC_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.adc(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _ADC_A() uint {
	cpu.adc(cpu.a)
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
