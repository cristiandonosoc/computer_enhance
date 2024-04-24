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
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting pwd: %w", err)
	}

	rootDir, err := findRootDir(pwd)
	if err != nil {
		return nil, fmt.Errorf("searching for root dir: %w", err)
	}

	nasm, err := findNasm(rootDir)
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

const (
	kMarkerFilename = "ROOT_MARKER"
	// How much parent directories we're willing to search
	kMaxSearchIterations = 32
)

// findRootDir searches for a file that marks the ROOT directory of the project, so we can search
// nasm from there, regardless of wherever in the project we executed Go.
func findRootDir(dir string) (string, error) {
	originalDir := dir
	for i := 0; i < kMaxSearchIterations; i++ {
		path := filepath.Join(dir, kMarkerFilename)
		if _, exists, err := files.StatFile(path); err != nil {
			return "", fmt.Errorf("statting %q: %w", path, err)
		} else if exists {
			return dir, nil
		}

		// If it doesn't exist, we try to go to the parent directory.
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("marker file %q not found from %q", kMarkerFilename, originalDir)
		}

		dir = parentDir
	}

	return "", fmt.Errorf("max search iterations exceeded")
}

func findNasm(rootDir string) (string, error) {
	path := filepath.Join(rootDir, "extras", "nasm", "nasm.exe")
	if _, exists, err := files.StatFile(path); err != nil || !exists {
		return "", files.StatFileErrorf(err, "statting %q", path)
	}

	return path, nil
}
