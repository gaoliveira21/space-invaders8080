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
		case 0x02:
			fmt.Printf("%.4X %.2X "+colorize("STAX B", Green)+"\n", i, opcode)
		case 0x03:
			fmt.Printf("%.4X %.2X "+colorize("INX B", Green)+"\n", i, opcode)
		case 0x04:
			fmt.Printf("%.4X %.2X "+colorize("INR B", Green)+"\n", i, opcode)
		case 0x05:
			fmt.Printf("%.4X %.2X "+colorize("DCR B", Green)+"\n", i, opcode)
		case 0x06:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI B,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x07:
			fmt.Printf("%.4X %.2X "+colorize("RLC", Green)+"\n", i, opcode)
		case 0x08:
			fmt.Printf("%.4X %.2X "+colorize("NOP", Green)+"\n", i, opcode)
		case 0x09:
			fmt.Printf("%.4X %.2X "+colorize("DAD B", Green)+"\n", i, opcode)
		case 0x0A:
			fmt.Printf("%.4X %.2X "+colorize("LDAX B", Green)+"\n", i, opcode)
		case 0x0B:
			fmt.Printf("%.4X %.2X "+colorize("DCX B", Green)+"\n", i, opcode)
		case 0x0C:
			fmt.Printf("%.4X %.2X "+colorize("INR C", Green)+"\n", i, opcode)
		case 0x0D:
			fmt.Printf("%.4X %.2X "+colorize("DCR C", Green)+"\n", i, opcode)
		case 0x0E:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI C,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x0F:
			fmt.Printf("%.4X %.2X "+colorize("RRC", Green)+"\n", i, opcode)

		case 0x11:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI D,", Green)+colorize(" #$%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x12:
			fmt.Printf("%.4X %.2X "+colorize("STAX D", Green)+"\n", i, opcode)
		case 0x13:
			fmt.Printf("%.4X %.2X "+colorize("INX D", Green)+"\n", i, opcode)
		case 0x14:
			fmt.Printf("%.4X %.2X "+colorize("INR D", Green)+"\n", i, opcode)
		case 0x15:
			fmt.Printf("%.4X %.2X "+colorize("DCR D", Green)+"\n", i, opcode)
		case 0x16:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI D,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x19:
			fmt.Printf("%.4X %.2X "+colorize("DAD D", Green)+"\n", i, opcode)
		case 0x1A:
			fmt.Printf("%.4X %.2X "+colorize("LDAX D", Green)+"\n", i, opcode)

		case 0x21:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI H,", Green)+colorize(" #$%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x22:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("SHLD", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x23:
			fmt.Printf("%.4X %.2X "+colorize("INX H", Green)+"\n", i, opcode)
		case 0x24:
			fmt.Printf("%.4X %.2X "+colorize("INR H", Green)+"\n", i, opcode)
		case 0x25:
			fmt.Printf("%.4X %.2X "+colorize("DCR H", Green)+"\n", i, opcode)
		case 0x26:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI H,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x27:
			fmt.Printf("%.4X %.2X "+colorize("DAA", Green)+"\n", i, opcode)
		case 0x28:
			fmt.Printf("%.4X %.2X "+colorize("NOP", Green)+"\n", i, opcode)
		case 0x29:
			fmt.Printf("%.4X %.2X "+colorize("DAD H", Green)+"\n", i, opcode)
		case 0x2A:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LHLD", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x2B:
			fmt.Printf("%.4X %.2X "+colorize("DCX H", Green)+"\n", i, opcode)
		case 0x2C:
			fmt.Printf("%.4X %.2X "+colorize("INR L", Green)+"\n", i, opcode)
		case 0x2D:
			fmt.Printf("%.4X %.2X "+colorize("DCR L", Green)+"\n", i, opcode)
		case 0x2E:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI L,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x2F:
			fmt.Printf("%.4X %.2X "+colorize("CMA", Green)+"\n", i, opcode)

		case 0x31:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("LXI SP,", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x32:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("STA", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0x33:
			fmt.Printf("%.4X %.2X "+colorize("INX SP", Green)+"\n", i, opcode)
		case 0x34:
			fmt.Printf("%.4X %.2X "+colorize("INR M", Green)+"\n", i, opcode)
		case 0x35:
			fmt.Printf("%.4X %.2X "+colorize("DCR M", Green)+"\n", i, opcode)
		case 0x36:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("MVI M,", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0x37:
			fmt.Printf("%.4X %.2X "+colorize("STC", Green)+"\n", i, opcode)
		case 0x38:
			fmt.Printf("%.4X %.2X "+colorize("NOP", Green)+"\n", i, opcode)
		case 0x39:
			fmt.Printf("%.4X %.2X "+colorize("DAD SP", Green)+"\n", i, opcode)
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

		case 0x40:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,B", Green)+"\n", i, opcode)
		case 0x41:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,C", Green)+"\n", i, opcode)
		case 0x42:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,D", Green)+"\n", i, opcode)
		case 0x43:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,E", Green)+"\n", i, opcode)
		case 0x44:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,H", Green)+"\n", i, opcode)
		case 0x45:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,L", Green)+"\n", i, opcode)
		case 0x46:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,M", Green)+"\n", i, opcode)
		case 0x47:
			fmt.Printf("%.4X %.2X "+colorize("MOV B,A", Green)+"\n", i, opcode)
		case 0x48:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,B", Green)+"\n", i, opcode)
		case 0x49:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,C", Green)+"\n", i, opcode)
		case 0x4A:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,D", Green)+"\n", i, opcode)
		case 0x4B:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,E", Green)+"\n", i, opcode)
		case 0x4C:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,H", Green)+"\n", i, opcode)
		case 0x4E:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,M", Green)+"\n", i, opcode)
		case 0x4F:
			fmt.Printf("%.4X %.2X "+colorize("MOV C,A", Green)+"\n", i, opcode)

		case 0x56:
			fmt.Printf("%.4X %.2X "+colorize("MOV D,M", Green)+"\n", i, opcode)
		case 0x57:
			fmt.Printf("%.4X %.2X "+colorize("MOV D,A", Green)+"\n", i, opcode)
		case 0x5E:
			fmt.Printf("%.4X %.2X "+colorize("MOV E,M", Green)+"\n", i, opcode)
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
		case 0x68:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,B", Green)+"\n", i, opcode)
		case 0x69:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,C", Green)+"\n", i, opcode)
		case 0x6A:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,D", Green)+"\n", i, opcode)
		case 0x6B:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,E", Green)+"\n", i, opcode)
		case 0x6C:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,H", Green)+"\n", i, opcode)
		case 0x6D:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,L", Green)+"\n", i, opcode)
		case 0x6E:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,M", Green)+"\n", i, opcode)
		case 0x6F:
			fmt.Printf("%.4X %.2X "+colorize("MOV L,A", Green)+"\n", i, opcode)

		case 0x70:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,B", Green)+"\n", i, opcode)
		case 0x71:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,C", Green)+"\n", i, opcode)
		case 0x72:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,D", Green)+"\n", i, opcode)
		case 0x73:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,E", Green)+"\n", i, opcode)
		case 0x74:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,H", Green)+"\n", i, opcode)
		case 0x75:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,L", Green)+"\n", i, opcode)
		case 0x76:
			fmt.Printf("%.4X %.2X "+colorize("HLT", Green)+"\n", i, opcode)
		case 0x77:
			fmt.Printf("%.4X %.2X "+colorize("MOV M,A", Green)+"\n", i, opcode)
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

		case 0x80:
			fmt.Printf("%.4X %.2X "+colorize("ADD B", Green)+"\n", i, opcode)
		case 0x81:
			fmt.Printf("%.4X %.2X "+colorize("ADD C", Green)+"\n", i, opcode)
		case 0x82:
			fmt.Printf("%.4X %.2X "+colorize("ADD D", Green)+"\n", i, opcode)
		case 0x83:
			fmt.Printf("%.4X %.2X "+colorize("ADD E", Green)+"\n", i, opcode)
		case 0x84:
			fmt.Printf("%.4X %.2X "+colorize("ADD H", Green)+"\n", i, opcode)
		case 0x85:
			fmt.Printf("%.4X %.2X "+colorize("ADD L", Green)+"\n", i, opcode)
		case 0x86:
			fmt.Printf("%.4X %.2X "+colorize("ADD M", Green)+"\n", i, opcode)
		case 0x87:
			fmt.Printf("%.4X %.2X "+colorize("ADD A", Green)+"\n", i, opcode)
		case 0x88:
			fmt.Printf("%.4X %.2X "+colorize("ADC B", Green)+"\n", i, opcode)
		case 0x89:
			fmt.Printf("%.4X %.2X "+colorize("ADC C", Green)+"\n", i, opcode)
		case 0x8A:
			fmt.Printf("%.4X %.2X "+colorize("ADC D", Green)+"\n", i, opcode)
		case 0x8B:
			fmt.Printf("%.4X %.2X "+colorize("ADC E", Green)+"\n", i, opcode)
		case 0x8C:
			fmt.Printf("%.4X %.2X "+colorize("ADC H", Green)+"\n", i, opcode)
		case 0x8D:
			fmt.Printf("%.4X %.2X "+colorize("ADC L", Green)+"\n", i, opcode)
		case 0x8E:
			fmt.Printf("%.4X %.2X "+colorize("ADC M", Green)+"\n", i, opcode)

		case 0x90:
			fmt.Printf("%.4X %.2X "+colorize("SUB B", Green)+"\n", i, opcode)
		case 0x91:
			fmt.Printf("%.4X %.2X "+colorize("SUB C", Green)+"\n", i, opcode)
		case 0x92:
			fmt.Printf("%.4X %.2X "+colorize("SUB D", Green)+"\n", i, opcode)
		case 0x93:
			fmt.Printf("%.4X %.2X "+colorize("SUB E", Green)+"\n", i, opcode)
		case 0x94:
			fmt.Printf("%.4X %.2X "+colorize("SUB H", Green)+"\n", i, opcode)
		case 0x95:
			fmt.Printf("%.4X %.2X "+colorize("SUB L", Green)+"\n", i, opcode)
		case 0x96:
			fmt.Printf("%.4X %.2X "+colorize("SUB M", Green)+"\n", i, opcode)
		case 0x97:
			fmt.Printf("%.4X %.2X "+colorize("SUB A", Green)+"\n", i, opcode)

		case 0xA0:
			fmt.Printf("%.4X %.2X "+colorize("ANA B", Green)+"\n", i, opcode)
		case 0xA1:
			fmt.Printf("%.4X %.2X "+colorize("ANA C", Green)+"\n", i, opcode)
		case 0xA2:
			fmt.Printf("%.4X %.2X "+colorize("ANA D", Green)+"\n", i, opcode)
		case 0xA3:
			fmt.Printf("%.4X %.2X "+colorize("ANA E", Green)+"\n", i, opcode)
		case 0xA4:
			fmt.Printf("%.4X %.2X "+colorize("ANA H", Green)+"\n", i, opcode)
		case 0xA5:
			fmt.Printf("%.4X %.2X "+colorize("ANA L", Green)+"\n", i, opcode)
		case 0xA6:
			fmt.Printf("%.4X %.2X "+colorize("ANA M", Green)+"\n", i, opcode)
		case 0xA7:
			fmt.Printf("%.4X %.2X "+colorize("ANA A", Green)+"\n", i, opcode)
		case 0xA8:
			fmt.Printf("%.4X %.2X "+colorize("XRA B", Green)+"\n", i, opcode)
		case 0xA9:
			fmt.Printf("%.4X %.2X "+colorize("XRA C", Green)+"\n", i, opcode)
		case 0xAA:
			fmt.Printf("%.4X %.2X "+colorize("XRA D", Green)+"\n", i, opcode)
		case 0xAB:
			fmt.Printf("%.4X %.2X "+colorize("XRA E", Green)+"\n", i, opcode)
		case 0xAC:
			fmt.Printf("%.4X %.2X "+colorize("XRA H", Green)+"\n", i, opcode)
		case 0xAD:
			fmt.Printf("%.4X %.2X "+colorize("XRA L", Green)+"\n", i, opcode)
		case 0xAE:
			fmt.Printf("%.4X %.2X "+colorize("XRA M", Green)+"\n", i, opcode)
		case 0xAF:
			fmt.Printf("%.4X %.2X "+colorize("XRA A", Green)+"\n", i, opcode)

		case 0xB0:
			fmt.Printf("%.4X %.2X "+colorize("ORA B", Green)+"\n", i, opcode)
		case 0xB1:
			fmt.Printf("%.4X %.2X "+colorize("ORA C", Green)+"\n", i, opcode)
		case 0xB2:
			fmt.Printf("%.4X %.2X "+colorize("ORA D", Green)+"\n", i, opcode)
		case 0xB3:
			fmt.Printf("%.4X %.2X "+colorize("ORA E", Green)+"\n", i, opcode)
		case 0xB4:
			fmt.Printf("%.4X %.2X "+colorize("ORA H", Green)+"\n", i, opcode)
		case 0xB5:
			fmt.Printf("%.4X %.2X "+colorize("ORA L", Green)+"\n", i, opcode)
		case 0xB6:
			fmt.Printf("%.4X %.2X "+colorize("ORA M", Green)+"\n", i, opcode)
		case 0xB7:
			fmt.Printf("%.4X %.2X "+colorize("ORA A", Green)+"\n", i, opcode)
		case 0xB8:
			fmt.Printf("%.4X %.2X "+colorize("CMP B", Green)+"\n", i, opcode)
		case 0xB9:
			fmt.Printf("%.4X %.2X "+colorize("CMP C", Green)+"\n", i, opcode)
		case 0xBA:
			fmt.Printf("%.4X %.2X "+colorize("CMP D", Green)+"\n", i, opcode)
		case 0xBB:
			fmt.Printf("%.4X %.2X "+colorize("CMP E", Green)+"\n", i, opcode)
		case 0xBC:
			fmt.Printf("%.4X %.2X "+colorize("CMP H", Green)+"\n", i, opcode)
		case 0xBD:
			fmt.Printf("%.4X %.2X "+colorize("CMP L", Green)+"\n", i, opcode)
		case 0xBE:
			fmt.Printf("%.4X %.2X "+colorize("CMP M", Green)+"\n", i, opcode)
		case 0xBF:
			fmt.Printf("%.4X %.2X "+colorize("CMP A", Green)+"\n", i, opcode)

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

		case 0xD0:
			fmt.Printf("%.4X %.2X "+colorize("RNC", Green)+"\n", i, opcode)
		case 0xD1:
			fmt.Printf("%.4X %.2X "+colorize("POP D", Green)+"\n", i, opcode)
		case 0xD2:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JNC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xD3:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("OUT ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xD4:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("CNC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xD5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH D", Green)+"\n", i, opcode)
		case 0xD6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("SUI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xD7:
			fmt.Printf("%.4X %.2X "+colorize("RST 2", Green)+"\n", i, opcode)
		case 0xD8:
			fmt.Printf("%.4X %.2X "+colorize("RC", Green)+"\n", i, opcode)
		case 0xD9:
			fmt.Printf("%.4X %.2X "+colorize("RET", Green)+"\n", i, opcode)
		case 0xDA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JC", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
		case 0xDB:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("IN ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xDE:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("SBI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)

		case 0xE1:
			fmt.Printf("%.4X %.2X "+colorize("POP H", Green)+"\n", i, opcode)
		case 0xE3:
			fmt.Printf("%.4X %.2X "+colorize("XTHL", Green)+"\n", i, opcode)
		case 0xE5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH H", Green)+"\n", i, opcode)
		case 0xE6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("ANI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xE9:
			fmt.Printf("%.4X %.2X "+colorize("PCHL", Green)+"\n", i, opcode)
		case 0xEB:
			fmt.Printf("%.4X %.2X "+colorize("XCHG", Green)+"\n", i, opcode)

		case 0xF1:
			fmt.Printf("%.4X %.2X "+colorize("POP PSW", Green)+"\n", i, opcode)
		case 0xF5:
			fmt.Printf("%.4X %.2X "+colorize("PUSH PSW", Green)+"\n", i, opcode)
		case 0xF6:
			i++
			b := uint16(rom[i])
			fmt.Printf("%.4X %.2X %.2X "+colorize("ORI ", Green)+colorize(" #$0x%.2X\n", Cyan), i-1, opcode, b, b)
		case 0xFA:
			lb, hb, addr := getAddr(&i, rom)
			fmt.Printf("%.4X %.2X %.2X %.2X "+colorize("JM", Green)+colorize(" $%.4X\n", Cyan), i-2, opcode, lb, hb, addr)
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
