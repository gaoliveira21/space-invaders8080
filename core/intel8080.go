package core

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
}

func NewIntel8080() *Intel8080 {
	return &Intel8080{}
}
