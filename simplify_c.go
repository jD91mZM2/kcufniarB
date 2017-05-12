package main

import "strconv"

type csimplifier struct{}

func (s *csimplifier) simplify(code string, indent *string, i int, repeats int) (string, int) {
	if i == 0 && *indent == "" {
		*indent = "\t"
	}

	switch code[i] {
	case '<':
		if repeats == 0 {
			return *indent + "i--;", 0
		}
		return *indent + "i-=" + strconv.Itoa(repeats+1) + ";", repeats
	case '>':
		if repeats == 0 {
			return *indent + "i++;", 0
		}
		return *indent + "i+=" + strconv.Itoa(repeats+1) + ";", repeats
	case '+':
		if repeats == 0 {
			return *indent + "*i++;", 0
		}
		return *indent + "*i+=" + strconv.Itoa(repeats+1) + ";", repeats
	case '-':
		if repeats == 0 {
			return *indent + "*i--;", 0
		}
		return *indent + "*i-=" + strconv.Itoa(repeats+1) + ";", repeats
	case '.':
		return *indent + "putchar(*i);", 0
	case ',':
		return *indent + "*i=getchar();", 0
	case '[':
		if i+4 < len(code) && code[i+1] == '-' && code[i+2] == '>' && code[i+3] == '+' {
			factor := 0
			for i+factor+3 < len(code) && code[i+factor+3] == '+' {
				factor++
			}

			if i+factor+4 < len(code) && code[i+factor+3] == '<' && code[i+factor+4] == ']' {
				return *indent + "*(i+1)=*i*" + strconv.Itoa(factor), factor + 4
			}
		}

		str := *indent + "while (*i)\n" + *indent + "{"
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

func (s *csimplifier) finalize(code string) string {
	return `#include <stdio.h>

int main() {
	// Here would be an 'i' array with enough memory allocated.
` + code + `}
`
}
