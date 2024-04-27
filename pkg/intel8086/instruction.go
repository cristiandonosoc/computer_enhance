package intel8086

import (
	"fmt"
	"strings"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type InstructionKind byte

const (
	InstructionKind_MOV InstructionKind = 0b00100010
)

type InstructionData struct {
	D   bool
	W   bool
	MOD byte
	REG byte
	RM  byte
}

func (id *InstructionData) DebugString() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("; W: %t\n", id.W))
	sb.WriteString(fmt.Sprintf("; D: %t\n", id.D))
	sb.WriteString(fmt.Sprintf("; MOD: %02b\n", id.MOD))
	sb.WriteString(fmt.Sprintf("; REG: %03b\n", id.REG))
	sb.WriteString(fmt.Sprintf("; R/M: %03b\n", id.RM))

	return sb.String()
}

// Instruction is an interface that all opcode handlers must return.
type Instruction interface {
	Bytes() []byte
	Data() InstructionData
	String() string
	DebugString() string
}

// Generic Instruction -----------------------------------------------------------------------------

type genericPrintFunc func(gi *genericInstruction) string

// genericInstruction is an implementation of |Instruction| that conforms for a common usecase of
// just having some basic flag parsing and internal bytes.
type genericInstruction struct {
	data          InstructionData
	internalBytes []byte
	printFunc     genericPrintFunc

	// Mostly for debugging purposes.
	oh *OpcodeHandler
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

func (gi *genericInstruction) DebugString() string {
	var sb strings.Builder

	sb.WriteString(";\n")
	fmtStr := fmt.Sprintf("; OPCODE: %%0%db\n", gi.oh.OpcodeLength)
	sb.WriteString(fmt.Sprintf(fmtStr, gi.oh.Opcode))

	sb.WriteString("; BYTES: ")
	for _, b := range gi.internalBytes {
		sb.WriteString(fmt.Sprintf("%08b ", b))
	}
	sb.WriteString("\n")
	sb.WriteString(gi.data.DebugString())

	return sb.String()
}

// wrapPrintFunc is a helper wrapper that prints the debug data of each instruction before the
// actual printed instruction. This helps debugging of the generated code.
func wrapPrintFunc(pf genericPrintFunc) genericPrintFunc {
	return func(gi *genericInstruction) string {
		var sb strings.Builder

		sb.WriteString(gi.DebugString())
		sb.WriteString(pf(gi))

		return sb.String()
	}
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

	instruction, err := oh.Handler(oh, b1, b2, bs)
	if err != nil {
		return nil, fmt.Errorf("handling opcode %q (%08b): %w", oh.Name, oh.Opcode, err)
	}

	return instruction, nil
}
