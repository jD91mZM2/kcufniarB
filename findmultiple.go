package main

type multiple struct {
	x     int
	y     int
	value int
}
type valdiff struct {
	value   int
	diff    int
	diffAbs int
}

func findmultiple(n int) (int, int, int) {
	if n < 0 {
		return 0, 0, 0
	}
	arr := make([]*multiple, 0)
	found := false

	half := n / 2
	for x := 1; x < half; x++ {
		for y := 1; y < half; y++ {
			result := x * y
			if result == n {
				if !found {
					arr = []*multiple{&multiple{
						x: x,
						y: y,
					}}
					found = true
					continue
				}
			} else if found {
				continue
			}

			arr = append(arr, &multiple{
				x:     x,
				y:     y,
				value: result,
			})
		}
	}

	if found {
		var lowest *multiple

		for _, multi := range arr {
			if lowest == nil || multi.x+multi.y < lowest.x+lowest.y {
				lowest = multi
			}
		}
		return lowest.x, lowest.y, 0
	}
	var closest *valdiff

	for _, multi := range arr {
		diff := n - multi.value
		diffAbs := diff
		if diffAbs < 0 {
			diffAbs = -diffAbs
		}

		if closest == nil || diffAbs < closest.diffAbs {
			closest = &valdiff{
				value:   multi.value,
				diff:    diff,
				diffAbs: diffAbs,
			}
			continue
		}
	}

	var lowest *multiple

	for _, multi := range arr {
		if multi.value != closest.value {
			continue
		}

		if lowest == nil || multi.x+multi.y < lowest.x+lowest.y {
			lowest = multi
		}
	}

	if lowest == nil {
		return 0, 0, 0
	}

	return lowest.x, lowest.y, closest.diff
}
