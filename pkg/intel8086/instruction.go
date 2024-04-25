package intel8086

import (
	"fmt"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type InstructionKind byte

const (
	InstructionKind_MOV InstructionKind = 0b00100010
)

type InstructionData struct {
	W   bool
	D   bool
	MOD byte
	REG byte
	RM  byte
}

type Instruction interface {
	Bytes() []byte
	Data() InstructionData
	String() string
}

func ParseInstruction(bs *bytes.ByteStream) (Instruction, error) {
	if bs.IsEOF() {
		return nil, fmt.Errorf("stream EOF at byte 1")
	}
	b1 := bs.Advance()

	if bs.IsEOF() {
		return nil, fmt.Errorf("stream EOF at byte 2")
	}
	b2 := bs.Advance()

	or := GlobalOpcodeRegistry()
	oh, err := or.FindOpcodeHandler(b1)
	if err != nil {
		return nil, fmt.Errorf("finding opcode handler: %w", err)
	}

	if oh.Handler == nil {
		return nil, fmt.Errorf("opcode handler for %q (%08b) not defined!", oh.Name, oh.Opcode)
	}

	instruction, err := oh.Handler(b1, b2, bs)
	if err != nil {
		return nil, fmt.Errorf("handling opcode %q (%08b): %w", oh.Name, oh.Opcode, err)
	}

	return instruction, nil
}
