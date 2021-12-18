package day19

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Runner struct {

}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(0), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(0), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		scanner.Text()
	}

	return nil
}
