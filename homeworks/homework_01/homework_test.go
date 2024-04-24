package homework_01

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cristiandonosoc/computer_enhance/internal/nasm"
	"github.com/cristiandonosoc/computer_enhance/pkg/intel8086"
	"github.com/cristiandonosoc/golib/pkg/test_support"

	"github.com/stretchr/testify/require"
)

func TestHomework(t *testing.T) {
	runfiles, err := test_support.Runfiles("testdata")
	require.NoError(t, err)

	for _, runfile := range runfiles {
		t.Run(filepath.Base(runfile), func(t *testing.T) {
			wantData, err := nasm.RunNasm(runfile)
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

			tmpFile := filepath.Join(tmp, filepath.Base(runfile))
			err = os.WriteFile(tmpFile, []byte(asm), 0644)
			require.NoError(t, err)

			// Run asm on that file.
			gotData, err := nasm.RunNasm(tmpFile)
			require.NoError(t, err)

			require.Equal(t, wantData, gotData)
		})
	}
}
