package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cristiandonosoc/computer_enhance/internal/nasm"
	"github.com/cristiandonosoc/computer_enhance/pkg/intel8086"
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

	hs := hex.EncodeToString(data)
	fmt.Println(hs)

	fmt.Println("--------------------------")

	instructions, err := intel8086.ParseInstructions(data)
	if err != nil {
		return fmt.Errorf("parsing instructions: %w", err)
	}

	for _, instruction := range instructions {
		str, err := instruction.Print()
		if err != nil {
			return fmt.Errorf("printing instruction: %w", err)
		}

		fmt.Println("INSTRUCTION:", str)
	}

	return nil
}

func main() {
	if err := internalMain(); err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
