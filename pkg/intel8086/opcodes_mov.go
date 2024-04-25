package intel8086

import (
	"fmt"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type genericInstruction struct {
	data          InstructionData
	internalBytes []byte
	printFunc     func(gi *genericInstruction) string
}

var _ Instruction = (*genericInstruction)(nil)

func (gi *genericInstruction) Bytes() []byte {
	return gi.internalBytes
}

func (gi *genericInstruction) Data() InstructionData {
	return gi.data
}

func (gi *genericInstruction) String() string {
	return gi.printFunc(gi)
}

// Register/memory to/from register ----------------------------------------------------------------

func opcodeHandler_RegisterMemoryToFromRegister(b1, b2 byte, bs *bytes.ByteStream) (Instruction, error) {
	mod := (b2 >> 6) & 0b11

	if mod != 0b11 {
		return nil, fmt.Errorf("only register mode is supported for now")
	}

	printFunc := func(gi *genericInstruction) string {
		if gi.data.MOD != 0b11 {
			panic("only register mode supported for now")
		}

		dst := InterpretRegister(gi.data.RM, gi.data.W)
		src := InterpretRegister(gi.data.REG, gi.data.W)

		return fmt.Sprintf("mov %s, %s", dst, src)
	}

	return &genericInstruction{
		data: InstructionData{
			W:   (b1 & 0b1) > 0,
			D:   (b1 & 0b10) > 0,
			MOD: mod,
			REG: (b2 >> 3) & 0b111,
			RM:  b2 & 0b111,
		},
		internalBytes: []byte{b1, b2},
		printFunc:     printFunc,
	}, nil
}

// Immediate to Register ---------------------------------------------------------------------------

func opcodeHandler_ImmediateToRegister(b1, b2 byte, bs *bytes.ByteStream) (Instruction, error) {
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
		dst := InterpretRegister(gi.data.REG, gi.data.W)

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
		printFunc:     printFunc,
	}, nil
}
