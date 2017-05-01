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
	index int
	debug bool
}

func run(e *env) error {
	if stop {
		return errInterrupt
	}
	if e.vars == nil {
		e.vars = make([]rune, 1)
	}

	stdin := bufio.NewReader(os.Stdin)
	reader := strings.NewReader(e.code)

	// ReadRune can only throw EOF
	for c, _, err := reader.ReadRune(); err == nil; c, _, err = reader.ReadRune() {
		if stop {
			return errInterrupt
		}
		switch c {
		case '<':
			if e.index <= 0 {
				return errBounds
			}
			e.index--

			if e.debug {
				os.Stdout.Write([]byte("Index now on #" + strconv.Itoa(e.index) + "\n"))
			}
		case '>':
			e.index++

			// I usually like to check >= instead of ==, but in this case
			// that'd require a for loop and everything. Not worth it.
			if e.index == len(e.vars) {
				e.vars = append(e.vars, 0)
			}

			if e.debug {
				os.Stdout.Write([]byte("Index now on #" + strconv.Itoa(e.index) + "\n"))
			}
		case '+':
			e.vars[e.index]++
		case '-':
			e.vars[e.index]--
		case '.':
			char := e.vars[e.index]
			len := utf8.RuneLen(char)

			if len < 0 {
				continue
			}

			bytes := make([]byte, len)
			utf8.EncodeRune(bytes, char)

			if e.debug {
				suffix := " (char #" + strconv.Itoa(int(char)) + ", e.index [" + strconv.Itoa(e.index) + "])\n"
				bytes = append(bytes, []byte(suffix)...)
			}

			os.Stdout.Write(bytes)
		case ',':
			e.vars[e.index], err = getchar(stdin)
			if err != nil {
				return err
			}
		case '[':
			code := ""
			brackets := 0
			i := e.index

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
				e2 := &env{}
				*e2 = *e
				e2.code = code
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
