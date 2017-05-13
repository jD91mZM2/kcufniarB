package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
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
		code, ok := readFileOrStdin(args[1:])
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
		simplifyFileOrStdin(new(normsimplifier), args)
	case "genc":
		simplifyFileOrStdin(new(csimplifier), args)
	case "genasm64":
		simplifyFileOrStdin(new(asm64simplifier), args)
	case "genval":
		if len(args) < 2 {
			stdutil.PrintErr("Usage: genval <number>\nGenerate code for number", nil)
			return
		}

		n, err := strconv.Atoi(args[1])
		if err != nil {
			stdutil.PrintErr("Not a number", nil)
			return
		}
		if n <= 0 {
			stdutil.PrintErr("Number must be more than 0", nil)
			return
		}

		x, y, diff := findmultiple(n)

		if x == 0 && y == 0 && diff == 0 {
			os.Stdout.Write([]byte("Multiple not found!\n"))
			return
		}

		os.Stdout.Write([]byte(genmultiple(x, y, diff) + "\n"))
	case "genstr":
		if len(args) < 2 {
			stdutil.PrintErr("Usage: genstr <string>\nGenerate code for string", nil)
			return
		}

		s := ""

		last := len(args[1]) - 1
		for i, c := range args[1] {
			x, y, diff := findmultiple(int(c))

			if x == 0 && y == 0 && diff == 0 {
				os.Stdout.Write([]byte("Multiple not found for '" + string(c) + "'"))
				return
			}
			s += genmultiple(x, y, diff) + "."

			if i != last {
				s += ">"
			}
		}
		os.Stdout.Write([]byte(s + "\n"))
	default:
		printActions()
	}
}

func printActions() {
	stdutil.PrintErr("Actions: run, debug, genval, genstr, simplify, genc", nil)
}

func simplifyFileOrStdin(s simplifier, args []string) {
	code, ok := readFileOrStdin(args[1:])
	if !ok {
		return
	}
	os.Stdout.Write([]byte(simplify(code, s)))
}

func readFileOrStdin(args []string) (str string, ok bool) {
	if len(args) < 1 {
		scanner := bufio.NewScanner(os.Stdin)
		str, ok = readUntilEOF(scanner)
		return
	}

	bytes, err := ioutil.ReadFile(args[0])
	if err != nil {
		stdutil.PrintErr("Could not read file", err)
		return
	}
	str = string(bytes)
	ok = true
	return
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

func genmultiple(x int, y int, diff int) string {
	// Always take the lowest value as iteration number
	if y < x {
		x, y = y, x
	}

	s := ""
	for i := 0; i < x; i++ {
		s += "+"
	}
	s += "[->"
	for i := 0; i < y; i++ {
		s += "+"
	}
	s += "<]>"

	if diff >= 0 {
		for i := 0; i < diff; i++ {
			s += "+"
		}
	} else {
		for i := 0; i > diff; i-- {
			s += "-"
		}
	}
	return s
}
