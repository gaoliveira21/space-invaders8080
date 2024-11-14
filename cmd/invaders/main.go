package main

import (
	"log"
	"os"
	"time"

	"github.com/gaoliveira21/intel8080-space-invaders/pkg/cpu"
)

func main() {
	log.Println("Starting Space Invaders...")
	log.Println("Reading ROM...")

	rom, err := os.ReadFile("roms/space-invaders/invaders")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	log.Printf("%d bytes loaded\n", len(rom))

	cpu := cpu.NewIntel8080()
	cpu.LoadProgram(rom, 0)

	// debugger := debug.NewDebugger(cpu)
	// debugger.Disassemble8080(rom)
	// debugger.DumpMemory()

	// api.Start()

	var instructionCycles uint
	lastFrame := time.Now()
	frameRate := time.Second / 60
	interruptType := 1

	for {
		if instructionCycles > 0 {
			instructionCycles--
			time.Sleep(time.Duration(0))
			continue
		}

		instructionCycles = cpu.Run()

		if time.Since(lastFrame) >= frameRate {
			if interruptType == 1 {
				interruptType = 2
			} else {
				interruptType = 1
			}

			// Draw

			if cpu.InterruptEnabled {
				cpu.Interrupt(interruptType)
			}

			lastFrame = time.Now()
		}
	}
}
