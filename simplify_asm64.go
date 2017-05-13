package main

import "strconv"

type asm64simplifier struct {
	num int
}

func (s *asm64simplifier) simplify(code string, indent *string, i int, repeats int) (string, int) {
	*indent = "\t"

	switch code[i] {
	case '<':
		return *indent + "pop %rbx", 0
	case '>':
		return *indent + "push %rbx", 0
	case '+':
		if repeats > 0 {
			return *indent + "add $" + strconv.Itoa(repeats+1) + ", %rbx", repeats
		}
		return *indent + "inc %rbx", 0
	case '-':
		if repeats > 0 {
			return *indent + "sub $" + strconv.Itoa(repeats+1) + ", %rbx", repeats
		}
		return *indent + "dec %rbx", 0
	case '.':
		return "\n" + *indent + "mov %rbx, char\n" +
			*indent + "mov $char, %rsi\n" +
			*indent + "mov $1, %rax\n" +
			*indent + "mov $1, %rdi\n" +
			*indent + "mov $1, %rdx\n" +
			*indent + "syscall", 0
	case ',':
		return *indent + "// reading not implemented", 0
	case '[':
		s.num++
		str := strconv.Itoa(s.num)
		return "\n" +
			"loop" + str + ":\n" +
			*indent + "cmp $0, %rbx\n" +
			*indent + "jne loop" + str + "_end\n", 0
	case ']':
		str := strconv.Itoa(s.num)
		return *indent + "jmp loop " + str + "\n" +
			"loop" + str + "_end:", 0
	}
	return "", 0
}

func (s *asm64simplifier) finalize(output string) string {
	return `.text
	.global _start

_start:
` + output + `
	mov $60, %rax
	mov $0, %rdi
	syscall

.data
char:
	.space 1
`
}
