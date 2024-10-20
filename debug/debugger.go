package debug

import (
	"log"
	"os"
)

type Debugger struct{}

func NewDebugger() *Debugger {
	return &Debugger{}
}

func (d *Debugger) DumpMemory(mem []byte) {
	if _, err := os.Stat(".dump"); os.IsNotExist(err) {
		os.Mkdir(".dump", os.ModePerm)
	}

	err := os.WriteFile(".dump/memory", mem, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
