package day02

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type inst struct {
	direction string
	value     int
}
type Runner struct {
	instructions []inst
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	h := 0
	d := 0

	for _, inst := range r.instructions {
		switch inst.direction {
		case "forward":
			h += inst.value
		case "down":
			d += inst.value
		case "up":
			d -= inst.value
		default:
			return "", fmt.Errorf("Invalid direction: %v", inst.direction)
		}
	}

	return strconv.Itoa(h * d), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	h := 0
	d := 0
	a := 0

	for _, inst := range r.instructions {
		switch inst.direction {
		case "forward":
			h += inst.value
			d += a * inst.value
		case "down":
			a += inst.value
		case "up":
			a -= inst.value
		default:
			return "", fmt.Errorf("Invalid direction: %v", inst.direction)
		}
	}

	return strconv.Itoa(h * d), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.instructions = make([]inst, 0)

	for scanner.Scan() {
		tokens := strings.Split(scanner.Text(), " ")
		if len(tokens) != 2 {
			return fmt.Errorf("Error parsing line: %v\n", scanner.Text())
		}

		if i, err := strconv.Atoi(tokens[1]); err != nil {
			return err
		} else {
			r.instructions = append(r.instructions, inst{direction: tokens[0], value: i})
		}
	}

	return nil
}
