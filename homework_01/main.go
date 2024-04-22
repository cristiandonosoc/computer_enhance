package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cristiandonosoc/computer_enhance/internal/nasm"
)


func internalMain() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: main <BINARY_FILE>")
	}

	path := os.Args[1]

	data, err := nasm.RunNasm(path)
	if err != nil {
		return fmt.Errorf("running nasm on %q: %w", path, err)
	}

	for i, b := range data {
		fmt.Printf("BYTE %03d: %08b\n", i, b)
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
