package main

import (
	"log"
	"os"

	"github.com/gaoliveira21/intel8080-space-invaders/api"
	"github.com/gaoliveira21/intel8080-space-invaders/core"
	"github.com/gaoliveira21/intel8080-space-invaders/debug"
)

func main() {
	log.Println("Starting Space Invaders...")
	log.Println("Reading ROM...")

	rom, err := os.ReadFile("roms/space-invaders/invaders")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	log.Printf("%d bytes loaded\n", len(rom))

	debugger := debug.NewDebugger()
	debugger.Disassemble8080(rom)

	cpu := core.NewIntel8080(debugger)
	cpu.LoadProgram(rom)

	api.Start()
}
