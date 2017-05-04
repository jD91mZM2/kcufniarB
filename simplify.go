package main

type simplifier interface {
	simplify(string, *string, int, int) (string, bool)
	finalize(string) string
}

func simplify(code string, s simplifier) (output string) {
	indent := new(string)

	jumps := 0
	for i := range code {
		if jumps > 0 {
			jumps--
			continue
		}

		repeats := 0
		for i+repeats+1 < len(code) && code[i+repeats+1] == code[i] {
			repeats++
		}

		line, skip := s.simplify(code, indent, i, repeats)
		if line == "" {
			continue
		}
		output += line + "\n"

		if skip {
			jumps = repeats
		}
	}
	return
}
