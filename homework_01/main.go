package main

import (
	"fmt"
	"os"
	"encoding/hex"
)

func internalMain() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: main <BINARY_FILE>")
	}

	path := os.Args[1]
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %q: %w", path, err)
	}

	hs := hex.EncodeToString(data)
	fmt.Println(hs)

	return fmt.Errorf("TODO: Conversion!")
}

func main() {
	if err := internalMain(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
