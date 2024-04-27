package intel8086

import "fmt"

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

func InterpretREG(reg byte, w bool) Register {
	index := reg & 0b111

	if !w {
		return kRegistersByte[index]
	} else {
		return kRegistersWord[index]
	}
}

// Effective Address Calculator --------------------------------------------------------------------

type EAC string

const (
	EAC_Invalid       = "<invalid>"
	EAC_BXSI          = "bx + si"
	EAC_BXDI          = "bx + di"
	EAC_BPSI          = "bp + si"
	EAC_BPDI          = "bp + di"
	EAC_SI            = "si"
	EAC_DI            = "di"
	EAC_BP            = "bp"
	EAC_BX            = "bx"
	EAC_DirectAddress = "<direct_access>"
)

var (
	kEffectiveAddressCalculations = []EAC{
		EAC_BXSI,          // 0b000
		EAC_BXDI,          // 0b001
		EAC_BPSI,          // 0b010
		EAC_BPDI,          // 0b011
		EAC_SI,            // 0b100
		EAC_DI,            // 0b101
		EAC_BP,            // 0b110
		EAC_BX,            // 0b111
		EAC_DirectAddress, // 0b110
	}
)

func InterpretRM(rm byte, mod byte) EAC {
	// Check for direct address mode.
	if mod == 0b00 && rm == 0b110 {
		return EAC_DirectAddress
	}

	return kEffectiveAddressCalculations[rm]
}

func ToEACNotation(eac EAC, offset uint16) string {
	if offset == 0 {
		return fmt.Sprintf("[%s]", eac)
	}

	return fmt.Sprintf("[%s + %d]", eac, offset)
}
