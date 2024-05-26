package main

import (
	"log"
	"os"

	"github.com/gaoliveira21/intel8080-space-invaders/debug"
)

func main() {
	log.Println("Starting Space Invaders...")
	log.Println("Reading ROM...")

	rom, err := os.ReadFile("roms/space-invaders/invaders.h")

	if err != nil {
		log.Fatalln("Cannot read ROM", err)
	}

	log.Printf("%d bytes loaded\n", len(rom))

	debug.Disassemble8080(rom)
}
