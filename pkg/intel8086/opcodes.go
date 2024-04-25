package intel8086

import (
	"fmt"
	"sort"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type OpcodeHandler func(b1, b2 byte, bs *bytes.ByteStream) (*Instruction, error)

type Opcode struct {
	Name    string
	Opcode  byte
	Length  int
	Handler OpcodeHandler
}

var (
	kOpcodes = []*Opcode{
		{
			Name:   "Register/memory to/from register",
			Opcode: 0b100010,
		},
		{
			Name: "Immediate to register/memory",
			Opcode: 0b1100011,
		},
		{
			Name: "Immediate to register",
			Opcode: 0b1011,
		},
	}
)

func init() {
	// Calculate the length of all the defined opcodes.
	for _, opcode := range kOpcodes {
		opcode.Length = calculateOpcodeLength(opcode.Opcode)
	}
}

// OpcodeDecoder -----------------------------------------------------------------------------------

type OpcodeDecoder struct {
	Opcodes []*Opcode
}

var kOpcodeDecoder *OpcodeDecoder

func GlobalOpcodeDecoder() *OpcodeDecoder {
	if kOpcodeDecoder != nil {
		return kOpcodeDecoder
	}

	od := &OpcodeDecoder{}
	od.Opcodes = make([]*Opcode, 0, len(kOpcodes))
	od.Opcodes = append(od.Opcodes, kOpcodes...)

	// Sort opcodes by value.
	sort.Slice(od.Opcodes, func(i, j int) bool {
		return od.Opcodes[i].Opcode < od.Opcodes[j].Opcode
	})

	// Set the global pointer.
	kOpcodeDecoder = od
	return kOpcodeDecoder
}

func (od *OpcodeDecoder) FindOpcode(b byte) (*Opcode, error) {
	for _, opcode := range od.Opcodes {
		shift := 8 - opcode.Length
		candidateOpcode := b >> shift
		if candidateOpcode == opcode.Opcode {
			return opcode, nil
		}
	}

	return nil, fmt.Errorf("opcode in byte %08b not found", b)
}

func calculateOpcodeLength(opcode byte) int {
	for i := 7; i >= 0; i-- {
		bit := (opcode >> i) & 0b1
		if bit == 1 {
			return i + 1
		}
	}

	return 0
}
