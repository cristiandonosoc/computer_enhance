// Package nasm is a helper package for running nasm over input files that are in assembly.
package nasm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cristiandonosoc/golib/pkg/files"
)

func RunNasm(input string) ([]byte, error) {
	nasm, err := findNasm()
	if err != nil {
		return nil, fmt.Errorf("finding nasm: %w", err)
	}

	tmp, err := os.MkdirTemp("", "nasm")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	output := filepath.Join(tmp, "output")
	cmd := exec.Command(nasm, "-o", output, input)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("running %v: %w", cmd.Args, err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		return nil, fmt.Errorf("reading output at %q: %w", output, err)
	}

	return data, nil
}

func findNasm() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting pwd: %w", err)
	}

	path := filepath.Join(pwd, "extras", "nasm", "nasm.exe")
	if _, exists, err := files.StatFile(path); err != nil || !exists {
		return "", files.StatFileErrorf(err, "statting %q", path)
	}

	return path, nil
}
