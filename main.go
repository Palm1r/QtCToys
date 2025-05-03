package main

import (
	"fmt"
	"os"

	"github.com/Palm1r/QtCToys/qtcreator"
)

func main() {
	info, err := qtcreator.GetInfo()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Qt Creator Version: %s\n", info.Version)
	fmt.Printf("Qt Creator Path: %s\n", info.Path)
}
