package intel8086

import (
	"fmt"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

// Register/memory to/from register ----------------------------------------------------------------

var _ Instruction = (*Instruction_RegisterMemoryToFromRegister)(nil)

type Instruction_RegisterMemoryToFromRegister struct {
	data          InstructionData
	internalBytes []byte
}

func (i *Instruction_RegisterMemoryToFromRegister) Bytes() []byte {
	return i.internalBytes
}

func (i *Instruction_RegisterMemoryToFromRegister) Data() InstructionData {
	return i.data
}

func (i *Instruction_RegisterMemoryToFromRegister) String() string {

	if i.data.MOD != 0b11 {
		panic("only register mode supported for now")
	}

	dst := InterpretRegister(i.data.RM, i.data.W)
	src := InterpretRegister(i.data.REG, i.data.W)

	return fmt.Sprintf("mov %s, %s", dst, src)
}

func opcodeHandler_RegisterMemoryToFromRegister(b1, b2 byte, bs *bytes.ByteStream) (Instruction, error) {
	mod := (b2 >> 6) & 0b11

	if mod != 0b11 {
		return nil, fmt.Errorf("only register mode is supported for now")
	}

	i := &Instruction_RegisterMemoryToFromRegister{
		data: InstructionData{
			W:   (b1 & 0b1) > 0,
			D:   (b1 & 0b10) > 0,
			MOD: mod,
			REG: (b2 >> 3) & 0b111,
			RM:  b2 & 0b111,
		},
		internalBytes: []byte{b1, b2},
	}

	return i, nil
}
