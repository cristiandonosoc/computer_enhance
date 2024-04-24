package homework_01

import (
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
			data, err := nasm.RunNasm(runfile)
			require.NoError(t, err)

			instructions, err := intel8086.ParseInstructions(data)
			require.NoError(t, err)
			require.NotEmpty(t, instructions)
		})
	}
}
