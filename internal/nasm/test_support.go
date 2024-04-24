package nasm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cristiandonosoc/computer_enhance/pkg/intel8086"

	"github.com/stretchr/testify/require"
)

func RunNasmTest(t *testing.T, asmFile string) {
	wantData, err := RunNasm(asmFile)
	require.NoError(t, err)

	instructions, err := intel8086.ParseInstructions(wantData)
	require.NoError(t, err)
	require.NotEmpty(t, instructions)

	tmp, err := os.MkdirTemp("", "homework_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	// Write the parsed instructions back to file.
	asm, err := intel8086.PrintInstructionsToAsmFormat(instructions)
	require.NoError(t, err)

	tmpFile := filepath.Join(tmp, filepath.Base(asmFile))
	err = os.WriteFile(tmpFile, []byte(asm), 0644)
	require.NoError(t, err)

	// Run asm on that file.
	gotData, err := RunNasm(tmpFile)
	require.NoError(t, err)

	require.Equal(t, wantData, gotData)
}
