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

func assertCycles(t *testing.T, cpu *Intel8080, expected uint) {
	if cpu.cycles != expected {
		t.Errorf("cpu cycles have not been set correctly")
	}
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
		t.Errorf("LXI B did not set registers correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LXI B did not increment PCs correctly")
	}

	assertCycles(t, cpu, 10)
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

	assertCycles(t, cpu, 7)
}

func Test_INX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x03, 0x01})

	cpu.b = 0x03
	cpu.c = 0x01

	cpu.Run()

	if cpu.c != 0x02 || cpu.b != 0x03 {
		t.Errorf("INX B did not increment the program correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_B(f *testing.F) {
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

		if cpu.b != d.value+1 {
			t.Errorf("INR B did not increment the program correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_B(f *testing.F) {
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

		if cpu.b != d.value-1 {
			t.Errorf("DCR B did not decrement B register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x06, 0x42})

	cpu.Run()

	if cpu.b != 0x42 {
		t.Errorf("MVI B did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI B did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_RLC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x07, 0x01})
	cpu.a = 0x80

	cpu.Run()

	if cpu.a != 0x01 {
		t.Errorf("RLC did not set the correct value to register A")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RLC did not set the carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_DAD_B(t *testing.T) {
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

	assertCycles(t, cpu, 10)
}

func Test_LDAX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0a, 0x01, 0x01, 0x01, 0x01, 0x99})
	cpu.b = 0x00
	cpu.c = 0x05

	cpu.Run()

	if cpu.a != 0x99 {
		t.Errorf("LDAX B did not set the A register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_DCX_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0b, 0x01})
	cpu.b = 0x55
	cpu.c = 0x00

	cpu.Run()

	if cpu.b != 0x54 || cpu.c != 0xFF {
		t.Errorf("DCX B did not set the BC register pair correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_C(f *testing.F) {
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

		if cpu.c != d.value+1 {
			t.Errorf("INR C did not increment the program correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_C(f *testing.F) {
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

		if cpu.c != d.value-1 {
			t.Errorf("DCR C did not decrement C register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x0e, 0x42})

	cpu.Run()

	if cpu.c != 0x42 {
		t.Errorf("MVI C did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI C did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
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

	assertCycles(t, cpu, 4)
}

func Test_LXI_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x11, 0x02, 0x03})

	cpu.Run()

	if cpu.e != 0x02 || cpu.d != 0x03 {
		t.Errorf("LXI D did not set registers correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LXI B did not increment PC correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_STAX_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x12, 0x01})

	cpu.d = 0x03
	cpu.e = 0x01
	cpu.a = 0x08

	cpu.Run()

	if cpu.memory[0x0301] != 0x08 {
		t.Errorf("STAX D did not store the program correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_INX_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x13, 0x01})

	cpu.d = 0x03
	cpu.e = 0x01

	cpu.Run()

	if cpu.e != 0x02 || cpu.d != 0x03 {
		t.Errorf("INX D did not increment the program correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_D(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x14, 0x01})

		cpu.d = d.value

		cpu.Run()

		if cpu.d != d.value+1 {
			t.Errorf("INR D did not increment D register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_D(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x15, 0x01})

		cpu.d = d.value

		cpu.Run()

		if cpu.d != d.value-1 {
			t.Errorf("DCR D did not decrement D register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x16, 0x42})

	cpu.Run()

	if cpu.d != 0x42 {
		t.Errorf("MVI D did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI D did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
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

	assertCycles(t, cpu, 4)
}

func Test_RALWithCarryFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x17, 0x01})
	cpu.flags.Set(Carry, false)
	cpu.a = 0x8F

	cpu.Run()

	if cpu.a != 0x1E {
		t.Errorf("RAL did not rotate the A register correctly")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RAL did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_RALWithCarryFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x17, 0x01})
	cpu.flags.Set(Carry, true)
	cpu.a = 0xf

	cpu.Run()

	if cpu.a != 0x1F {
		t.Errorf("RAL did not rotate the A register correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("RAL did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_DAD_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x19, 0x01})
	cpu.h = 0xF0
	cpu.l = 0x12
	cpu.d = 0x44
	cpu.e = 0x55

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("DAD D did not set the carry flag correctly")
	}

	if cpu.l != 0x67 {
		t.Errorf("DAD D did not set the L register correctly")
	}

	if cpu.h != 0x34 {
		t.Errorf("DAD D did not set the H register correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_LDAX_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x1a, 0x01, 0x01, 0x01, 0x01, 0x99})
	cpu.d = 0x00
	cpu.e = 0x05

	cpu.Run()

	if cpu.a != 0x99 {
		t.Errorf("LDAX D did not set the A register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_DCX_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x1b, 0x01})
	cpu.d = 0x55
	cpu.e = 0x00

	cpu.Run()

	if cpu.d != 0x54 || cpu.e != 0xFF {
		t.Errorf("DCX D did not set the DE register pair correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_E(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x1c, 0x01})

		cpu.e = d.value

		cpu.Run()

		if cpu.e != d.value+1 {
			t.Errorf("INR E did not increment E register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_E(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x1d, 0x01})

		cpu.e = d.value

		cpu.Run()

		if cpu.e != d.value-1 {
			t.Errorf("DCR E did not decrement E register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_E(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x1e, 0x42})

	cpu.Run()

	if cpu.e != 0x42 {
		t.Errorf("MVI E did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI E did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
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

func Test_LXI_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x21, 0x02, 0x03})

	cpu.Run()

	if cpu.l != 0x02 || cpu.h != 0x03 {
		t.Errorf("LXI H did not set registers correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LXI H did not increment PC correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_SHLD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x22, 0x03, 0x00, 0x00, 0x00})
	cpu.l = 0x55
	cpu.h = 0x66

	cpu.Run()

	if cpu.memory[0x0003] != 0x55 || cpu.memory[0x0004] != 0x66 {
		t.Errorf("SHLD did not write into memory correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("SHLD did not increment PC correctly")
	}

	assertCycles(t, cpu, 16)
}

func Test_INX_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x23, 0x01})

	cpu.h = 0x03
	cpu.l = 0x01

	cpu.Run()

	if cpu.l != 0x02 || cpu.h != 0x03 {
		t.Errorf("INX H did not increment the program correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_H(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x24, 0x01})

		cpu.h = d.value

		cpu.Run()

		if cpu.h != d.value+1 {
			t.Errorf("INR H did not increment H register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_H(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x25, 0x01})

		cpu.h = d.value

		cpu.Run()

		if cpu.h != d.value-1 {
			t.Errorf("DCR H did not decrement H register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x26, 0x42})

	cpu.Run()

	if cpu.h != 0x42 {
		t.Errorf("MVI H did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI H did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_DAD_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x29, 0x01})
	cpu.h = 0xF0
	cpu.l = 0x12

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("DAD H did not set the carry flag correctly")
	}

	if cpu.l != 0x24 {
		t.Errorf("DAD H did not set the L register correctly")
	}

	if cpu.h != 0xe0 {
		t.Errorf("DAD H did not set the H register correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_LHLD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x2a, 0x03, 0x00, 0x55, 0x66})

	cpu.Run()

	if cpu.l != 0x55 || cpu.h != 0x66 {
		t.Errorf("LHLD did not set HL registers correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LHLD did not increment PC correctly")
	}

	assertCycles(t, cpu, 16)
}

func Test_DCX_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x2b, 0x01})
	cpu.h = 0x55
	cpu.l = 0x00

	cpu.Run()

	if cpu.h != 0x54 || cpu.l != 0xFF {
		t.Errorf("DCX H did not set the HL register pair correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_L(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x2c, 0x01})

		cpu.l = d.value

		cpu.Run()

		if cpu.l != d.value+1 {
			t.Errorf("INR L did not increment L register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_L(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x2d, 0x01})

		cpu.l = d.value

		cpu.Run()

		if cpu.l != d.value-1 {
			t.Errorf("DCR L did not decrement L register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_L(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x2e, 0x42})

	cpu.Run()

	if cpu.l != 0x42 {
		t.Errorf("MVI L did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI L did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_CMA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x2f, 0x01})
	cpu.a = 0xdd

	cpu.Run()

	if cpu.a != 0x22 {
		t.Errorf("CMA did not set the A register correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_LXI_SP(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x31, 0x02, 0x03})

	cpu.Run()

	if cpu.sp != 0x0302 {
		t.Errorf("LXI SP did not set Stack Pointer correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LXI SP did not increment PC correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_STA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x32, 0x03, 0x00, 0x00})
	cpu.a = 0x99

	cpu.Run()

	if cpu.memory[0x0003] != 0x99 {
		t.Errorf("STA did not write A into memory correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("STA did not increment PC correctly")
	}

	assertCycles(t, cpu, 13)
}

func Test_INX_SP(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x33, 0x01})
	cpu.sp = 0x5

	cpu.Run()

	if cpu.sp != 0x6 {
		t.Errorf("INX SP did not increment Stack Pointer correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_M(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x34, 0x00, 0x00, d.value})
		cpu.h = 0x00
		cpu.l = 0x03

		cpu.Run()

		if cpu.memory[0x0003] != d.value+1 {
			t.Errorf("INR M did not increment value in memory correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 10)
	})
}

func Fuzz_DCR_M(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x35, 0x00, 0x00, 0x00, d.value})
		cpu.h = 0x00
		cpu.l = 0x04

		cpu.Run()

		if cpu.memory[0x0004] != d.value-1 {
			t.Errorf("INR M did not decrement value in memory correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 10)
	})
}

func Test_MVI_M(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x36, 0xED, 0x00, 0x00, 0x00, 0x00, 0x01FF: 0x00})
	cpu.h = 0x01
	cpu.l = 0xFF

	cpu.Run()

	if cpu.memory[0x01FF] != 0xED {
		t.Errorf("MVI M did not store the correct value to memory")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI M did not increment PC correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_STC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x37, 0x00})
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("STC did not set Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_DAD_SP(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x39, 0x01})
	cpu.h = 0xF0
	cpu.l = 0x12
	cpu.sp = 0x4455

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("DAD SP did not set the carry flag correctly")
	}

	if cpu.l != 0x67 {
		t.Errorf("DAD SP did not set the L register correctly")
	}

	if cpu.h != 0x34 {
		t.Errorf("DAD SP did not set the H register correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_LDA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x3A, 0x33, 0x20, 0x2033: 0x76})

	cpu.Run()

	if cpu.a != 0x76 {
		t.Errorf("LDA did not set A register correctly")
	}

	if cpu.pc != 3 {
		t.Errorf("LDA did not increment PC correctly")
	}

	assertCycles(t, cpu, 13)
}

func Test_DCX_SP(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x3b, 0x01})
	cpu.sp = 0x5

	cpu.Run()

	if cpu.sp != 0x4 {
		t.Errorf("dcx SP did not decrement Stack Pointer correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_INR_A(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x3c, 0x01})
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value+1 {
			t.Errorf("INR A did not increment A register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Fuzz_DCR_A(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0x3d, 0x01})

		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-1 {
			t.Errorf("DCR A did not decrement A register correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("DCR A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 5)
	})
}

func Test_MVI_A(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x3e, 0x68})

	cpu.Run()

	if cpu.a != 0x68 {
		t.Errorf("MVI A did not load the correct value to register")
	}

	if cpu.pc != 2 {
		t.Errorf("MVI A did not increment PC correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_CMCWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x3f, 0x01})
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if !cpu.flags.Get(Carry) {
		t.Errorf("CMC did not set Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_CMCWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x3f, 0x01})
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.flags.Get(Carry) {
		t.Errorf("CMC did not set Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_MOV_BB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x40, 0x01})
	cpu.b = 0x5

	cpu.Run()

	if cpu.b != 0x5 {
		t.Errorf("MOV B,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x41, 0x01})
	cpu.b = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x42, 0x01})
	cpu.b = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x43, 0x01})
	cpu.b = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x44, 0x01})
	cpu.b = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x45, 0x01})
	cpu.b = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,L did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_BM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x46, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.b = 0x1

	cpu.Run()

	if cpu.b != 0x89 {
		t.Errorf("MOV B,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_BA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x47, 0x01})
	cpu.b = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.b != 0x8 {
		t.Errorf("MOV B,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x48, 0x01})
	cpu.c = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x49, 0x01})
	cpu.c = 0x5

	cpu.Run()

	if cpu.c != 0x5 {
		t.Errorf("MOV C,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4a, 0x01})
	cpu.c = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4b, 0x01})
	cpu.c = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4c, 0x01})
	cpu.c = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4d, 0x01})
	cpu.c = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_CM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4e, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.c = 0x1

	cpu.Run()

	if cpu.c != 0x89 {
		t.Errorf("MOV C,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_CA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x4f, 0x01})
	cpu.c = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.c != 0x8 {
		t.Errorf("MOV C,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x50, 0x01})
	cpu.d = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x51, 0x01})
	cpu.d = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x52, 0x01})
	cpu.d = 0x5

	cpu.Run()

	if cpu.d != 0x5 {
		t.Errorf("MOV D,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x53, 0x01})
	cpu.d = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x54, 0x01})
	cpu.d = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x55, 0x01})
	cpu.d = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,L did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_DM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x56, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.d = 0x1

	cpu.Run()

	if cpu.d != 0x89 {
		t.Errorf("MOV D,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_DA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x57, 0x01})
	cpu.d = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.d != 0x8 {
		t.Errorf("MOV D,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x58, 0x01})
	cpu.e = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x59, 0x01})
	cpu.e = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_ED(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5a, 0x01})
	cpu.e = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5b, 0x01})
	cpu.e = 0x5

	cpu.Run()

	if cpu.e != 0x5 {
		t.Errorf("MOV E,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5c, 0x01})
	cpu.e = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5d, 0x01})
	cpu.e = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,L did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_EM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5e, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.e = 0x1

	cpu.Run()

	if cpu.e != 0x89 {
		t.Errorf("MOV E,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_EA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x5f, 0x01})
	cpu.e = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.e != 0x8 {
		t.Errorf("MOV E,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x60, 0x01})
	cpu.h = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x61, 0x01})
	cpu.h = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x62, 0x01})
	cpu.h = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x63, 0x01})
	cpu.h = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x64, 0x01})
	cpu.h = 0x5

	cpu.Run()

	if cpu.h != 0x5 {
		t.Errorf("MOV H,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x65, 0x01})
	cpu.h = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,L did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_HM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x66, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33

	cpu.Run()

	if cpu.h != 0x89 {
		t.Errorf("MOV H,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_HA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x67, 0x01})
	cpu.h = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.h != 0x8 {
		t.Errorf("MOV H,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x68, 0x01})
	cpu.l = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x69, 0x01})
	cpu.l = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6a, 0x01})
	cpu.l = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6b, 0x01})
	cpu.l = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6c, 0x01})
	cpu.l = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6d, 0x01})
	cpu.l = 0x5

	cpu.Run()

	if cpu.l != 0x5 {
		t.Errorf("MOV L,L did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_LM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6e, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33

	cpu.Run()

	if cpu.l != 0x89 {
		t.Errorf("MOV L,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_LA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x6f, 0x01})
	cpu.l = 0x1
	cpu.a = 0x8

	cpu.Run()

	if cpu.l != 0x8 {
		t.Errorf("MOV L,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_MB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x70, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.b = 0x90

	cpu.Run()

	if cpu.memory[0x2233] != cpu.b {
		t.Errorf("MOV M,B did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_MC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x71, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.c = 0x90

	cpu.Run()

	if cpu.memory[0x2233] != cpu.c {
		t.Errorf("MOV M,C did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_MD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x72, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.d = 0x90

	cpu.Run()

	if cpu.memory[0x2233] != cpu.d {
		t.Errorf("MOV M,D did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_ME(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x73, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.e = 0x90

	cpu.Run()

	if cpu.memory[0x2233] != cpu.e {
		t.Errorf("MOV M,E did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_MH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x74, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33

	cpu.Run()

	if cpu.memory[0x2233] != cpu.h {
		t.Errorf("MOV M,H did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_ML(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x75, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33

	cpu.Run()

	if cpu.memory[0x2233] != cpu.l {
		t.Errorf("MOV M,L did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_MA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x77, 0x01, 0x2233: 0x00})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.a = 0x90

	cpu.Run()

	if cpu.memory[0x2233] != cpu.a {
		t.Errorf("MOV M,A did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_AB(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x78, 0x01})
	cpu.a = 0x1
	cpu.b = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,B did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AC(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x79, 0x01})
	cpu.a = 0x1
	cpu.c = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,C did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AD(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7a, 0x01})
	cpu.a = 0x1
	cpu.d = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,D did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AE(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7b, 0x01})
	cpu.a = 0x1
	cpu.e = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,E did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AH(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7c, 0x01})
	cpu.a = 0x1
	cpu.h = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,H did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7d, 0x01})
	cpu.a = 0x1
	cpu.l = 0x8

	cpu.Run()

	if cpu.a != 0x8 {
		t.Errorf("MOV A,LL did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_MOV_AM(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7e, 0x01, 0x2233: 0x89})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.a = 0x1

	cpu.Run()

	if cpu.a != 0x89 {
		t.Errorf("MOV A,M did not move register correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_MOV_AA(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x7f, 0x01})
	cpu.a = 0x5

	cpu.Run()

	if cpu.a != 0x5 {
		t.Errorf("MOV A,A did not move register correctly")
	}

	assertCycles(t, cpu, 5)
}

func Fuzz_ADD_B(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x80, 0x00, 0x00, 0x00})
		cpu.b = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD B did not add A + B correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_C(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x81, 0x00, 0x00, 0x00})
		cpu.c = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD C did not add A + C correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_D(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x82, 0x00, 0x00, 0x00})
		cpu.d = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD D did not add A + D correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_E(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x83, 0x00, 0x00, 0x00})
		cpu.e = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD E did not add A + E correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_H(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x84, 0x00, 0x00, 0x00})
		cpu.h = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD H did not add A + H correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_L(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x85, 0x00, 0x00, 0x00})
		cpu.l = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD L did not add A + L correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADD_M(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0xFB, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x86, 0x00, 0x00, 0x00, 0x2233: d.value})
		cpu.h = 0x22
		cpu.l = 0x33
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADD M did not add A + (HL) correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Fuzz_ADD_A(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0x80, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x87, 0x00, 0x00, 0x00})
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value+d.value {
			t.Errorf("ADD A did not add A + A correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADD A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_B(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x88, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.b = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC B did not add A + B + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_C(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x89, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.c = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC C did not add A + C + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_D(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8a, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.d = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC D did not add A + D + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_E(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8b, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.e = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC E did not add A + E + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_H(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8c, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.h = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC H did not add A + H + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_L(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8d, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.l = d.value
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC L did not add A + L + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_ADC_M(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA3, flagName: "Parity", flagMask: Parity},
		{value: 0xFA, flagName: "Zero", flagMask: Zero},
		{value: 0x0A, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7A, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8e, 0x00, 0x00, 0x00, 0x2233: d.value})
		cpu.flags.Set(Carry, true)
		cpu.h = 0x22
		cpu.l = 0x33
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ADC M did not add A + (HL) + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Fuzz_ADC_A(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA4, flagName: "Parity", flagMask: Parity},
		{value: 0x80, flagName: "Zero", flagMask: Zero},
		{value: 0x0B, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7B, flagName: "Sign", flagMask: Sign},
		{value: 0xFC, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x8f, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, false)
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value+d.value {
			t.Errorf("ADC A did not add A + A + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADC A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_B(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x90, 0x00, 0x00, 0x00})
		cpu.b = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB B did not subtract A - B correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_C(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x91, 0x00, 0x00, 0x00})
		cpu.c = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB C did not subtract A - C correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_D(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x92, 0x00, 0x00, 0x00})
		cpu.d = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB D did not subtract A - D correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_E(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x93, 0x00, 0x00, 0x00})
		cpu.e = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB E did not subtract A - E correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_H(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x94, 0x00, 0x00, 0x00})
		cpu.h = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB H did not subtract A - H correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_L(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x95, 0x00, 0x00, 0x00})
		cpu.l = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB L did not subtract A - L correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SUB_M(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0A, flagName: "Parity", flagMask: Parity},
		{value: 0x05, flagName: "Zero", flagMask: Zero},
		{value: 0x14, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x89, flagName: "Sign", flagMask: Sign},
		{value: 0x04, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x96, 0x00, 0x00, 0x00, 0x2233: 0x05})
		cpu.h = 0x22
		cpu.l = 0x33
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUB M did not subtract A - (HL) correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Fuzz_SUB_A(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x05, flagName: "Zero", flagMask: Zero},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x97, 0x00, 0x00, 0x00})
		cpu.a = d.value

		cpu.Run()

		if cpu.a != 0 {
			t.Errorf("SUB A did not subtract A - A correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUB A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_B(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x98, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.b = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB B did not subtract A - B - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_C(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x99, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.c = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB C did not subtract A - C - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_D(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x9a, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.d = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB D did not subtract A - D - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_E(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x9b, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.e = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB E did not subtract A - E - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_H(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x9c, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.h = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB H did not subtract A - H - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_L(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x9d, 0x00, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.l = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB L did not subtract A - L - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_SBB_M(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x0B, flagName: "Parity", flagMask: Parity},
		{value: 0x06, flagName: "Zero", flagMask: Zero},
		{value: 0x15, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x8a, flagName: "Sign", flagMask: Sign},
		{value: 0x05, flagName: "Carry", flagMask: Carry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0x9e, 0x00, 0x00, 0x00, 0x2233: 0x05})
		cpu.flags.Set(Carry, true)
		cpu.h = 0x22
		cpu.l = 0x33
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBB M did not subtract A - M - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBB M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Test_SBB_AWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x9f, 0x00, 0x00, 0x00})
	cpu.flags.Set(Carry, true)
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0xFF {
		t.Errorf("SBB A did not subtract A - A - Carry correctly")
	}

	if cpu.flags.Get(Zero) {
		t.Errorf("SBB A did not set the Zero flag correctly")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("SBB A did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_SBB_AWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x9f, 0x00, 0x00, 0x00})
	cpu.flags.Set(Carry, false)
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x00 {
		t.Errorf("SBB A did not subtract A - A - Carry correctly")
	}

	if !cpu.flags.Get(Zero) {
		t.Errorf("SBB A did not set the Zero flag correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("SBB A did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_JNZ_ZeroFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc2, 0x88, 0xff})
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JNZ dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
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

	assertCycles(t, cpu, 10)
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

	assertCycles(t, cpu, 10)
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

	assertCycles(t, cpu, 17)
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

	assertCycles(t, cpu, 7)
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

		assertCycles(t, cpu, 7)
	})
}
