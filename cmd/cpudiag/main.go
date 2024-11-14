package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gaoliveira21/intel8080-space-invaders/pkg/cpu"
	"github.com/gaoliveira21/intel8080-space-invaders/pkg/io"
)

func onInput(cpu *cpu.Intel8080) {
	os.Exit(0)
}

func onOutput(cpu *cpu.Intel8080) {
	registers := cpu.GetRegisters()
	switch registers.C {
	// C = 0x02 signals printing the value of register E as an ASCII value
	case 0x02:
		fmt.Printf("%s", string(registers.E))
		// C = 0x09 signals printing the value of memory pointed to by DE until a '$' character is encountered
	case 0x09:
		addr := uint16(registers.D)<<8 | uint16(registers.E)
		for {
			c := cpu.ReadFromMemory(addr)
			if string(c) == "$" {
				break
			}
			fmt.Printf("%s", string(c))
			addr++
		}
	}
}

func main() {
	fmt.Println("Running a test ROM - roms/tests/TST8080.COM")
	rom, err := os.ReadFile("roms/tests/TST8080.COM")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	fmt.Printf("%d bytes loaded\n", len(rom))

	ioBus := io.NewIOBus()
	cpu := cpu.NewIntel8080(ioBus)
	cpu.LoadProgram(rom, 0x100)
	cpu.SetPC(0x100)

	in := byte(0xDB)  // IN pa
	out := byte(0xD3) // OUT
	ret := byte(0xC9) // RET

	cpu.WriteIntoMemory(0x0000, in)
	cpu.WriteIntoMemory(0x0005, out)
	cpu.WriteIntoMemory(0x0007, ret)

	cpu.SetInputListener(onInput)
	cpu.SetOutputListener(onOutput)

	for {
		cpu.Run()
	}
}
