package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

var errBounds = errors.New("can't '<', already minimum value")
var errInvalid = errors.New("invalid token")
var errEOF = errors.New("EOF hit while searching for matching ']'")
var errUnmatch = errors.New("unmatching bracket. ']' without previous '['")
var errReading = errors.New("could not read char")
var errInterrupt = errors.New("interrupted")

type env struct {
	code  string
	vars  []rune
	index int
	debug bool

	debugCode    string
	debugSpeedup bool
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

		debugDelay := time.Millisecond * 100
		if e.debug {
			e.debugCode += string(c)

			if !e.debugSpeedup {
				clear()
				renderdebugger(e)
			}
		}

		switch c {
		case '<':
			if e.index <= 0 {
				return errBounds
			}
			e.index--
		case '>':
			e.index++

			// I usually like to check >= instead of ==, but in this case
			// that'd require a for loop and everything. Not worth it.
			if e.index == len(e.vars) {
				e.vars = append(e.vars, 0)
			}
		case '+':
			e.vars[e.index]++
		case '-':
			e.vars[e.index]--
		case '.':
			if e.debugSpeedup {
				continue
			}

			char := e.vars[e.index]
			if char < 0 {
				continue
			}

			bytes := make([]byte, utf8.RuneLen(char))
			utf8.EncodeRune(bytes, char)

			if e.debug {
				suffix := " (char #" + strconv.Itoa(int(char)) + ", e.index [" + strconv.Itoa(e.index) + "])"
				bytes = append(bytes, []byte(suffix)...)

				debugDelay = time.Second * 3
			}
			os.Stdout.Write(bytes)
		case ',':
			if e.debug && !e.debugSpeedup {
				os.Stdout.Write([]byte("Awaiting input...\n"))
			}

			e.vars[e.index], err = getchar(stdin)
			if err != nil {
				return err
			}
		case '!':
			if !e.debug {
				return errInvalid
			}

			if e.debugSpeedup {
				clear()
				renderdebugger(e)

				e.debugSpeedup = false
			}

			s := "["
			for i, c := range e.vars {
				if i > 0 {
					s += ", "
				}
				s += strconv.Itoa(i) + ": "

				if c >= 0 {
					s += string(c)
				}
				s += " (#" + strconv.Itoa(int(c)) + ")"
			}
			s += "]\n"

			os.Stdout.Write([]byte(s))
			os.Stdout.Write([]byte("[c]ontinue, [s]peedup"))

			char, err := getchar(stdin)
			if err != nil {
				return err
			}
			for {
				switch char {
				case 'c':
				case 's':
					e.debugSpeedup = true
				default:
					continue
				}
				break
			}
		case '[':
			code := ""
			brackets := 0
			found := false

			for c, _, err := reader.ReadRune(); err == nil; c, _, err = reader.ReadRune() {
				if c == '[' {
					brackets++
				} else if c == ']' {
					if brackets <= 0 {
						found = true
						break
					}
					brackets--
				}
				code += string(c)
			}

			if !found {
				return errEOF
			}

			for e.vars[e.index] != 0 {
				backup := e.code

				e.code = code
				err := run(e)
				if err != nil {
					return err
				}

				e.code = backup
			}
		case ']':
			return errUnmatch
		case ' ', '\n', '\t':
		default:
			return errInvalid
		}

		if e.debug && !e.debugSpeedup {
			time.Sleep(debugDelay)
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

func renderdebugger(e *env) {
	os.Stdout.Write([]byte("Index: #" + strconv.Itoa(e.index) + "\n"))
	os.Stdout.Write([]byte("Value: #" + strconv.Itoa(int(e.vars[e.index])) + "\n\n"))
	os.Stdout.Write([]byte(e.debugCode + "\n"))
}
