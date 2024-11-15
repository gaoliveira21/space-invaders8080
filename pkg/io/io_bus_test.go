package io

import "testing"

func TestReadInput1(t *testing.T) {
	bus := NewIOBus(nil)
	bus.input1 = 0x05

	if bus.Read(0x1) != bus.input1 {
		t.Errorf("bus.Read did not return input 1 correctly")
	}
}

func TestReadInput2(t *testing.T) {
	bus := NewIOBus(nil)
	bus.input2 = 0x05

	if bus.Read(0x2) != bus.input2 {
		t.Errorf("bus.Read did not return input 2 correctly")
	}
}

func TestReadShiftRegisters(t *testing.T) {
	bus := NewIOBus(nil)
	bus.shiftH = 0xff
	bus.shiftL = 0xaa
	bus.offset = 2

	if bus.Read(0x3) != 0xFE {
		t.Errorf("bus.Read did not shift bytes correctly")
	}
}

func TestWriteOffset(t *testing.T) {
	bus := NewIOBus(nil)
	bus.Write(0x2, 0x62)

	if bus.offset != 0x2 {
		t.Errorf("bus.Write did not set offset correctly")
	}
}

func TestWriteShiftRegisters(t *testing.T) {
	bus := NewIOBus(nil)
	bus.shiftH = 0xff
	bus.Write(0x4, 0xaa)

	if bus.shiftL != 0xff || bus.shiftH != 0xaa {
		t.Errorf("bus.Write did not set shift registers correctly")
	}
}
