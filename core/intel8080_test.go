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
