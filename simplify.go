package main

import "strconv"

func simplify(code string) (output string) {
	indent := ""

	jumps := 0
	for i, c := range code {
		if jumps > 0 {
			jumps--
			continue
		}

		switch c {
		case '<':
			repeats := 0
			i++
			for i < len(code) && code[i] == '<' {
				repeats++
				i++
			}

			if repeats == 0 {
				output += indent + "i--\n"
			} else {
				jumps = repeats
				output += indent + "c[i] -= " + strconv.Itoa(repeats+1) + "\n"
			}
		case '>':
			repeats := 0
			i++
			for i < len(code) && code[i] == '>' {
				repeats++
				i++
			}

			if repeats == 0 {
				output += indent + "i++\n"
			} else {
				jumps = repeats
				output += indent + "i += " + strconv.Itoa(repeats+1) + "\n"
			}
		case '+':
			repeats := 0
			i++
			for i < len(code) && code[i] == '+' {
				repeats++
				i++
			}

			if repeats == 0 {
				output += indent + "c[i]++\n"
			} else {
				jumps = repeats
				output += indent + "c[i] += " + strconv.Itoa(repeats+1) + "\n"
			}
		case '-':
			repeats := 0
			i++
			for i < len(code) && code[i] == '-' {
				repeats++
				i++
			}

			if repeats == 0 {
				output += indent + "c[i]--\n"
			} else {
				jumps = repeats
				output += indent + "c[i] -= " + strconv.Itoa(repeats+1) + "\n"
			}
		case '.':
			output += indent + "print(c[i])\n"
		case ',':
			output += indent + "c[i] = getchar()\n"
		case '[':
			output += indent + "while (c[i] != 0) {\n"
			indent += "\t"
		case ']':
			if len(indent) > 0 {
				indent = indent[:len(indent)-1]
			}
			output += indent + "}\n"
		}
	}
	return
}
