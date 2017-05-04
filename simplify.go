package main

type simplifier interface {
	simplify(string, *string, int, int) (string, int)
	finalize(string) string
}

func simplify(code string, s simplifier) (output string) {
	indent := new(string)

	skip := 0
	for i := range code {
		if skip > 0 {
			skip--
			continue
		}

		repeats := 0
		for i+repeats+1 < len(code) && code[i+repeats+1] == code[i] {
			repeats++
		}

		var line string
		line, skip = s.simplify(code, indent, i, repeats)

		if line == "" {
			continue
		}
		output += line + "\n"
	}
	return s.finalize(output)
}
