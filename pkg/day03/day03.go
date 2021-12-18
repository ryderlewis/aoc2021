package day03

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Runner struct {
	inputs []string
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	gamma := make([]int, len(r.inputs[0]))

	for _, s := range r.inputs {
		for j, c := range s {
			if c == '1' {
				gamma[j] += 1
			}
		}
	}

	gsum := 0
	esum := 0

	for i := 0; i < len(gamma); i++ {
		power := len(gamma) - i - 1

		if gamma[i] > len(r.inputs)/2 {
			gsum += 1 << power
		} else {
			esum += 1 << power
		}
	}

	return strconv.Itoa(gsum * esum), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// split the inputs into oxygen and co2 sets
	oxygen := make([]string, len(r.inputs))
	copy(oxygen, r.inputs)
	for pos := 0; len(oxygen) > 1; pos++ {
		newOxygen := make([]string, 0)
		g := r.gamma(oxygen, pos)
		// fmt.Printf("pos=%d, g=%d, oxygen=%v\n", pos, g, oxygen)
		for _, s := range oxygen {
			if s[pos] == strconv.Itoa(g)[0] {
				newOxygen = append(newOxygen, s)
			}
		}
		oxygen = newOxygen
	}

	co2 := make([]string, len(r.inputs))
	copy(co2, r.inputs)
	for pos := 0; len(co2) > 1; pos++ {
		newCo2 := make([]string, 0)
		e := r.epsilon(co2, pos)
		for _, s := range co2 {
			if s[pos] == strconv.Itoa(e)[0] {
				newCo2 = append(newCo2, s)
			}
		}
		co2 = newCo2
	}

	o, _ := strconv.ParseInt(oxygen[0], 2, 64)
	c, _ := strconv.ParseInt(co2[0], 2, 64)
	return fmt.Sprintf("%d", o * c), nil
}

func (r *Runner) gamma(inputs []string, index int) int {
	v := r.mostCommon(inputs, index)
	if v < 0 {
		return 1
	}
	return v
}

func (r *Runner) epsilon(inputs []string, index int) int {
	v := r.mostCommon(inputs, index)
	switch v {
	case 1:
		return 0
	case 0:
		return 1
	default:
		return 0
	}
}

func (*Runner) mostCommon(inputs []string, index int) int {
	count := 0

	for _, s := range inputs {
		if s[index] == '1' {
			count += 1
		}
	}
	
	if len(inputs) % 2 == 0 && count == len(inputs)/2 {
		return -1
	} else if count > len(inputs)/2 {
		return 1
	}
	
	return 0
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.inputs = make([]string, 0)

	for scanner.Scan() {
		r.inputs = append(r.inputs, scanner.Text())
	}

	return nil
}
