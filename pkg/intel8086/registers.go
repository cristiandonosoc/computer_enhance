package intel8086

import ()

type Register string

const (
	Register_AL Register = "al"
	Register_CL          = "cl"
	Register_DL          = "dl"
	Register_BL          = "bl"
	Register_AH          = "ah"
	Register_CH          = "ch"
	Register_DH          = "dh"
	Register_BH          = "bh"

	Register_AX Register = "ax"
	Register_CX          = "cx"
	Register_DX          = "dx"
	Register_BX          = "bx"
	Register_SP          = "sp"
	Register_BP          = "bp"
	Register_SI          = "si"
	Register_DI          = "di"
)

var (
	// kRegistersByte is the table implementation of the W=0 part of the REG table.
	kRegistersByte = []Register{
		Register_AL, // 0b000
		Register_CL, // 0b001
		Register_DL, // 0b010
		Register_BL, // 0b011
		Register_AH, // 0b100
		Register_CH, // 0b101
		Register_DH, // 0b110
		Register_BH, // 0b111
	}

	// kRegistersByte is the table implementation of the W=1 part of the REG table.
	kRegistersWord = []Register{
		Register_AX, // 0b000
		Register_CX, // 0b001
		Register_DX, // 0b010
		Register_BX, // 0b011
		Register_SP, // 0b100
		Register_BP, // 0b101
		Register_SI, // 0b110
		Register_DI, // 0b111
	}
)

type RegisterWord byte

const ()

func InterpretRegister(reg byte, w bool) Register {
	index := reg & 0b111

	if !w {
		return kRegistersByte[index]
	} else {
		return kRegistersWord[index]
	}
}
