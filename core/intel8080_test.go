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
	cpu.LoadProgram(p, 0)

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

func Test_ANA_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa0, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.b = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA B did not A & B correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA B did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa1, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.c = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA C did not A & C correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA C did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa2, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.d = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA D did not A & D correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA D did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_E(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa3, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.e = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA E did not A & E correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA E did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa4, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.h = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA H did not A & H correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA H did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_L(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa5, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.l = 0x09

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA L did not A & L correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA L did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ANA_M(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa6, 0x00, 0x00, 0x00, 0x2233: 0x09})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("ANA M did not A & (HL) correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA M did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_ANA_A(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa7, 0x00, 0x00, 0x00})
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x06 {
		t.Errorf("ANA A did not A & A correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ANA A did not set the Carry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa8, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.b = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA B did not A ^ B correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA B did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA B did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xa9, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.c = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA C did not A ^ C correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA C did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA C did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xaa, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.d = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA D did not A ^ D correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA D did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA D did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_E(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xab, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.e = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA E did not A ^ E correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA E did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA E did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xac, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.h = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA H did not A ^ H correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA H did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA H did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_L(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xad, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.l = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA L did not A ^ L correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA L did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA L did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_XRA_M(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xae, 0x00, 0x00, 0x00, 0x2233: 0x09})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRA L did not A ^ L correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA L did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA L did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_XRA_A(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xaf, 0x00, 0x00, 0x00})
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0 {
		t.Errorf("XRA L did not A ^ L correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRA L did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRA L did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb0, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.b = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA B did not A | B correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA B did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA B did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_C(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb1, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.c = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA C did not A | C correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA C did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA C did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb2, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.d = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA D did not A | D correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA D did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA D did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_E(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb3, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.e = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA E did not A | E correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA E did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA E did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb4, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.h = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA H did not A | H correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA H did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA H did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_L(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb5, 0x00, 0x00, 0x00})
	cpu.a = 0x06
	cpu.l = 0x09

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA L did not A | L correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA L did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA L did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_ORA_M(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb6, 0x00, 0x00, 0x00, 0x2233: 0x09})
	cpu.h = 0x22
	cpu.l = 0x33
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORA M did not A | (HL) correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA M did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA M did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_ORA_A(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xb7, 0x00, 0x00, 0x00})
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x06 {
		t.Errorf("ORA A did not A | A correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORA A did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORA A did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 4)
}

func Fuzz_CMP_B(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xb8, 0x00, 0x00, 0x00})
		cpu.b = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP B changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP B did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_C(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xb9, 0x00, 0x00, 0x00})
		cpu.c = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP C changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP C did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_D(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xba, 0x00, 0x00, 0x00})
		cpu.d = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP D changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP D did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_E(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xbb, 0x00, 0x00, 0x00})
		cpu.e = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP E changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP E did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_H(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xbc, 0x00, 0x00, 0x00})
		cpu.h = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP H changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP H did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_L(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xbd, 0x00, 0x00, 0x00})
		cpu.l = 0x05
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP L changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP L did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Fuzz_CMP_M(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xbe, 0x00, 0x00, 0x00, 0x2233: 0x05})
		cpu.h = 0x22
		cpu.l = 0x33
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP M changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP M did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Fuzz_CMP_A(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x05, flagName: "Zero", flagMask: Zero},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0xbf, 0x00, 0x00, 0x00})
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value {
			t.Errorf("CMP A changed A register")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("CMP A did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 4)
	})
}

func Test_RNZWithZeroUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RNZ did not set PC correctly when Zero flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RNZ did not set SP correctly when Zero flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RNZWithZeroSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RNZ modified PC when Zero flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RNZ modified SP when Zero flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_POP_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc1, 0x34, 0x12})
	cpu.sp = 1

	cpu.Run()

	if cpu.b != 0x12 {
		t.Errorf("POP B did not set B register correctly")
	}
	if cpu.c != 0x34 {
		t.Errorf("POP B did not set C register correctly")
	}
	if cpu.sp != 3 {
		t.Errorf("POP B did not set SP correctly")
	}
	assertCycles(t, cpu, 10)
}

func Test_JNZ_ZeroFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc2, 0x88, 0xff})
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JNZ dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JNZ_ZeroFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc2, 0x88, 0xff, 0x01})
	cpu.flags.Set(Zero, true)

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

func Test_CNZWithZeroUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CNZ did not set PC correctly when Zero flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CNZ did not set SP correctly when Zero flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CNZ did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CNZWithZeroSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CNZ did not increment PC correctly when Zero flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CNZ modified SP when Zero flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_PUSH_B(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc5, 0x0, 0x0, 0x0, 0x0})
	cpu.sp = 4
	cpu.b = 0x12
	cpu.c = 0x34

	cpu.Run()

	if cpu.memory[2] != 0x34 {
		t.Errorf("PUSH B did not store C correctly")
	}
	if cpu.memory[3] != 0x12 {
		t.Errorf("PUSH B did not store B correctly")
	}
	if cpu.sp != 2 {
		t.Errorf("PUSH B did not decrement SP correctly")
	}
	assertCycles(t, cpu, 11)
}

func Fuzz_ADI(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xc6, d.value, 0x00, 0x00})
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+5 {
			t.Errorf("ADI did not add A + B correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ADI did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Test_RST_0(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xc7})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0000 {
		t.Errorf("RST 0 did not set PC to 0x0000, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 0 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 0 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RZWithZeroSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc8, 0x34, 0x12, 0x00})
	cpu.flags.Set(Zero, true)
	cpu.sp = 1

	cpu.Run()

	if cpu.sp != 3 {
		t.Errorf("RZ dit not set SP correctly")
	}

	if cpu.pc != 0x1234 {
		t.Errorf("RZ dit not set PC correctly")
	}

	assertCycles(t, cpu, 11)
}

func Test_RZWithZeroUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc8, 0x34, 0x12, 0x00})
	cpu.flags.Set(Zero, false)
	cpu.sp = 1

	cpu.Run()

	if cpu.sp != 1 {
		t.Errorf("RZ dit not set SP correctly")
	}

	if cpu.pc != 0x1 {
		t.Errorf("RZ dit not set PC correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_RET(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xc9, 0x34, 0x12, 0x00})
	cpu.sp = 1

	cpu.Run()

	if cpu.sp != 3 {
		t.Errorf("RET dit not set SP correctly")
	}

	if cpu.pc != 0x1234 {
		t.Errorf("RET dit not set PC correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JZ_ZeroFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xca, 0x88, 0xff, 0x01})
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JNZ dit not set pc correctly")
	}
}

func Test_JZ_ZeroFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xca, 0x88, 0xff})
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JNZ dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_CZWithZeroSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xcc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Zero, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CZ did not set PC correctly when Zero flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CZ did not set SP correctly when Zero flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CZ did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CZWithZeroUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xcc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Zero, false)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CZ did not increment PC correctly when Zero flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CZ modified SP when Zero flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_CALL(t *testing.T) {
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

func Fuzz_ACI(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xce, d.value, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.a = 0x05

		cpu.Run()

		if cpu.a != d.value+6 {
			t.Errorf("ACI did not add A + Data + Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("ACI did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Test_RST_1(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xcf})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0008 {
		t.Errorf("RST 1 did not set PC to 0x0008, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 1 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 1 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RNCWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RNC did not set PC correctly when Carry flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RNC did not set SP correctly when Carry flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RNCWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RNC modified PC when Carry flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RNC modified SP when Carry flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_POP_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd1, 0x34, 0x12})
	cpu.sp = 1

	cpu.Run()

	if cpu.d != 0x12 {
		t.Errorf("POP D did not set D register correctly")
	}
	if cpu.e != 0x34 {
		t.Errorf("POP D did not set E register correctly")
	}
	if cpu.sp != 3 {
		t.Errorf("POP D did not set SP correctly")
	}
	assertCycles(t, cpu, 10)
}

func Test_JNC_CarryFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd2, 0x88, 0xff})
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JNC dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JNC_CarryFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd2, 0x88, 0xff, 0x01})
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JNC dit not set pc correctly")
	}
}

func Test_OUT(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd3, 0x88, 0x01, 0x01})

	cpu.Run()

	if cpu.pc != 0x0002 {
		t.Errorf("OUT dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_CNCWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CNC did not set PC correctly when Carry flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CNC did not set SP correctly when Carry flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CNC did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CNCWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CNC did not increment PC correctly when Carry flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CNC modified SP when Carry flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_PUSH_D(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd5, 0x0, 0x0, 0x0, 0x0})
	cpu.sp = 4
	cpu.d = 0x12
	cpu.e = 0x34

	cpu.Run()

	if cpu.memory[2] != 0x34 {
		t.Errorf("PUSH D did not store E correctly")
	}
	if cpu.memory[3] != 0x12 {
		t.Errorf("PUSH D did not store D correctly")
	}
	if cpu.sp != 2 {
		t.Errorf("PUSH D did not decrement SP correctly")
	}
	assertCycles(t, cpu, 11)
}

func Fuzz_SUI(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xd6, 0x05, 0x00, 0x00})
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-5 {
			t.Errorf("SUI did not subtract A - Data correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SUI did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Test_RST_2(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xd7})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0010 {
		t.Errorf("RST 2 did not set PC to 0x0010, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 2 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 2 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RCWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RC did not set PC correctly when Carry flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RC did not set SP correctly when Carry flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RCWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xd8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RC modified PC when Carry flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RC modified SP when Carry flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_JC_CarryFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xda, 0x88, 0xff})
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JC dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JC_CarryFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xda, 0x88, 0xff, 0x01})
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JC dit not set pc correctly")
	}
}

func Test_IN(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xdb, 0x88, 0x01, 0x01})

	cpu.Run()

	if cpu.pc != 0x0002 {
		t.Errorf("IN dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_CCWithCarrySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xdc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CC did not set PC correctly when Carry flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CC did not set SP correctly when Carry flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CC did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CCWithCarryUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xdc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Carry, false)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CC did not increment PC correctly when Carry flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CC modified SP when Carry flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Fuzz_SBI(f *testing.F) {
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
		cpu := createCPUWithProgramLoaded([]byte{0xde, 0x05, 0x00, 0x00})
		cpu.flags.Set(Carry, true)
		cpu.a = d.value

		cpu.Run()

		if cpu.a != d.value-6 {
			t.Errorf("SBI did not subtract A - Data - Carry correctly")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("SBI did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 7)
	})
}

func Test_RST_3(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xdf})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0018 {
		t.Errorf("RST 3 did not set PC to 0x0018, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 3 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 3 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RPOWithParityUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RPO did not set PC correctly when Parity flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RPO did not set SP correctly when Parity flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RPOWithParitySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RPO modified PC when Parity flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RPO modified SP when Parity flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_POP_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe1, 0x34, 0x12})
	cpu.sp = 1

	cpu.Run()

	if cpu.h != 0x12 {
		t.Errorf("POP H did not set H register correctly")
	}
	if cpu.l != 0x34 {
		t.Errorf("POP H did not set L register correctly")
	}
	if cpu.sp != 3 {
		t.Errorf("POP H did not set SP correctly")
	}
	assertCycles(t, cpu, 10)
}

func Test_JPO_ParityFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe2, 0x88, 0xff})
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JPO dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JPO_ParityFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe2, 0x88, 0xff, 0x01})
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JPO dit not set pc correctly")
	}
}

func Test_XTHL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe3, 0x01, 0x88, 0xff, 0x01})
	cpu.sp = 2
	cpu.h = 0x99
	cpu.l = 0x05

	cpu.Run()

	if cpu.h != 0xff {
		t.Errorf("XTHL dit not set H register correctly")
	}

	if cpu.l != 0x88 {
		t.Errorf("XTHL dit not set L register correctly")
	}

	if cpu.memory[cpu.sp] != 0x05 || cpu.memory[cpu.sp+1] != 0x99 {
		t.Errorf("XTHL dit not write HL values into memory correctly")
	}
	assertCycles(t, cpu, 18)
}

func Test_CPOWithParityUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CPO did not set PC correctly when Parity flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CPO did not set SP correctly when Parity flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CPO did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CPOWithParitySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CPO did not increment PC correctly when Parity flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CPO modified SP when Parity flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_PUSH_H(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe5, 0x0, 0x0, 0x0, 0x0})
	cpu.sp = 4
	cpu.h = 0x12
	cpu.l = 0x34

	cpu.Run()

	if cpu.memory[2] != 0x34 {
		t.Errorf("PUSH H did not store L correctly")
	}
	if cpu.memory[3] != 0x12 {
		t.Errorf("PUSH H did not store H correctly")
	}
	if cpu.sp != 2 {
		t.Errorf("PUSH H did not decrement SP correctly")
	}
	assertCycles(t, cpu, 11)
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

func Test_RST_4(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xe7})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0020 {
		t.Errorf("RST 4 did not set PC to 0x0020, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 4 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 4 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RPEWithParitySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RPE did not set PC correctly when Parity flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RPE did not set SP correctly when Parity flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RPEWithParityUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RPE modified PC when Parity flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RPE modified SP when Parity flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_PCHL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xe9, 0x00, 0x00})
	cpu.h = 0x33
	cpu.l = 0x0F

	cpu.Run()

	if cpu.pc != 0x330F {
		t.Errorf("PCHL did not set PC correctly")
	}
	assertCycles(t, cpu, 5)
}

func Test_JPE_ParityFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xea, 0x88, 0xff})
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JPE dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JPE_ParityFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xea, 0x88, 0xff, 0x01})
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JPE dit not set pc correctly")
	}
}

func Test_XCHG(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xeb, 0x00, 0x00})
	cpu.h = 0x11
	cpu.l = 0x22
	cpu.d = 0x33
	cpu.e = 0x44

	cpu.Run()

	if cpu.h != 0x33 || cpu.l != 0x44 {
		t.Errorf("XCHG dit not swap HL register pair correctly")
	}

	if cpu.d != 0x11 || cpu.e != 0x22 {
		t.Errorf("XCHG dit not swap DE register pair correctly")
	}

	assertCycles(t, cpu, 5)
}

func Test_CPEWithParitySet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xec, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Parity, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CPE did not set PC correctly when Parity flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CPE did not set SP correctly when Parity flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CPE did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CPEWithParityUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xec, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Parity, false)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CPE did not increment PC correctly when Parity flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CPE modified SP when Parity flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_XRI(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xee, 0x09, 0x00, 0x00})
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("XRI did not A ^ Data correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("XRI did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("XRI did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_RST_5(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xef})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0028 {
		t.Errorf("RST 5 did not set PC to 0x0028, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 5 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 5 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RPWithSignUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RP did not set PC correctly when Sign flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RP did not set SP correctly when Sign flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RPWithSignSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf0, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RP modified PC when Sign flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RP modified SP when Sign flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Fuzz_POP_PSW_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0x04, flagName: "Parity", flagMask: Parity},
		{value: 0x40, flagName: "Zero", flagMask: Zero},
		{value: 0x80, flagName: "Sign", flagMask: Sign},
		{value: 0x01, flagName: "Carry", flagMask: Carry},
		{value: 0x10, flagName: "AuxCarry", flagMask: AuxCarry},
	}

	for i := range tData {
		f.Add(i)
	}

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := createCPUWithProgramLoaded([]byte{0xf1, d.value, 0x0A, 0x0})
		cpu.sp = 1

		cpu.Run()

		if cpu.sp != 3 {
			t.Errorf("POP PSW did not set SP correcty")
		}

		if cpu.a != 0x0A {
			t.Errorf("POP PSW did not set A register correcty")
		}

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("POP PSW did not set the %s flag correctly", d.flagName)
		}

		assertCycles(t, cpu, 10)
	})
}

func Test_JP_SignFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf2, 0x88, 0xff})
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JP dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JP_SignFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf2, 0x88, 0xff, 0x01})
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JP dit not set pc correctly")
	}
}

func Test_DI(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf3, 0x00, 0x00})
	cpu.InterruptEnabled = true

	cpu.Run()

	if cpu.InterruptEnabled {
		t.Errorf("DI dit not disable Interrupts correctly")
	}

	assertCycles(t, cpu, 4)
}

func Test_CPWithSignUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CP did not set PC correctly when Sign flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CP did not set SP correctly when Sign flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CP did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CPWithSignSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf4, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CP did not increment PC correctly when Sign flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CP modified SP when Sign flag was set")
	}
	assertCycles(t, cpu, 11)
}

func Test_PUSH_PSW(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf5, 0x00, 0x00, 0x00})
	cpu.sp = 3
	cpu.a = 0x5
	cpu.flags.Set(Sign, true)
	cpu.flags.Set(Zero, true)
	cpu.flags.Set(Parity, true)
	cpu.flags.Set(AuxCarry, true)
	cpu.flags.Set(Carry, true)

	cpu.Run()

	if cpu.sp != 1 {
		t.Errorf("PUSH PSW did not set sp correctly")
	}

	if cpu.memory[2] != 0x5 {
		t.Errorf("PUSH PSW did not store A register correctly")
	}

	if cpu.memory[1] != 0xD7 {
		t.Errorf("PUSH PSW did not store program status correctly")
	}

	assertCycles(t, cpu, 11)
}

func Test_ORI(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf6, 0x09, 0x00, 0x00})
	cpu.a = 0x06

	cpu.Run()

	if cpu.a != 0x0F {
		t.Errorf("ORI did not A | Data correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("ORI did not set the Carry flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("ORI did not set the AuxCarry flag correctly")
	}

	assertCycles(t, cpu, 7)
}

func Test_RST_6(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xf7})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0030 {
		t.Errorf("RST 6 did not set PC to 0x0030, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 6 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 6 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}

func Test_RMWithSignSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("RM did not set PC correctly when Sign flag was not set")
	}
	if cpu.sp != 3 {
		t.Errorf("RM did not set SP correctly when Sign flag was not set")
	}
	assertCycles(t, cpu, 11)
}

func Test_RMWithSignUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf8, 0x34, 0x12})
	cpu.sp = 1
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 1 {
		t.Errorf("RM modified PC when Sign flag was set")
	}
	if cpu.sp != 1 {
		t.Errorf("RM modified SP when Sign flag was set")
	}
	assertCycles(t, cpu, 5)
}

func Test_SPHL(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xf9, 0x00, 0x00})
	cpu.h = 0x33
	cpu.l = 0x0F

	cpu.Run()

	if cpu.sp != 0x330F {
		t.Errorf("SPHL did not set PC correctly")
	}
	assertCycles(t, cpu, 5)
}

func Test_JM_SignFlagSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xfa, 0x88, 0xff})
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 0xff88 {
		t.Errorf("JM dit not set pc correctly")
	}

	assertCycles(t, cpu, 10)
}

func Test_JM_SignFlagUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xfa, 0x88, 0xff, 0x01})
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 0x0003 {
		t.Errorf("JM dit not set pc correctly")
	}
}

func Test_EI(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xfb, 0x00, 0x00})
	cpu.enableInterruptDeferred = false
	cpu.InterruptEnabled = false

	cpu.Run()

	if !cpu.enableInterruptDeferred {
		t.Errorf("EI dit not deferred interrupt correctly")
	}

	assertCycles(t, cpu, 4)

	cpu.Run()

	if !cpu.InterruptEnabled {
		t.Errorf("EI dit not enable interrupt correctly")
	}
}

func Test_CMWithSignSet(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xfc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Sign, true)

	cpu.Run()

	if cpu.pc != 0x1234 {
		t.Errorf("CM did not set PC correctly when Sign flag was not set")
	}
	if cpu.sp != 1 {
		t.Errorf("CM did not set SP correctly when Sign flag was not set")
	}

	if cpu.memory[1] != 0x03 || cpu.memory[2] != 0x00 {
		t.Errorf("CM did not store return address correctly")
	}
	assertCycles(t, cpu, 17)
}

func Test_CMWithSignUnset(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0xfc, 0x34, 0x12, 0x00})
	cpu.sp = 3
	cpu.flags.Set(Sign, false)

	cpu.Run()

	if cpu.pc != 3 {
		t.Errorf("CM did not increment PC correctly when Sign flag was set")
	}
	if cpu.sp != 3 {
		t.Errorf("CM modified SP when Sign flag was set")
	}
	assertCycles(t, cpu, 11)
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

func Test_RST_7(t *testing.T) {
	cpu := createCPUWithProgramLoaded([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x1233: 0xff})
	cpu.sp = 0x04
	cpu.pc = 0x1233

	cpu.Run()

	if cpu.pc != 0x0038 {
		t.Errorf("RST 7 did not set PC to 0x0038, got 0x%04x", cpu.pc)
	}

	if cpu.memory[cpu.sp] != 0x34 || cpu.memory[cpu.sp+1] != 0x12 {
		t.Errorf("RST 7 did not save return address correctly")
	}

	if cpu.sp != 0x02 {
		t.Errorf("RST 7 did not adjust SP correctly, got 0x%04x", cpu.sp)
	}

	assertCycles(t, cpu, 11)
}
