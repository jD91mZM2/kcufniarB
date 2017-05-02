//+build !windows

package main

import (
	"os"
	"os/exec"
)

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	// What should I do with the potential error?...
}
