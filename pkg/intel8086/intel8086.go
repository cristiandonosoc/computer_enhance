// Package intel8086 represents the parser for the Intel 8086 CPU.
package intel8086

import (
	"fmt"
	"strings"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

func ParseInstructions(asm []byte) ([]*Instruction, error) {
	bs := bytes.NewByteStream(asm)

	var instructions []*Instruction
	for !bs.IsEOF() {
		instruction, err := ParseInstruction(bs)
		if err != nil {
			return nil, fmt.Errorf("parsing instruction: %w", err)
		}
		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

// printInstructions prints the .asm content.
func PrintInstructionsToAsmFormat(instructions []*Instruction) (string, error) {
	var sb strings.Builder

	sb.WriteString("bits 16\n\n")

	for i, instruction := range instructions {
		str, err := instruction.Print()
		if err != nil {
			return "", fmt.Errorf("printing instruction %d: %w", i, err)
		}

		sb.WriteString(strings.TrimSpace(str) + "\n")
	}

	return sb.String(), nil
}
