package intel8086

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateOpcodeLength(t *testing.T) {
	testcases := []struct {
		opcode byte
		want   int
	}{
		{
			opcode: 0,
			want:   0,
		},
		{
			opcode: 0b1,
			want:   1,
		},
		{
			opcode: 0b1111,
			want:   4,
		},
		{
			opcode: 0b00100010,
			want:   6,
		},
		{
			opcode: 0b10100010,
			want:   8,
		},
	}

	for i, testcase := range testcases {
		got := calculateOpcodeLength(testcase.opcode)
		assert.Equal(t, testcase.want, got, "testcase %d", i)
	}
}

func TestFindOpcode(t *testing.T) {
	testcases := []struct {
		b        byte
		wantName string
		wantErr  string
	}{
		{
			b:       0b11111111,
			wantErr: "not found",
		},
		// Immediate to register.
		{
			b:        0b10110000,
			wantName: "Immediate to register",
		},
		{
			b:        0b10110001,
			wantName: "Immediate to register",
		},
		{
			b:        0b10110010,
			wantName: "Immediate to register",
		},
		{
			b:        0b10110101,
			wantName: "Immediate to register",
		},
		// Immediate to register/memory.
		{
			b:        0b11000110,
			wantName: "Immediate to register/memory",
		},
		{
			b:        0b11000111,
			wantName: "Immediate to register/memory",
		},
		{
			b:        0b10001000,
			wantName: "Register/memory to/from register",
		},
		{
			b:        0b10001001,
			wantName: "Register/memory to/from register",
		},
		{
			b:        0b10001010,
			wantName: "Register/memory to/from register",
		},
		{
			b:        0b10001011,
			wantName: "Register/memory to/from register",
		},
	}

	od := GlobalOpcodeDecoder()
	for _, opcode := range od.Opcodes {
		t.Logf("OPCODE: %08b, LENGTH: %d\n", opcode.Opcode, opcode.Length)
	}

	for i, testcase := range testcases {
		got, err := od.FindOpcode(testcase.b)

		if testcase.wantErr != "" {
			if assert.Error(t, err, "testcase %d", i) {
				assert.Contains(t, err.Error(), testcase.wantErr, "testcase %d", i)
			}
			continue
		}

		if assert.NoError(t, err, "testcase %d", i) {
			assert.Contains(t, got.Name, testcase.wantName, "testcase %d", i)
		}
	}
}
