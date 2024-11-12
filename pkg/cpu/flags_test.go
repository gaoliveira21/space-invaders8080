package cpu

import (
	"testing"
)

func TestSetFlagsToTrue(t *testing.T) {
	flags := &intel8080Flags{}

	flags.Set(Sign, true)
	flags.Set(Zero, true)
	flags.Set(AuxCarry, true)
	flags.Set(Parity, true)
	flags.Set(Carry, true)

	if flags.Get(Sign) != true {
		t.Errorf("Sign flag is not set")
	}

	if flags.Get(Zero) != true {
		t.Errorf("Zero flag is not set")
	}

	if flags.Get(AuxCarry) != true {
		t.Errorf("AuxCarry flag is not set")
	}

	if flags.Get(Parity) != true {
		t.Errorf("Parity flag is not set")
	}

	if flags.Get(Carry) != true {
		t.Errorf("Carry flag is not set")
	}
}

func TestSetFlagsToFalse(t *testing.T) {
	flags := &intel8080Flags{}

	flags.Set(Sign, false)
	flags.Set(Zero, false)
	flags.Set(AuxCarry, false)
	flags.Set(Parity, false)
	flags.Set(Carry, false)

	if flags.Get(Sign) != false {
		t.Errorf("Sign flag is not set")
	}

	if flags.Get(Zero) != false {
		t.Errorf("Zero flag is not set")
	}

	if flags.Get(AuxCarry) != false {
		t.Errorf("AuxCarry flag is not set")
	}

	if flags.Get(Parity) != false {
		t.Errorf("Parity flag is not set")
	}

	if flags.Get(Carry) != false {
		t.Errorf("Carry flag is not set")
	}
}

func TestSetFlagsToMixedValues(t *testing.T) {
	flags := &intel8080Flags{}

	flags.Set(Sign, true)
	flags.Set(Zero, false)
	flags.Set(AuxCarry, true)
	flags.Set(Parity, false)
	flags.Set(Carry, true)

	if flags.Get(Sign) != true {
		t.Errorf("Sign flag is not set")
	}

	if flags.Get(Zero) != false {
		t.Errorf("Zero flag is not set")
	}

	if flags.Get(AuxCarry) != true {
		t.Errorf("AuxCarry flag is not set")
	}

	if flags.Get(Parity) != false {
		t.Errorf("Parity flag is not set")
	}

	if flags.Get(Carry) != true {
		t.Errorf("Carry flag is not set")
	}
}
