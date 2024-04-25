package intel8086

import (
	"fmt"
	"sort"

	"github.com/cristiandonosoc/computer_enhance/internal/bytes"
)

type OpcodeDecoderFunc func(b1, b2 byte, bs *bytes.ByteStream) (Instruction, error)

type OpcodeHandler struct {
	Name         string
	Opcode       byte
	OpcodeLength int
	Handler      OpcodeDecoderFunc
}

var (
	// This are all the opcodes supported by our disassembler.
	// The length will be calculated at the package init function.
	kOpcodeHandlers = []*OpcodeHandler{
		{
			Name:    "Register/memory to/from register",
			Opcode:  0b100010,
			Handler: opcodeHandler_RegisterMemoryToFromRegister,
		},
		{
			Name:   "Immediate to register/memory",
			Opcode: 0b1100011,
		},
		{
			Name:   "Immediate to register",
			Opcode: 0b1011,
		},
	}
)

func init() {
	// Calculate the length of all the defined opcodes.
	for _, oh := range kOpcodeHandlers {
		oh.OpcodeLength = calculateOpcodeLength(oh.Opcode)
	}
}

// OpcodeRegistry -----------------------------------------------------------------------------------

type OpcodeRegistry struct {
	Opcodes []*OpcodeHandler
}

var kOpcodeRegistry *OpcodeRegistry

func GlobalOpcodeRegistry() *OpcodeRegistry {
	if kOpcodeRegistry != nil {
		return kOpcodeRegistry
	}

	od := &OpcodeRegistry{}
	od.Opcodes = make([]*OpcodeHandler, 0, len(kOpcodeHandlers))
	od.Opcodes = append(od.Opcodes, kOpcodeHandlers...)

	// Sort opcodes by value.
	sort.Slice(od.Opcodes, func(i, j int) bool {
		return od.Opcodes[i].Opcode < od.Opcodes[j].Opcode
	})

	// Set the global pointer.
	kOpcodeRegistry = od
	return kOpcodeRegistry
}

func (od *OpcodeRegistry) FindOpcodeHandler(b byte) (*OpcodeHandler, error) {
	for _, opcode := range od.Opcodes {
		shift := 8 - opcode.OpcodeLength
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
