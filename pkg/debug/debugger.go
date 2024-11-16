package debug

import (
	"fmt"
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
	fmt.Println("Dumping memory...")
	mem := d.cpu.GetMemory()

	if _, err := os.Stat(".dump"); os.IsNotExist(err) {
		os.Mkdir(".dump", os.ModePerm)
	}

	err := os.WriteFile(".dump/memory", mem, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
