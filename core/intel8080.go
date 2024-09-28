package core

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
	cpu := &Intel8080{}

	cpu.instructions = [256]*Intel8080Instruction{
		{cpu._NOP, "NOP", 1}, {cpu._LXI_B, "LXI B", 3}, {cpu._STAX_B, "STAX B", 1}, {cpu._INX_B, "INX B", 1},
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
