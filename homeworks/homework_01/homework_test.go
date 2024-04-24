package homework_01

import (
	"path/filepath"
	"testing"

	"github.com/cristiandonosoc/computer_enhance/internal/nasm"
	"github.com/cristiandonosoc/golib/pkg/test_support"

	"github.com/stretchr/testify/require"
)

func TestHomework(t *testing.T) {
	runfiles, err := test_support.Runfiles("testdata")
	require.NoError(t, err)

	for _, runfile := range runfiles {
		t.Run(filepath.Base(runfile), func(t *testing.T) {
			nasm.RunNasmTest(t, runfile)
		})
	}
}
