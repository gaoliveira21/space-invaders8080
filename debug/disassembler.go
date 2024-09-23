package debug

import (
	"fmt"
)

var Reset = "\033[0m"
var Green = "\033[32m"
var Cyan = "\033[36m"

func Disassemble8080(rom []byte) {
	fmt.Println("########## 8080 Opcode ##########")

L:
	for i := 0; i < len(rom); i++ {
		opcode := rom[i]

		switch opcode {
		case 0x00:
			fmt.Printf("%.4X %.2X "+colorize("NOP", Green)+"\n", i, opcode)
		case 0x01:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI B,", Green)+colorize(" #$%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x05:
			fmt.Printf("%.4X %.2X "+colorize("DCR B", Green)+"\n", i, opcode)
		case 0x06:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI B,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x07:
			fmt.Printf("%.4X %.2X "+colorize("RLC", Green)+"\n", i, opcode)
		case 0x0F:
			fmt.Printf("%.4X %.2X "+colorize("RRC", Green)+"\n", i, opcode)

		case 0x16:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI D,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x19:
			fmt.Printf("%.4X %.2X "+colorize("DAD D", Green)+"\n", i, opcode)

		case 0x21:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI H,", Green)+colorize(" #$%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x22:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("SHLD", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x23:
			fmt.Printf("%.4X %.2X "+colorize("INX H", Green)+"\n", i, opcode)
		case 0x27:
			fmt.Printf("%.4X %.2X "+colorize("DAA", Green)+"\n", i, opcode)
		case 0x2A:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LHLD", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x2B:
			fmt.Printf("%.4X %.2X "+colorize("DCX H", Green)+"\n", i, opcode)
		case 0x2F:
			fmt.Printf("%.4X %.2X "+colorize("CMA", Green)+"\n", i, opcode)

		case 0x32:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("STA", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x35:
			fmt.Printf("%.4X %.2X "+colorize("DCR M", Green)+"\n", i, opcode)
		case 0x3A:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LDA", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x3C:
			fmt.Printf("%.4X %.2X "+colorize("INR A", Green)+"\n", i, opcode)
		case 0x3D:
			fmt.Printf("%.4X %.2X "+colorize("DCR A", Green)+"\n", i, opcode)
		case 0x3E:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI A,", Green)+colorize(" #0x%.2X\n", Cyan), i-1, opcode, b, b)

		case 0x46:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,M", Green)+"\n", i, opcode)
		case 0x47:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,A", Green)+"\n", i, opcode)
		case 0x4F:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,A", Green)+"\n", i, opcode)

		case 0x56:
			fmt.Printf("%.4X %.2X "+colorize("MOV D,M", Green)+"\n", i, opcode)
		case 0x57:
			fmt.Printf("%.4X %.2X "+colorize("MOV D,A", Green)+"\n", i, opcode)
		case 0x5F:
			fmt.Printf("%.4X %.2X "+colorize("MOV E,A", Green)+"\n", i, opcode)

		case 0x60:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,B", Green)+"\n", i, opcode)
		case 0x61:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,C", Green)+"\n", i, opcode)
		case 0x62:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,D", Green)+"\n", i, opcode)
		case 0x63:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,E", Green)+"\n", i, opcode)
		case 0x64:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,H", Green)+"\n", i, opcode)
		case 0x65:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,L", Green)+"\n", i, opcode)
		case 0x66:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,M", Green)+"\n", i, opcode)
		case 0x67:
			fmt.Printf("%.4X %.2X "+colorize("MOV H,A", Green)+"\n", i, opcode)
		case 0x6F:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,A", Green)+"\n", i, opcode)

		case 0x78:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,B", Green)+"\n", i, opcode)
		case 0x79:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,C", Green)+"\n", i, opcode)
		case 0x7A:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,D", Green)+"\n", i, opcode)
		case 0x7B:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,E", Green)+"\n", i, opcode)
		case 0x7C:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,H", Green)+"\n", i, opcode)
		case 0x7D:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,L", Green)+"\n", i, opcode)
		case 0x7E:
			fmt.Printf("%.4X %.2X "+colorize("MOV A,M", Green)+"\n", i, opcode)

		case 0xA7:
			fmt.Printf("%.4X %.2X "+colorize("ANA A", Green)+"\n", i, opcode)
		case 0xAF:
			fmt.Printf("%.4X %.2X "+colorize("XRA A", Green)+"\n", i, opcode)

		case 0xC0:
			fmt.Printf("%.4X %.2X "+colorize("RNZ", Green)+"\n", i, opcode)
		case 0xC1:
			fmt.Printf("%.4X %.2X "+colorize("POP B", Green)+"\n", i, opcode)
		case 0xC2:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JNZ", Green)+colorize(" $%.4X", Cyan)+"\n", i-2, opcode, lb, hb, addr)
		case 0xC3:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JMP", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xC4:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("CNZ", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xC5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH B", Green)+"\n", i, opcode)
		case 0xC6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("ADI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xC8:
			fmt.Printf("%.4X %.2X "+colorize("RZ", Green)+"\n", i, opcode)
		case 0xC9:
			fmt.Printf("%.4X %.2X "+colorize("RET", Green)+"\n", i, opcode)
		case 0xCA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JZ", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xCC:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("CC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xCD:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("CALL", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)

		case 0xD1:
			fmt.Printf("%.4X %.2X "+colorize("POP D", Green)+"\n", i, opcode)
		case 0xD2:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JNC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xD5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH D", Green)+"\n", i, opcode)
		case 0xDA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xDB:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("IN ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)

		case 0xE1:
			fmt.Printf("%.4X %.2X "+colorize("POP H", Green)+"\n", i, opcode)
		case 0xE5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH H", Green)+"\n", i, opcode)
		case 0xE6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("ANI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xEB:
			fmt.Printf("%.4X %.2X "+colorize("XCHG", Green)+"\n", i, opcode)

		case 0xF1:
			fmt.Printf("%.4X %.2X "+colorize("POP PSW", Green)+"\n", i, opcode)
		case 0xF5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH PSW", Green)+"\n", i, opcode)
		case 0xFB:
			fmt.Printf("%.4X %.2X "+colorize("EI", Green)+"\n", i, opcode)
		case 0xFE:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("CPI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)

		default:
			fmt.Printf("Unknown Opcode %.2X\n", opcode)
			break L
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
