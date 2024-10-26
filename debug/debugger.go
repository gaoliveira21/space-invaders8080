package debug

import (
	"log"
	"os"
)

type Cpu interface {
	GetMemory() []byte
}

type Debugger struct {
	cpu Cpu
}

func NewDebugger(c Cpu) *Debugger {
	return &Debugger{
		cpu: c,
	}
}

func (d *Debugger) DumpMemory() {
	mem := d.cpu.GetMemory()

	if _, err := os.Stat(".dump"); os.IsNotExist(err) {
		os.Mkdir(".dump", os.ModePerm)
	}

	err := os.WriteFile(".dump/memory", mem, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
