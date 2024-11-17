package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gaoliveira21/intel8080-space-invaders/pkg/cpu"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/debug"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/io"
	"github.com/veandco/go-sdl2/sdl"
)

func onSignal(c chan os.Signal, running *bool, debugger *debug.Debugger) {
	for signal := range c {
		log.Printf("signal %s received\n", signal)
		if debugger != nil {
			debugger.Dump()
		}
		*running = false
	}
}

func main() {
	debugEnabled := flag.Bool("debug", false, "Run emulator in Debug Mode")
	audioDisabled := flag.Bool("sound-off", false, "Turn audio On/Off")

	flag.Parse()

	log.Println("Starting Space Invaders...")
	log.Println("Reading ROM...")

	rom, err := os.ReadFile("roms/space-invaders/invaders")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	log.Printf("%d bytes loaded\n", len(rom))

	var soundManager *io.SoundManager
	if !(*audioDisabled) {
		soundManager = io.NewSoundManager()
		defer soundManager.Cleanup()
	}

	ioBus := io.NewIOBus(soundManager)
	cpu := cpu.NewIntel8080(ioBus)
	cpu.LoadProgram(rom, 0)

	var debugger *debug.Debugger
	if *debugEnabled {
		debugger = debug.NewDebugger(cpu)
		go debugger.StartHttpServer()
	}

	running := true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go onSignal(c, &running, debugger)

	io.InitDisplay()
	defer io.DestroyDisplay()

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

			io.Draw(cpu.GetVRAM())

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
					case sdl.K_2:
						ioBus.OnInput(1, 1, pressed) // 2P start
					case sdl.K_1:
						ioBus.OnInput(1, 2, pressed) // 1P start
					case sdl.K_w:
						ioBus.OnInput(1, 4, pressed) // 1P shot
					case sdl.K_a:
						ioBus.OnInput(1, 5, pressed) // 1P left
					case sdl.K_d:
						ioBus.OnInput(1, 6, pressed) // 1P right
					case sdl.K_t:
						ioBus.OnInput(2, 2, pressed) // Tilt (Game over)
					case sdl.K_UP:
						ioBus.OnInput(2, 4, pressed) // 2P shot
					case sdl.K_LEFT:
						ioBus.OnInput(2, 5, pressed) // 2P left
					case sdl.K_RIGHT:
						ioBus.OnInput(2, 6, pressed) // 2P right
					}
				case *sdl.QuitEvent:
					running = false
				}
			}

			lastFrame = time.Now()
		}
	}
}
