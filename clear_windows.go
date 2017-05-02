package main

import (
	"os"
	"os/exec"
)

func clear() {
	cmd := exec.Command("cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// What should I do with the potential error?...
}
