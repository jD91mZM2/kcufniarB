package main

import "strconv"

type normsimplifier struct{}

func (s *normsimplifier) simplify(token rune, indent *string, repeats int) (string, bool) {
	switch token {
	case '<':
		if repeats == 0 {
			return *indent + "i--", false
		}
		return *indent + "i -= " + strconv.Itoa(repeats+1) + "", true
	case '>':
		if repeats == 0 {
			return *indent + "i++", false
		}
		return *indent + "i += " + strconv.Itoa(repeats+1) + "", true
	case '+':
		if repeats == 0 {
			return *indent + "c[i]++", false
		}
		return *indent + "c[i] += " + strconv.Itoa(repeats+1) + "", true
	case '-':
		if repeats == 0 {
			return *indent + "c[i]--", false
		}
		return *indent + "c[i] -= " + strconv.Itoa(repeats+1) + "", true
	case '.':
		return *indent + "print(c[i])", false
	case ',':
		return *indent + "c[i] = getchar()", false
	case '[':
		str := *indent + "while (c[i] != 0) {"
		*indent += "\t"
		return str, false
	case ']':
		if len(*indent) > 0 {
			*indent = (*indent)[:len(*indent)-1]
		}
		return *indent + "}", false
	}
	return "", false
}

func (s *normsimplifier) finalize(output string) string {
	return output
}
