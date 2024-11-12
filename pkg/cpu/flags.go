/**
 * 7 	6 	5 	4 	3 	2 	1 	0
 * S 	Z 	0 	A 	0 	P 	1 	C
 *
 * S - Sign Flag
 * Z - Zero Flag
 * 0 - Not used, always zero
 * A - also called AC, Auxiliary Carry Flag
 * 0 - Not used, always zero
 * P - Parity Flag
 * 1 - Not used, always one
 * C - Carry Flag
 */

package cpu

type Flag = byte

const (
	Sign     Flag = 1 << 7
	Zero     Flag = 1 << 6
	AuxCarry Flag = 1 << 4
	Parity   Flag = 1 << 2
	Carry    Flag = 1 << 0
)

type intel8080Flags struct {
	value byte
}

func (flags *intel8080Flags) Set(flag Flag, value bool) {
	if value {
		flags.value |= flag
	} else {
		flags.value &= ^flag
	}
}

func (flags *intel8080Flags) Get(flag Flag) bool {
	return flags.value&flag != 0
}
