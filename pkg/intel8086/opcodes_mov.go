package intel8086

import (
	"fmt"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

// Register/memory to/from register ----------------------------------------------------------------

func opcodeHandler_RegisterMemoryToFromRegister(oh *OpcodeHandler, b1, b2 byte, bs *bytes.ByteStream) (Instruction, error) {
	data := InstructionData{
		D:   (b1 & 0b10) > 0,
		W:   (b1 & 0b1) > 0,
		MOD: (b2 >> 6) & 0b11,
		REG: (b2 >> 3) & 0b111,
		RM:  b2 & 0b111,
	}

	var internalBytes []byte
	switch data.MOD {
	case 0b00:
		// Check for direct addressing mode.
		if data.RM == 0b110 {
			b3, b4, err := bs.ReadWordAsBytes()
			if err != nil {
				return nil, fmt.Errorf("mod 00: %w", err)
			}
			internalBytes = []byte{b1, b2, b3, b4}
		} else {
			internalBytes = []byte{b1, b2}
		}
	case 0b01:
		b3, err := bs.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("mod 01: %w", err)
		}
		internalBytes = []byte{b1, b2, b3}
	case 0b10:
		b3, b4, err := bs.ReadWordAsBytes()
		if err != nil {
			return nil, fmt.Errorf("mod 10: %w", err)
		}
		internalBytes = []byte{b1, b2, b3, b4}
	case 0b11:
		internalBytes = []byte{b1, b2}
	}

	printFunc := func(gi *genericInstruction) string {
		// Check for the Effective Address Calculation.
		var eac EAC = EAC_Invalid
		if gi.data.MOD != 0b11 {
			eac = InterpretRM(gi.data.RM, gi.data.MOD)
		}

		switch gi.data.MOD {
		case 0b00:
			if eac != EAC_DirectAddress {
				dst := InterpretREG(gi.data.REG, gi.data.W)
				return fmt.Sprintf("mov %s, %s", dst, ToEACNotation(eac, 0))
			} else {
				src := InterpretREG(gi.data.REG, gi.data.W)
				return fmt.Sprintf("mov %s, %s", ToEACNotation(eac, 0), src)
			}
		case 0b01:
			dst := InterpretREG(gi.data.REG, gi.data.W)
			return fmt.Sprintf("mov %s, %s", dst, ToEACNotation(eac, uint16(gi.internalBytes[2])))
		case 0b10:
			dst := InterpretREG(gi.data.REG, gi.data.W)
			offset := bytes.ToUint16(gi.internalBytes[2], gi.internalBytes[3])
			return fmt.Sprintf("mov %s, %s", dst, ToEACNotation(eac, offset))
		case 0b11:
			dst := InterpretREG(gi.data.RM, gi.data.W)
			src := InterpretREG(gi.data.REG, gi.data.W)
			return fmt.Sprintf("mov %s, %s", dst, src)
		}
		panic(fmt.Sprintf("invalid mod value %08b", gi.data.MOD))
	}

	return &genericInstruction{
		data:          data,
		internalBytes: internalBytes,
		printFunc:     wrapPrintFunc(printFunc),
		oh:            oh,
	}, nil
}

// Immediate to Register ---------------------------------------------------------------------------

func opcodeHandler_ImmediateToRegister(oh *OpcodeHandler, b1, b2 byte, bs *bytes.ByteStream) (Instruction, error) {
	data := InstructionData{
		W:   ((b1 >> 3) & 0b1) > 0,
		REG: (b1 & 0b111),
	}

	// See if this is a wide instruction, and if so, load the extra byte.
	var internalBytes []byte
	if !data.W {
		internalBytes = []byte{b1, b2}
	} else {
		if bs.IsEOF() {
			return nil, fmt.Errorf("stream EOF when loading byte 3")
		}

		b3 := bs.Advance()
		internalBytes = []byte{b1, b2, b3}
	}

	printFunc := func(gi *genericInstruction) string {
		dst := InterpretREG(gi.data.REG, gi.data.W)

		if !gi.data.W {
			return fmt.Sprintf("mov %s, %d", dst, gi.internalBytes[1])
		}

		// Return the 16 bit value, since this is wide.
		value := bytes.ToUint16(gi.internalBytes[1], gi.internalBytes[2])
		return fmt.Sprintf("mov %s, %d", dst, value)
	}

	return &genericInstruction{
		data:          data,
		internalBytes: internalBytes,
		printFunc:     wrapPrintFunc(printFunc),
		oh:            oh,
	}, nil
}
