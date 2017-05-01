package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

var errBounds = errors.New("can't '<', already minimum value")
var errInvalid = errors.New("invalid token")
var errUnmatch = errors.New("unmatching bracket. ']' without previous '['")
var errReading = errors.New("could not read char")
var errInterrupt = errors.New("interrupted")

type env struct {
	code  string
	vars  []rune
	debug bool
}

func run(e *env) error {
	if stop {
		return errInterrupt
	}
	if e.vars == nil {
		e.vars = make([]rune, 1)
	}

	index := 0

	stdin := bufio.NewReader(os.Stdin)
	reader := strings.NewReader(e.code)

	// ReadRune can only throw EOF
	for c, _, err := reader.ReadRune(); err == nil; c, _, err = reader.ReadRune() {
		if stop {
			return errInterrupt
		}
		switch c {
		case '<':
			if index <= 0 {
				return errBounds
			}
			index--
		case '>':
			index++

			// I usually like to check >= instead of ==, but in this case
			// that'd require a for loop and everything. Not worth it.
			if index == len(e.vars) {
				e.vars = append(e.vars, 0)
			}
		case '+':
			e.vars[index]++
		case '-':
			e.vars[index]--
		case '.':
			char := e.vars[index]
			bytes := make([]byte, utf8.RuneLen(char))
			utf8.EncodeRune(bytes, char)

			if e.debug {
				suffix := " (char #" + strconv.Itoa(int(char)) + ", index [" + strconv.Itoa(index) + "])\n"
				bytes = append(bytes, []byte(suffix)...)
			}

			os.Stdout.Write(bytes)
		case ',':
			e.vars[index], err = getchar(stdin)
			if err != nil {
				return err
			}
		case '[':
			code := ""
			brackets := 0
			i := index

			for c, _, err := reader.ReadRune(); err == nil; c, _, err = reader.ReadRune() {
				if c == '[' {
					brackets++
				} else if c == ']' {
					if brackets <= 0 {
						break
					}
					brackets--
				}
				code += string(c)
			}

			for e.vars[i] != 0 {
				e2 := &env{
					code: code,
					vars: e.vars,
				}
				err := run(e2)
				if err != nil {
					return err
				}
				e.vars = e2.vars
			}
		case ']':
			return errUnmatch
		case ' ', '\n', '\t':
		default:
			return errInvalid
		}
	}
	return nil
}

func getchar(reader *bufio.Reader) (char rune, err error) {
	old, err := terminal.MakeRaw(0)
	// What do I do with the error though?
	// I don't wanna return it...
	if err == nil {
		defer terminal.Restore(0, old)
	}

	char, _, err = reader.ReadRune()
	if err != nil {
		err = errReading
		return
	}

	if char == 3 {
		err = io.EOF
		return
	}
	return
}
