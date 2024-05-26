package debug

import (
	"fmt"
)

var Reset = "\033[0m"
var Green = "\033[32m"
var Cyan = "\033[36m"

func Disassemble8080(rom []byte) {
	fmt.Println("########## 8080 Opcode ##########")

	for i := 0; i < len(rom); i++ {
		opcode := rom[i]

		switch opcode {
		case 0x00:
			fmt.Printf("%.4X %.2X "+colorize("NOP", Green)+"\n", i, opcode)
		case 0x0F:
			fmt.Printf("%.4X %.2X "+colorize("RRC", Green)+"\n", i, opcode)

		case 0x21:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI H,", Green)+colorize(" #$%.4X\n", Cyan), i-2, opcode, lb, hb, addr)

		case 0x32:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("STA", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x35:
			fmt.Printf("%.4X %.2X "+colorize("DCR M", Green)+"\n", i, opcode)
		case 0x3A:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LDA", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x3E:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI A,", Green)+colorize(" #0x%.2X\n", Cyan), i-1, opcode, b, b)

		case 0xA7:
			fmt.Printf("%.4X %.2X "+colorize("ANA A", Green)+"\n", i, opcode)

		case 0xC3:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JMP", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xC5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH B", Green)+"\n", i, opcode)
		case 0xC6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("ADI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xCA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JZ", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xCD:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("CALL", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)

		case 0xD5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH D", Green)+"\n", i, opcode)
		case 0xDA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xDB:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("IN ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)

		case 0xE5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH H", Green)+"\n", i, opcode)

		case 0xF5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH PSW", Green)+"\n", i, opcode)
		case 0xFE:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("CPI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)

		default:
			fmt.Printf("Unknown Opcode %.2X\n", opcode)
		}
	}
	fmt.Println("#################################")
}

func colorize(instruction string, color string) string {
	return color + instruction + Reset
}

func getAddr(index *int, rom []byte) (uint16, uint16, uint16) {
	*index++
	lb := uint16(rom[*index])
	*index++
	hb := uint16(rom[*index])

	addr := (hb << 8) | lb

	return lb, hb, addr
}
