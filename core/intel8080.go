package core

import (
	"math/bits"
)

type Intel8080Instruction struct {
	operation func() uint
	mnemonic  string
	size      uint16
}

type Intel8080Registers struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte
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

	InterruptEnabled        bool
	enableInterruptDeferred bool

	// listeners
	onInput  func(cpu *Intel8080)
	onOutput func(cpu *Intel8080)
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
		0x27: {cpu._DAA, "DAA", 1},
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
		0x76: {cpu._HLT, "HLT", 1},
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

		0x90: {cpu._SUB_B, "SUB B", 1},
		0x91: {cpu._SUB_C, "SUB C", 1},
		0x92: {cpu._SUB_D, "SUB D", 1},
		0x93: {cpu._SUB_E, "SUB E", 1},
		0x94: {cpu._SUB_H, "SUB H", 1},
		0x95: {cpu._SUB_L, "SUB L", 1},
		0x96: {cpu._SUB_M, "SUB M", 1},
		0x97: {cpu._SUB_A, "SUB A", 1},
		0x98: {cpu._SBB_B, "SBB B", 1},
		0x99: {cpu._SBB_C, "SBB C", 1},
		0x9a: {cpu._SBB_D, "SBB D", 1},
		0x9b: {cpu._SBB_E, "SBB E", 1},
		0x9c: {cpu._SBB_H, "SBB H", 1},
		0x9d: {cpu._SBB_L, "SBB L", 1},
		0x9e: {cpu._SBB_M, "SBB M", 1},
		0x9f: {cpu._SBB_A, "SBB A", 1},

		0xa0: {cpu._ANA_B, "ANA B", 1},
		0xa1: {cpu._ANA_C, "ANA C", 1},
		0xa2: {cpu._ANA_D, "ANA D", 1},
		0xa3: {cpu._ANA_E, "ANA E", 1},
		0xa4: {cpu._ANA_H, "ANA H", 1},
		0xa5: {cpu._ANA_L, "ANA L", 1},
		0xa6: {cpu._ANA_M, "ANA M", 1},
		0xa7: {cpu._ANA_A, "ANA A", 1},
		0xa8: {cpu._XRA_B, "XRA B", 1},
		0xa9: {cpu._XRA_C, "XRA C", 1},
		0xaa: {cpu._XRA_D, "XRA D", 1},
		0xab: {cpu._XRA_E, "XRA E", 1},
		0xac: {cpu._XRA_H, "XRA H", 1},
		0xad: {cpu._XRA_L, "XRA L", 1},
		0xae: {cpu._XRA_M, "XRA M", 1},
		0xaf: {cpu._XRA_A, "XRA A", 1},

		0xb0: {cpu._ORA_B, "ORA B", 1},
		0xb1: {cpu._ORA_C, "ORA C", 1},
		0xb2: {cpu._ORA_D, "ORA D", 1},
		0xb3: {cpu._ORA_E, "ORA E", 1},
		0xb4: {cpu._ORA_H, "ORA H", 1},
		0xb5: {cpu._ORA_L, "ORA L", 1},
		0xb6: {cpu._ORA_M, "ORA M", 1},
		0xb7: {cpu._ORA_A, "ORA A", 1},
		0xb8: {cpu._CMP_B, "CMP B", 1},
		0xb9: {cpu._CMP_C, "CMP C", 1},
		0xba: {cpu._CMP_D, "CMP D", 1},
		0xbb: {cpu._CMP_E, "CMP E", 1},
		0xbc: {cpu._CMP_H, "CMP H", 1},
		0xbd: {cpu._CMP_L, "CMP L", 1},
		0xbe: {cpu._CMP_M, "CMP M", 1},
		0xbf: {cpu._CMP_A, "CMP A", 1},

		0xc0: {cpu._RNZ, "RNZ", 1},
		0xc1: {cpu._POP_B, "POP B", 1},
		0xc2: {cpu._JNZ, "JNZ addr", 3},
		0xc3: {cpu._JMP, "JMP addr", 3},
		0xc4: {cpu._CNZ, "CNZ", 3},
		0xc5: {cpu._PUSH_B, "PUSH B", 1},
		0xc6: {cpu._ADI, "ADI", 2},
		0xc7: {cpu._RST_0, "RST 0", 1},
		0xc8: {cpu._RZ, "RZ", 1},
		0xc9: {cpu._RET, "RET", 1},
		0xca: {cpu._JZ, "JZ", 1},
		0xcb: {cpu._JMP, "*JMP", 3},
		0xcc: {cpu._CZ, "CZ", 3},
		0xcd: {cpu._CALL, "CALL addr", 3},
		0xce: {cpu._ACI, "ACI", 2},
		0xcf: {cpu._RST_1, "RST 1", 1},

		0xd0: {cpu._RNC, "RNC", 1},
		0xd1: {cpu._POP_D, "POP D", 1},
		0xd2: {cpu._JNC, "JNC", 3},
		0xd3: {cpu._OUT, "OUT", 2},
		0xd4: {cpu._CNC, "CNC", 3},
		0xd5: {cpu._PUSH_D, "PUSH D", 1},
		0xd6: {cpu._SUI, "SUI", 2},
		0xd7: {cpu._RST_2, "RST 2", 1},
		0xd8: {cpu._RC, "RC", 1},
		0xd9: {cpu._RET, "*RET", 1},
		0xda: {cpu._JC, "JC", 3},
		0xdb: {cpu._IN, "IN", 2},
		0xdc: {cpu._CC, "CC", 3},
		0xdd: {cpu._CALL, "*CALL", 3},
		0xde: {cpu._SBI, "SBI", 2},
		0xdf: {cpu._RST_3, "RST 3", 1},

		0xe0: {cpu._RPO, "RPO", 1},
		0xe1: {cpu._POP_H, "POP H", 1},
		0xe2: {cpu._JPO, "JPO", 3},
		0xe3: {cpu._XTHL, "XTHL", 1},
		0xe4: {cpu._CPO, "CPO", 3},
		0xe5: {cpu._PUSH_H, "PUSH H", 1},
		0xe6: {cpu._ANI, "ANI", 2},
		0xe7: {cpu._RST_4, "RST 4", 1},
		0xe8: {cpu._RPE, "RPE", 1},
		0xe9: {cpu._PCHL, "PCHL", 1},
		0xea: {cpu._JPE, "JPE", 3},
		0xeb: {cpu._XCHG, "XCHG", 1},
		0xec: {cpu._CPE, "CPE", 3},
		0xed: {cpu._CALL, "*CALL", 3},
		0xee: {cpu._XRI, "XRI", 2},
		0xef: {cpu._RST_5, "RST 5", 1},

		0xf0: {cpu._RP, "RP", 1},
		0xf1: {cpu._POP_PSW, "POP PSW", 1},
		0xf2: {cpu._JP, "JP", 3},
		0xf3: {cpu._DI, "DI", 1},
		0xf4: {cpu._CP, "CP", 3},
		0xf5: {cpu._PUSH_PSW, "PUSH PSW", 1},
		0xf6: {cpu._ORI, "ORI", 2},
		0xf7: {cpu._RST_6, "RST 6", 1},
		0xf8: {cpu._RM, "RM", 1},
		0xf9: {cpu._SPHL, "SPHL", 1},
		0xfa: {cpu._JM, "JM", 3},
		0xfb: {cpu._EI, "EI", 1},
		0xfc: {cpu._CM, "CM", 3},
		0xfd: {cpu._CALL, "*CALL", 3},
		0xfe: {cpu._CPI, "CPI", 2},
		0xff: {cpu._RST_7, "RST 7", 1},
	}

	return cpu
}

func (cpu *Intel8080) GetMemory() []byte {
	return cpu.memory[:]
}

func (cpu *Intel8080) LoadProgram(program []byte, offset int) {
	for i, b := range program {
		cpu.memory[i+offset] = b
	}
}

func (cpu *Intel8080) SetPC(value uint16) {
	cpu.pc = value
}

func (cpu *Intel8080) WriteIntoMemory(addr uint16, b byte) {
	cpu.memory[addr] = b
}

func (cpu *Intel8080) ReadFromMemory(addr uint16) byte {
	return cpu.memory[addr]
}

func (cpu *Intel8080) SetInputListener(listener func(cpu *Intel8080)) {
	cpu.onInput = listener
}

func (cpu *Intel8080) SetOutputListener(listener func(cpu *Intel8080)) {
	cpu.onOutput = listener
}

func (cpu *Intel8080) GetRegisters() *Intel8080Registers {
	return &Intel8080Registers{
		A: cpu.a,
		B: cpu.b,
		C: cpu.c,
		D: cpu.d,
		E: cpu.e,
		H: cpu.h,
		L: cpu.l,
	}
}

func (cpu *Intel8080) Run() {
	opcode := cpu.memory[cpu.pc]
	cpu.pc++

	if cpu.enableInterruptDeferred {
		cpu.enableInterruptDeferred = false
		cpu.InterruptEnabled = true
	}

	instruction := cpu.instructions[opcode]
	cycles := instruction.operation()

	cpu.cycles += cycles
}

func hasParity(b byte) bool {
	return bits.OnesCount8(b)%2 == 0
}

func (cpu *Intel8080) add(value byte, carry uint16) {
	result := uint16(cpu.a) + uint16(value) + carry

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

	cpu.add(value, carryVal)
}

func (cpu *Intel8080) sub(value byte, carry uint16) {
	result := uint16(cpu.a) - uint16(value) - carry

	cpu.flags.Set(Carry, result>>8 > 0)
	cpu.flags.Set(AuxCarry, ((cpu.a^uint8(result)^value)&0x10) > 0)

	cpu.a = uint8(result & 0xFF)

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
}

func (cpu *Intel8080) sbb(value byte) {
	carryVal := uint16(0)
	if cpu.flags.Get(Carry) {
		carryVal = 1
	}

	cpu.sub(value, carryVal)
}

func (cpu *Intel8080) ana(value byte) {
	result := cpu.a & value

	cpu.flags.Set(AuxCarry, ((cpu.a|value)&0x08) != 0)
	cpu.flags.Set(Carry, false)

	cpu.a = result

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
}

func (cpu *Intel8080) xra(value byte) {
	cpu.a ^= value

	cpu.flags.Set(AuxCarry, false)
	cpu.flags.Set(Carry, false)
	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
}

func (cpu *Intel8080) ora(value byte) {
	cpu.a |= value

	cpu.flags.Set(AuxCarry, false)
	cpu.flags.Set(Carry, false)
	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))
}

func (cpu *Intel8080) cmp(value byte) {
	result := uint16(cpu.a) - uint16(value)

	cpu.flags.Set(Carry, result>>8 > 0)
	cpu.flags.Set(AuxCarry, ((cpu.a^uint8(result)^value)&0x10) > 0)

	cpu.flags.Set(Zero, uint8(result) == 0)
	cpu.flags.Set(Sign, uint8(result)&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(uint8(result)))
}

func (cpu *Intel8080) ret() {
	lb, hb := uint16(cpu.memory[cpu.sp]), uint16(cpu.memory[cpu.sp+1])
	cpu.sp += 2
	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) pop() (hb byte, lb byte) {
	lob, hib := cpu.memory[cpu.sp], cpu.memory[cpu.sp+1]
	cpu.sp += 2

	return hib, lob
}

func (cpu *Intel8080) push(hb byte, lb byte) {
	cpu.memory[cpu.sp-1] = hb
	cpu.memory[cpu.sp-2] = lb
	cpu.sp -= 2
}

func (cpu *Intel8080) call() {
	lb, hb := uint16(cpu.memory[cpu.pc]), uint16(cpu.memory[cpu.pc+1])

	ret := cpu.pc + 2
	cpu.memory[cpu.sp-1] = uint8((ret >> 8) & 0xff)
	cpu.memory[cpu.sp-2] = uint8(ret & 0xff)
	cpu.sp -= 2

	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) jump() {
	lb := uint16(cpu.memory[cpu.pc])
	hb := uint16(cpu.memory[cpu.pc+1])

	cpu.pc = (hb << 8) | lb
}

func (cpu *Intel8080) rst(addr uint16) {
	ret := cpu.pc
	cpu.memory[cpu.sp-1] = uint8((ret >> 8) & 0xFF)
	cpu.memory[cpu.sp-2] = uint8(ret & 0xFF)
	cpu.sp -= 2

	cpu.pc = addr
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

func (cpu *Intel8080) _DAA() uint {
	daa := uint16(cpu.a)

	if (daa&0x0F) > 0x09 || cpu.flags.Get(AuxCarry) {
		cpu.flags.Set(AuxCarry, (((daa&0x0F)+0x06)&0xF0) != 0)
		daa += 0x06
		if (daa & 0xFF00) != 0 {
			cpu.flags.Set(Carry, true)
		}
	}

	if (daa&0xF0) > 0x90 || cpu.flags.Get(Carry) {
		daa += 0x60
		if (daa & 0xFF00) != 0 {
			cpu.flags.Set(Carry, true)
		}
	}

	cpu.a = byte(daa & 0xFF)

	cpu.flags.Set(Zero, cpu.a == 0)
	cpu.flags.Set(Sign, cpu.a&0x80 != 0)
	cpu.flags.Set(Parity, hasParity(cpu.a))

	return 4
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

func (cpu *Intel8080) _HLT() uint {
	// Note: It starts an infinite loop here
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
	cpu.add(cpu.b, 0)
	return 4
}

func (cpu *Intel8080) _ADD_C() uint {
	cpu.add(cpu.c, 0)
	return 4
}

func (cpu *Intel8080) _ADD_D() uint {
	cpu.add(cpu.d, 0)
	return 4
}

func (cpu *Intel8080) _ADD_E() uint {
	cpu.add(cpu.e, 0)
	return 4
}

func (cpu *Intel8080) _ADD_H() uint {
	cpu.add(cpu.h, 0)
	return 4
}

func (cpu *Intel8080) _ADD_L() uint {
	cpu.add(cpu.l, 0)
	return 4
}

func (cpu *Intel8080) _ADD_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.add(cpu.memory[addr], 0)
	return 7
}

func (cpu *Intel8080) _ADD_A() uint {
	cpu.add(cpu.a, 0)
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

func (cpu *Intel8080) _SUB_B() uint {
	cpu.sub(cpu.b, 0)
	return 4
}

func (cpu *Intel8080) _SUB_C() uint {
	cpu.sub(cpu.c, 0)
	return 4
}

func (cpu *Intel8080) _SUB_D() uint {
	cpu.sub(cpu.d, 0)
	return 4
}

func (cpu *Intel8080) _SUB_E() uint {
	cpu.sub(cpu.e, 0)
	return 4
}

func (cpu *Intel8080) _SUB_H() uint {
	cpu.sub(cpu.h, 0)
	return 4
}

func (cpu *Intel8080) _SUB_L() uint {
	cpu.sub(cpu.l, 0)
	return 4
}

func (cpu *Intel8080) _SUB_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.sub(cpu.memory[addr], 0)
	return 7
}

func (cpu *Intel8080) _SUB_A() uint {
	cpu.sub(cpu.a, 0)
	return 4
}

func (cpu *Intel8080) _SBB_B() uint {
	cpu.sbb(cpu.b)
	return 4
}

func (cpu *Intel8080) _SBB_C() uint {
	cpu.sbb(cpu.c)
	return 4
}

func (cpu *Intel8080) _SBB_D() uint {
	cpu.sbb(cpu.d)
	return 4
}

func (cpu *Intel8080) _SBB_E() uint {
	cpu.sbb(cpu.e)
	return 4
}

func (cpu *Intel8080) _SBB_H() uint {
	cpu.sbb(cpu.h)
	return 4
}

func (cpu *Intel8080) _SBB_L() uint {
	cpu.sbb(cpu.l)
	return 4
}

func (cpu *Intel8080) _SBB_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.sbb(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _SBB_A() uint {
	cpu.sbb(cpu.a)
	return 4
}

func (cpu *Intel8080) _ANA_B() uint {
	cpu.ana(cpu.b)
	return 4
}

func (cpu *Intel8080) _ANA_C() uint {
	cpu.ana(cpu.c)
	return 4
}

func (cpu *Intel8080) _ANA_D() uint {
	cpu.ana(cpu.d)
	return 4
}

func (cpu *Intel8080) _ANA_E() uint {
	cpu.ana(cpu.e)
	return 4
}

func (cpu *Intel8080) _ANA_H() uint {
	cpu.ana(cpu.h)
	return 4
}

func (cpu *Intel8080) _ANA_L() uint {
	cpu.ana(cpu.l)
	return 4
}

func (cpu *Intel8080) _ANA_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.ana(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _ANA_A() uint {
	cpu.ana(cpu.a)
	return 4
}

func (cpu *Intel8080) _XRA_B() uint {
	cpu.xra(cpu.b)
	return 4
}

func (cpu *Intel8080) _XRA_C() uint {
	cpu.xra(cpu.c)
	return 4
}

func (cpu *Intel8080) _XRA_D() uint {
	cpu.xra(cpu.d)
	return 4
}

func (cpu *Intel8080) _XRA_E() uint {
	cpu.xra(cpu.e)
	return 4
}

func (cpu *Intel8080) _XRA_H() uint {
	cpu.xra(cpu.h)
	return 4
}

func (cpu *Intel8080) _XRA_L() uint {
	cpu.xra(cpu.l)
	return 4
}

func (cpu *Intel8080) _XRA_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.xra(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _XRA_A() uint {
	cpu.xra(cpu.a)
	return 4
}

func (cpu *Intel8080) _ORA_B() uint {
	cpu.ora(cpu.b)
	return 4
}

func (cpu *Intel8080) _ORA_C() uint {
	cpu.ora(cpu.c)
	return 4
}

func (cpu *Intel8080) _ORA_D() uint {
	cpu.ora(cpu.d)
	return 4
}

func (cpu *Intel8080) _ORA_E() uint {
	cpu.ora(cpu.e)
	return 4
}

func (cpu *Intel8080) _ORA_H() uint {
	cpu.ora(cpu.h)
	return 4
}

func (cpu *Intel8080) _ORA_L() uint {
	cpu.ora(cpu.l)
	return 4
}

func (cpu *Intel8080) _ORA_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.ora(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _ORA_A() uint {
	cpu.ora(cpu.a)
	return 4
}

func (cpu *Intel8080) _CMP_B() uint {
	cpu.cmp(cpu.b)
	return 4
}

func (cpu *Intel8080) _CMP_C() uint {
	cpu.cmp(cpu.c)
	return 4
}

func (cpu *Intel8080) _CMP_D() uint {
	cpu.cmp(cpu.d)
	return 4
}

func (cpu *Intel8080) _CMP_E() uint {
	cpu.cmp(cpu.e)
	return 4
}

func (cpu *Intel8080) _CMP_H() uint {
	cpu.cmp(cpu.h)
	return 4
}

func (cpu *Intel8080) _CMP_L() uint {
	cpu.cmp(cpu.l)
	return 4
}

func (cpu *Intel8080) _CMP_M() uint {
	addr := uint16(cpu.h)<<8 | uint16(cpu.l)
	cpu.cmp(cpu.memory[addr])
	return 7
}

func (cpu *Intel8080) _CMP_A() uint {
	cpu.cmp(cpu.a)
	return 4
}

func (cpu *Intel8080) _RNZ() uint {
	if !cpu.flags.Get(Zero) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _POP_B() uint {
	cpu.b, cpu.c = cpu.pop()
	return 10
}

func (cpu *Intel8080) _JNZ() uint {
	if !cpu.flags.Get(Zero) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _JMP() uint {
	cpu.jump()
	return 10
}

func (cpu *Intel8080) _CNZ() uint {
	if !cpu.flags.Get(Zero) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _PUSH_B() uint {
	cpu.push(cpu.b, cpu.c)
	return 11
}

func (cpu *Intel8080) _ADI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.add(value, 0)
	return 7
}

func (cpu *Intel8080) _RST_0() uint {
	cpu.rst(0x0000)
	return 11
}

func (cpu *Intel8080) _RZ() uint {
	if cpu.flags.Get(Zero) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _RET() uint {
	cpu.ret()
	return 10
}

func (cpu *Intel8080) _JZ() uint {
	if cpu.flags.Get(Zero) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _CZ() uint {
	if cpu.flags.Get(Zero) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _ACI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.adc(value)
	return 7
}

func (cpu *Intel8080) _CALL() uint {
	cpu.call()
	return 17
}

func (cpu *Intel8080) _RST_1() uint {
	cpu.rst(0x0008)
	return 11
}

func (cpu *Intel8080) _RNC() uint {
	if !cpu.flags.Get(Carry) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _POP_D() uint {
	cpu.d, cpu.e = cpu.pop()
	return 10
}

func (cpu *Intel8080) _JNC() uint {
	if !cpu.flags.Get(Carry) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _OUT() uint {
	//TODO: Revisiting this later to implement it correctly
	if cpu.onOutput != nil {
		cpu.onOutput(cpu)
	}

	cpu.pc++
	return 10
}

func (cpu *Intel8080) _CNC() uint {
	if !cpu.flags.Get(Carry) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _PUSH_D() uint {
	cpu.push(cpu.d, cpu.e)
	return 11
}

func (cpu *Intel8080) _SUI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.sub(value, 0)
	return 7
}

func (cpu *Intel8080) _RST_2() uint {
	cpu.rst(0x0010)
	return 11
}

func (cpu *Intel8080) _RC() uint {
	if cpu.flags.Get(Carry) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _JC() uint {
	if cpu.flags.Get(Carry) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _IN() uint {
	//TODO: Revisiting this later to implement it correctly
	if cpu.onInput != nil {
		cpu.onInput(cpu)
	}

	cpu.pc++
	return 10
}

func (cpu *Intel8080) _CC() uint {
	if cpu.flags.Get(Carry) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _SBI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.sbb(value)
	return 7
}

func (cpu *Intel8080) _RST_3() uint {
	cpu.rst(0x0018)
	return 11
}

func (cpu *Intel8080) _RPO() uint {
	if !cpu.flags.Get(Parity) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _POP_H() uint {
	cpu.h, cpu.l = cpu.pop()
	return 10
}

func (cpu *Intel8080) _JPO() uint {
	if !cpu.flags.Get(Parity) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _XTHL() uint {
	stackLb := cpu.memory[cpu.sp]
	stackHb := cpu.memory[cpu.sp+1]
	cpu.memory[cpu.sp] = cpu.l
	cpu.memory[cpu.sp+1] = cpu.h
	cpu.l = stackLb
	cpu.h = stackHb
	return 18
}

func (cpu *Intel8080) _CPO() uint {
	if !cpu.flags.Get(Parity) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _PUSH_H() uint {
	cpu.push(cpu.h, cpu.l)
	return 11
}

func (cpu *Intel8080) _ANI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.ana(value)
	return 7
}

func (cpu *Intel8080) _RST_4() uint {
	cpu.rst(0x0020)
	return 11
}

func (cpu *Intel8080) _RPE() uint {
	if cpu.flags.Get(Parity) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _PCHL() uint {
	cpu.pc = uint16(cpu.h)<<8 | uint16(cpu.l)
	return 5
}

func (cpu *Intel8080) _JPE() uint {
	if cpu.flags.Get(Parity) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _XCHG() uint {
	tmpD := cpu.d
	tmpE := cpu.e
	cpu.d = cpu.h
	cpu.e = cpu.l
	cpu.h = tmpD
	cpu.l = tmpE
	return 5
}

func (cpu *Intel8080) _CPE() uint {
	if cpu.flags.Get(Parity) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _XRI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.xra(value)
	return 7
}

func (cpu *Intel8080) _RST_5() uint {
	cpu.rst(0x0028)
	return 11
}

func (cpu *Intel8080) _RP() uint {
	if !cpu.flags.Get(Sign) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _POP_PSW() uint {
	hb, psw := cpu.pop()
	cpu.a = hb
	cpu.flags.Set(Sign, (psw>>7&0b1) > 0)
	cpu.flags.Set(Zero, (psw>>6&0b1) > 0)
	cpu.flags.Set(AuxCarry, (psw>>4&0b1) > 0)
	cpu.flags.Set(Parity, (psw>>2&0b1) > 0)
	cpu.flags.Set(Carry, (psw>>0&0b1) > 0)
	return 10
}

func (cpu *Intel8080) _JP() uint {
	if !cpu.flags.Get(Sign) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _DI() uint {
	cpu.InterruptEnabled = false
	return 4
}

func (cpu *Intel8080) _CP() uint {
	if !cpu.flags.Get(Sign) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _PUSH_PSW() uint {
	// S, Z, 0, AC, 0, P, 1, CY
	status := uint8(0b00000010)
	if cpu.flags.Get(Sign) {
		status |= 1 << 7
	}
	if cpu.flags.Get(Zero) {
		status |= 1 << 6
	}
	if cpu.flags.Get(AuxCarry) {
		status |= 1 << 4
	}
	if cpu.flags.Get(Parity) {
		status |= 1 << 2
	}
	if cpu.flags.Get(Carry) {
		status |= 1 << 0
	}

	cpu.push(cpu.a, status)

	return 11
}

func (cpu *Intel8080) _ORI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.ora(value)
	return 7
}

func (cpu *Intel8080) _RST_6() uint {
	cpu.rst(0x0030)
	return 11
}

func (cpu *Intel8080) _RM() uint {
	if cpu.flags.Get(Sign) {
		cpu.ret()
		return 11
	}
	return 5
}

func (cpu *Intel8080) _SPHL() uint {
	cpu.sp = uint16(cpu.h)<<8 | uint16(cpu.l)
	return 5
}

func (cpu *Intel8080) _JM() uint {
	if cpu.flags.Get(Sign) {
		cpu.jump()
	} else {
		cpu.pc += 2
	}

	return 10
}

func (cpu *Intel8080) _EI() uint {
	// Note: interrupt enable is deferred until the next instruction has completed
	cpu.enableInterruptDeferred = true
	return 4
}

func (cpu *Intel8080) _CM() uint {
	if cpu.flags.Get(Sign) {
		cpu.call()
		return 17
	}
	cpu.pc += 2
	return 11
}

func (cpu *Intel8080) _CPI() uint {
	value := cpu.memory[cpu.pc]
	cpu.pc++
	cpu.cmp(value)
	return 7
}

func (cpu *Intel8080) _RST_7() uint {
	cpu.rst(0x0038)
	return 11
}
