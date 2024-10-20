package core

import (
	"testing"
)

type flagDataTest struct {
	value    byte
	flagName string
	flagMask byte
}

func createCPUWithProgramLoaded(p []byte) *Intel8080 {
	cpu := NewIntel8080()
	cpu.LoadProgram(p)

	return cpu
}

func TestLoadProgram(t *testing.T) {
	program := []byte{0x00, 0x01, 0x02, 0x03}
	cpu := createCPUWithProgramLoaded(program)

	for i, v := range program {
		if cpu.memory[i] != v {
			t.Errorf("LoadProgram did not load the program correctly")
		}
	}
}

func Test_LXI_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x01, 0x02, 0x03})

	cpu.Run()

	if cpu.c != 0x02 || cpu.b != 0x03 {
		t.Errorf("LXI B did not load the program correctly")
	}
}

func Test_STAX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x02, 0x01})

	cpu.b = 0x03
	cpu.c = 0x01
	cpu.a = 0x08

	cpu.Run()

	if cpu.memory[0x0301] != 0x08 {
		t.Errorf("STAX B did not store the program correctly")
	}
}

func Test_INX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x03, 0x01})

	cpu.b = 0x03
	cpu.c = 0x01

	cpu.Run()

	if cpu.c != 0x02 || cpu.b != 0x03 {
		t.Errorf("INX B did not increment the program correctly")
	}
}

func Test_INR_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x04, 0x01})
	cpu.b = 0x03

	cpu.Run()

	if cpu.b != 0x04 {
		t.Errorf("INR B did not increment the program correctly")
	}
}

func Fuzz_INR_B_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA9, flagName: "Parity", flagMask: Parity},
		{value: 0xFF, flagName: "Zero", flagMask: Zero},
		{value: 0x2F, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7F, flagName: "Sign", flagMask: Sign},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x04, 0x01})

		cpu.b = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR B did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_DCR_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x05, 0x01})
	cpu.b = 0x05

	cpu.Run()

	if cpu.b != 0x04 {
		t.Errorf("DCR B did not increment the program correctly")
	}
}

func Fuzz_DCR_B_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xAB, flagName: "Parity", flagMask: Parity},
		{value: 0x01, flagName: "Zero", flagMask: Zero},
		{value: 0x31, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x81, flagName: "Sign", flagMask: Sign},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x05, 0x01})

		cpu.b = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR B did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_MVI_B_Flags(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x06, 0x42})

	cpu.Run()

	if cpu.b != 0x42 {
		t.Errorf("MVI B did not load the correct value to register")
	}
}

func Test_RLC_Flags(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x07, 0x01})
	cpu.a = 0x80

	cpu.Run()

	if cpu.a != 0x01 {
		t.Errorf("RLC did not set the correct value to register A")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RLC did not set the carry flag correctly")
	}
}

func Test_DAD_B_Flags(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x09, 0x01})
	cpu.h = 0xF0
	cpu.l = 0x12
	cpu.b = 0x44
	cpu.c = 0x55

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("DAD B did not set the carry flag correctly")
	}

	if cpu.l != 0x67 {
		t.Errorf("DAD B did not set the L register correctly")
	}

	if cpu.h != 0x34 {
		t.Errorf("DAD B did not set the H register correctly")
	}
}

func Test_LDAX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0a, 0x01, 0x01, 0x01, 0x01, 0x99})
	cpu.b = 0x00
	cpu.c = 0x05

	cpu.Run()

	if cpu.a != 0x99 {
		t.Errorf("LDAX B did not set the A register correctly")
	}
}

func Test_DCX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0b, 0x01})
	cpu.b = 0x55
	cpu.c = 0x00

	cpu.Run()

	if cpu.b != 0x54 || cpu.c != 0xFF {
		t.Errorf("DCX B did not set the BC register pair correctly")
	}
}

func Test_INR_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0c, 0x01})
	cpu.c = 0x03

	cpu.Run()

	if cpu.c != 0x04 {
		t.Errorf("INR C did not increment the program correctly")
	}
}

func Fuzz_INR_C_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA9, flagName: "Parity", flagMask: Parity},
		{value: 0xFF, flagName: "Zero", flagMask: Zero},
		{value: 0x2F, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7F, flagName: "Sign", flagMask: Sign},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x0c, 0x01})

		cpu.c = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR C did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_DCR_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0d, 0x01})
	cpu.c = 0x05

	cpu.Run()

	if cpu.c != 0x04 {
		t.Errorf("DCR C did not increment the program correctly")
	}
}

func Fuzz_DCR_C_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xAB, flagName: "Parity", flagMask: Parity},
		{value: 0x01, flagName: "Zero", flagMask: Zero},
		{value: 0x31, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x81, flagName: "Sign", flagMask: Sign},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x0d, 0x01})

		cpu.c = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR C did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_RRC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0f, 0x01})
	cpu.a = 0xf

	cpu.Run()

	if cpu.a != 0x87 {
		t.Errorf("RRC did not rotate the A register correctly")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RRC did not set the Carry flag correctly")
	}
}

func Test_RARWithCarryFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x1f, 0x01})
	cpu.flags.Set(Carry, true)
	cpu.a = 0xf

	cpu.Run()

	if cpu.a != 0x87 {
		t.Errorf("RRC did not rotate the A register correctly")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RRC did not set the Carry flag correctly")
	}
}

func Test_RARWithCarryFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x1f, 0x01})
	cpu.flags.Set(Carry, false)
	cpu.a = 0xf

	cpu.Run()

	if cpu.a != 0x7 {
		t.Errorf("RRC did not rotate the A register correctly")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RRC did not set the Carry flag correctly")
	}
}

func Test_CMA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x2f, 0x01})
	cpu.a = 0xdd

	cpu.Run()

	if cpu.a != 0x22 {
		t.Errorf("CMA did not set the A register correctly")
	}
}

func Test_JNZ_ZeroFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc2, 0x88, 0xff})
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JNZ dit not set pc correctly")
	}
}

func Test_JNZ_ZeroFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc2, 0x88, 0xff, 0x01})
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JNZ dit not set pc correctly")
	}
}

func Test_JMP(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc3, 0x96, 0xed})

	cpu.Run()

	if cpu.pc != 0xed96 {
		t.Errorf("JMP dit not set pc correctly")
	}
}

func TestCALL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xcd, 0x34, 0x12, 0x09, 0x00, 0x00, 0x00})
	cpu.sp = 6

	cpu.Run()

	if cpu.sp != 4 {
		t.Errorf("CALL dit not set SP correctly")
	}

	if cpu.pc != 0x1234 {
		t.Errorf("CALL dit not set PC correctly")
	}

	if cpu.memory[cpu.sp+1] != 0x00 || cpu.memory[cpu.sp] != 0x03 {
		t.Errorf("CALL dit not write correctly to memory")
	}
}

func TestRET(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc9, 0x34, 0x12, 0x00})
	cpu.sp = 1

	cpu.Run()

	if cpu.sp != 3 {
		t.Errorf("RET dit not set SP correctly")
	}

	if cpu.pc != 0x1234 {
		t.Errorf("CALL dit not set PC correctly")
	}
}

func Test_ANI(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe6, 0x05, 0x01})
	cpu.flags.Set(Carry, true)
	cpu.flags.Set(AuxCarry, true)
	cpu.a = 0x10

	cpu.Run()

	if cpu.a != 0x00 {
		t.Errorf("ANI did not set the A register correctly")
	}

	if cpu.flags.Get(Carry) || cpu.flags.Get(AuxCarry) {
		t.Errorf("ANI did not clear Carry and AuxCarry flags")
	}

	if !cpu.flags.Get(Zero) {
		t.Errorf("ANI did not set Zero flag correctly")
	}

	if cpu.flags.Get(Sign) {
		t.Errorf("ANI did not set Sign flag correctly")
	}

	if !cpu.flags.Get(Parity) {
		t.Errorf("ANI did not set Parity flag correctly")
	}
}

func Fuzz_CPI_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xAB, flagName: "Parity", flagMask: Parity},
		{value: 0x01, flagName: "Zero", flagMask: Zero},
		{value: 0x81, flagName: "Sign", flagMask: Sign},
		{value: 0x00, flagName: "Carry", flagMask: Carry},
		{value: 0x00, flagName: "AuxCarry", flagMask: AuxCarry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0xfe, 0x01, 0x01})
		cpu.a = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CPI did not set the %s flag correctly", d.flagName)
		}

		if cpu.pc != 2 {
			t.Errorf("CPI did not increment PC")
		}
	})
}
