package main

import (
	"bufio"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/legolord208/stdutil"
)

var stop bool

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		printActions()
		return
	}

	go func() {
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		stop = true
	}()

	debug := false
	switch strings.ToLower(args[0]) {
	case "debug":
		debug = true
		fallthrough
	case "run":
		scanner := bufio.NewScanner(os.Stdin)
		code, ok := readUntilEOF(scanner)
		if !ok {
			return
		}

		err := run(&env{
			code:  code,
			debug: debug,
		})

		os.Stdout.Write([]byte("\n"))
		if err != nil {
			stdutil.PrintErr("Error while executing code", err)
		}
	case "simplify":
		scanner := bufio.NewScanner(os.Stdin)
		code, ok := readUntilEOF(scanner)
		if !ok {
			return
		}

		os.Stdout.Write([]byte(simplify(code)))
	default:
		printActions()
	}
}

func printActions() {
	stdutil.PrintErr("Actions: run, debug, simplify", nil)
}

func readUntilEOF(scanner *bufio.Scanner) (str string, ok bool) {
	ok = true
	for scanner.Scan() {
		str += scanner.Text() // + "\n"
	}

	if err := scanner.Err(); err != nil {
		stdutil.PrintErr("Could not read line", err)
		ok = false
	}
	return
}
