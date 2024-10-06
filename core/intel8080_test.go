package core

import (
	"testing"
)

type flagDataTest struct {
	value    byte
	flagName string
	flagMask byte
}

func TestLoadProgram(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x00, 0x01, 0x02, 0x03}
	cpu.LoadProgram(program)

	for i, v := range program {
		if cpu.memory[i] != v {
			t.Errorf("LoadProgram did not load the program correctly")
		}
	}
}

func Test_LXI_B(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x01, 0x02, 0x03}
	cpu.LoadProgram(program)

	cpu.Run()

	if cpu.c != 0x02 || cpu.b != 0x03 {
		t.Errorf("LXI B did not load the program correctly")
	}
}

func Test_STAX_B(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x02, 0x01}
	cpu.LoadProgram(program)

	cpu.b = 0x03
	cpu.c = 0x01
	cpu.a = 0x08

	cpu.Run()

	if cpu.memory[0x0301] != 0x08 {
		t.Errorf("STAX B did not store the program correctly")
	}
}

func Test_INX_B(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x03, 0x01}
	cpu.LoadProgram(program)

	cpu.b = 0x03
	cpu.c = 0x01

	cpu.Run()

	if cpu.c != 0x02 || cpu.b != 0x03 {
		t.Errorf("INX B did not increment the program correctly")
	}
}

func Test_INR_B(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x04, 0x01}
	cpu.LoadProgram(program)

	cpu.b = 0x03

	cpu.Run()

	if cpu.b != 0x04 {
		t.Errorf("INR B did not increment the program correctly")
	}

	if cpu.flags.Get(Parity) {
		t.Errorf("INR B did not set the parity flag correctly")
	}

	if cpu.flags.Get(Zero) {
		t.Errorf("INR B did not set the zero flag correctly")
	}

	if cpu.flags.Get(Sign) {
		t.Errorf("INR B did not set the sign flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("INR B did not set the auxiliary carry flag correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("INR B did not set the carry flag correctly")
	}
}

func Fuzz_INR_B_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xA9, flagName: "Parity", flagMask: Parity},
		{value: 0xFF, flagName: "Zero", flagMask: Zero},
		{value: 0x2F, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x7F, flagName: "Sign", flagMask: Sign},
	}

	f.Add(0)
	f.Add(1)
	f.Add(2)
	f.Add(3)

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := NewIntel8080()

		program := []byte{0x04, 0x01}
		cpu.LoadProgram(program)

		cpu.b = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR B did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_DCR_B(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x05, 0x01}
	cpu.LoadProgram(program)

	cpu.b = 0x05

	cpu.Run()

	if cpu.b != 0x04 {
		t.Errorf("INR B did not increment the program correctly")
	}

	if cpu.flags.Get(Parity) {
		t.Errorf("INR B did not set the parity flag correctly")
	}

	if cpu.flags.Get(Zero) {
		t.Errorf("INR B did not set the zero flag correctly")
	}

	if cpu.flags.Get(Sign) {
		t.Errorf("INR B did not set the sign flag correctly")
	}

	if cpu.flags.Get(AuxCarry) {
		t.Errorf("INR B did not set the auxiliary carry flag correctly")
	}

	if cpu.flags.Get(Carry) {
		t.Errorf("INR B did not set the carry flag correctly")
	}
}

func Fuzz_DCR_B_Flags(f *testing.F) {
	tData := []flagDataTest{
		{value: 0xAB, flagName: "Parity", flagMask: Parity},
		{value: 0x01, flagName: "Zero", flagMask: Zero},
		{value: 0x31, flagName: "AuxCarry", flagMask: AuxCarry},
		{value: 0x81, flagName: "Sign", flagMask: Sign},
	}

	f.Add(0)
	f.Add(1)
	f.Add(2)
	f.Add(3)

	f.Fuzz(func(t *testing.T, i int) {
		d := tData[i]
		cpu := NewIntel8080()

		program := []byte{0x05, 0x01}
		cpu.LoadProgram(program)

		cpu.b = d.value

		cpu.Run()

		if !cpu.flags.Get(d.flagMask) {
			t.Errorf("INR B did not set the %s flag correctly", d.flagName)
		}
	})
}

func Test_MVI_B_Flags(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x06, 0x42}
	cpu.LoadProgram(program)

	cpu.Run()

	if cpu.b != 0x42 {
		t.Errorf("MVI B did not load the correct value to register")
	}
}

func Test_RLC_Flags(t *testing.T) {
	cpu := NewIntel8080()

	program := []byte{0x07, 0x01}
	cpu.LoadProgram(program)

	cpu.a = 0x80

	cpu.Run()

	if cpu.a != 0x01 {
		t.Errorf("RLC did not set the correct value to register A")
	}

	if !cpu.flags.Get(Carry) {
		t.Errorf("RLC did not set the carry flag correctly")
	}
}
