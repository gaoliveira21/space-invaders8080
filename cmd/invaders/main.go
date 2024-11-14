package main

import (
	"log"
	"os"
	"time"

	"github.com/gaoliveira21/intel8080-space-invaders/pkg/cpu"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/display"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	log.Println("Starting Space Invaders...")
	log.Println("Reading ROM...")

	rom, err := os.ReadFile("roms/space-invaders/invaders")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	log.Printf("%d bytes loaded\n", len(rom))

	running := true

	cpu := cpu.NewIntel8080()
	cpu.LoadProgram(rom, 0)

	display.Init()
	defer display.Destroy()

	var instructionCycles uint
	lastFrame := time.Now()
	frameRate := time.Second / 120
	interruptType := 1

	for running {
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

			display.Draw(cpu.GetVRAM())

			if cpu.InterruptEnabled {
				cpu.Interrupt(interruptType)
			}

			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					running = false
				}
			}

			lastFrame = time.Now()
		}
	}
}
