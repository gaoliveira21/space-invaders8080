package main

import (
	"log"
	"os"
	"time"

	"github.com/gaoliveira21/intel8080-space-invaders/pkg/cpu"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/display"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/io"
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

	ioBus := io.NewIOBus()
	cpu := cpu.NewIntel8080(ioBus)
	cpu.LoadProgram(rom, 0)

	display.Init()
	defer display.Destroy()

	var instructionCycles uint
	lastFrame := time.Now()
	frameRate := time.Second / 120
	interruptType := 1

	running := true

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
				switch t := event.(type) {
				case *sdl.KeyboardEvent:
					pressed := false
					if t.Type == sdl.KEYDOWN {
						pressed = true
					} else if t.Type == sdl.KEYUP {
						pressed = false
					}

					switch t.Keysym.Sym {
					case sdl.K_c:
						ioBus.OnInput(1, 0, pressed) // Coin
					case sdl.K_SPACE:
						ioBus.OnInput(1, 2, pressed) // 1P start
					case sdl.K_z:
						ioBus.OnInput(1, 4, pressed) // 1P shot
					case sdl.K_LEFT:
						ioBus.OnInput(1, 5, pressed) // 1P left
					case sdl.K_RIGHT:
						ioBus.OnInput(1, 6, pressed) // 1P right
					case sdl.K_t:
						ioBus.OnInput(2, 2, pressed) // Tilt (Game over)
					}
				case *sdl.QuitEvent:
					running = false
				}
			}

			lastFrame = time.Now()
		}
	}
}
