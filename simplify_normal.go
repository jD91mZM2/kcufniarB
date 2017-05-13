package main

import "strconv"

type normsimplifier struct{}

func (s *normsimplifier) simplify(code string, indent *string, i int, repeats int) (string, int) {
	switch code[i] {
	case '<':
		if repeats > 0 {
			return *indent + "i -= " + strconv.Itoa(repeats+1), repeats
		}
		return *indent + "i--", 0
	case '>':
		if repeats > 0 {
			return *indent + "i += " + strconv.Itoa(repeats+1), repeats
		}
		return *indent + "i++", 0
	case '+':
		if repeats > 0 {
			return *indent + "c[i] += " + strconv.Itoa(repeats+1), repeats
		}
		return *indent + "c[i]++", 0
	case '-':
		if repeats > 0 {
			return *indent + "c[i] -= " + strconv.Itoa(repeats+1), repeats
		}
		return *indent + "c[i]--", 0
	case '.':
		return *indent + "print(c[i])", 0
	case ',':
		return *indent + "c[i] = getchar()", 0
	case '[':
		str := *indent + "while (c[i] != 0) {"
		*indent += "\t"
		return str, 0
	case ']':
		if len(*indent) > 0 {
			*indent = (*indent)[:len(*indent)-1]
		}
		return *indent + "}", 0
	}
	return "", 0
}

func (s *normsimplifier) finalize(output string) string {
	return output
}
