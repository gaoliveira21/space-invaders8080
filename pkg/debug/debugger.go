package debug

import (
	"encoding/json"
	"log"
	"os"
)

type Cpu interface {
	GetMemory() []byte
	GetRegisters() map[string]byte
	GetPointers() map[string]uint16
}

type CpuState struct {
	Registers map[string]byte   `json:"registers"`
	Pointers  map[string]uint16 `json:"pointers"`
}

type Debugger struct {
	cpu Cpu
}

func NewDebugger(c Cpu) *Debugger {
	return &Debugger{
		cpu: c,
	}
}

func (d *Debugger) Dump() {
	d.dumpMemory()
	d.dumpCpuState()
}

func (d *Debugger) createDumpFolderIfNotExists() {
	if _, err := os.Stat(".dump"); os.IsNotExist(err) {
		os.Mkdir(".dump", os.ModePerm)
	}
}

func (d *Debugger) dumpMemory() {
	log.Println("Dumping memory...")
	mem := d.cpu.GetMemory()

	d.createDumpFolderIfNotExists()

	err := os.WriteFile(".dump/memory", mem, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Debugger) dumpCpuState() {
	log.Println("Dumping cpu state...")
	d.createDumpFolderIfNotExists()

	state := &CpuState{
		Registers: d.cpu.GetRegisters(),
		Pointers:  d.cpu.GetPointers(),
	}

	stateJson, _ := json.Marshal(state)
	err := os.WriteFile(".dump/cpu_state.json", stateJson, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
