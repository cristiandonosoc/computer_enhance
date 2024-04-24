package intel8086

import (
	"fmt"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type InstructionKind byte

const (
	InstructionKind_MOV InstructionKind = 0b00100010
)

type Instruction struct {
	Kind       InstructionKind
	B1         byte
	B2         byte
	ExtraBytes []byte
}

func (i *Instruction) Opcode() byte {
	return (i.B1 >> 2) & 0b00111111
}

func (i *Instruction) W() bool {
	return (i.B1 & 0b1) == 1
}

func (i *Instruction) D() bool {
	return (i.B1 & 0b01) == 1
}

func (i *Instruction) MOD() byte {
	return (i.B2 >> 6) & 0b11
}

func (i *Instruction) REG() byte {
	return (i.B2 >> 3) & 0b111
}

func (i *Instruction) RM() byte {
	return i.B2 & 0b111
}

func (i *Instruction) Print() (string, error) {
	switch kind := i.Kind; kind {
	case InstructionKind_MOV:
		return printMOV(i)
	default:
		return "", fmt.Errorf("unsupported kind: %v", kind)
	}
}

func ParseInstruction(bs *bytes.ByteStream) (*Instruction, error) {
	if bs.IsEOF() {
		return nil, fmt.Errorf("stream EOF")
	}

	i := &Instruction{
		B1: bs.Advance(),
		B2: bs.Advance(),
	}

	switch opcode := i.Opcode(); opcode {
	case byte(InstructionKind_MOV):
		i.Kind = InstructionKind_MOV
	default:
		return nil, fmt.Errorf("unknown opcode %08b", opcode)
	}

	return i, nil
}

// MOV ----------------------------------------------------------------------------------------------

func printMOV(i *Instruction) (string, error) {
	if i.MOD() != 0b11 {
		return "", fmt.Errorf("only register mode supported for now")
	}

	w := i.W()
	dst := InterpretRegister(i.RM(), w)
	src := InterpretRegister(i.REG(), w)

	return fmt.Sprintf("mov %s, %s", dst, src), nil
}
