package main

type simplifier interface {
	simplify(rune, *string, int) (string, bool)
	finalize(string) string
}

func simplify(code string, s simplifier) (output string) {
	indent := new(string)

	jumps := 0
	for i, c := range code {
		if jumps > 0 {
			jumps--
			continue
		}

		repeats := 0
		for i+repeats+1 < len(code) && rune(code[i+repeats+1]) == c {
			repeats++
		}

		line, skip := s.simplify(c, indent, repeats)
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
