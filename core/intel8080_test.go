package core

import (
	"testing"
)

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
