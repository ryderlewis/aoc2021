package day01

import (
	"bufio"
	"io"
	"strconv"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
)

type Runner struct {
	values []int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	lastVal := -1
	count := 0
	for idx, val := range r.values {
		if idx > 0 && val > lastVal {
			count++
		}
		lastVal = val
	}
	
	return strconv.Itoa(count), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}
	
	sum := 0
	count := 0
	for idx, val := range r.values {
		lastSum := sum
		sum += val
		if idx > 2 {
			sum -= r.values[idx-3]
			if sum > lastSum {
				count++
			}
		}
	}

	return strconv.Itoa(count), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.values = make([]int, 0)
	
	for scanner.Scan() {
		if i, err := strconv.Atoi(scanner.Text()); err != nil {
			return err
		} else {
			r.values = append(r.values, i)
		}
	}
	
	return nil
}
