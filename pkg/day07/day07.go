package day07

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Runner struct {
	crabs []int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	sort.Ints(r.crabs)
	median := r.crabs[len(r.crabs)/2]

	dist := 0
	for _, c := range r.crabs {
		if median > c {
			dist += median - c
		} else {
			dist += c - median
		}
	}

	return strconv.Itoa(dist), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	sort.Ints(r.crabs)
	fuel := make([]int, r.crabs[len(r.crabs)-1] + 1)
	for i := 1; i < len(fuel); i++ {
		fuel[i] = i + fuel[i-1]
	}

	var bestFuel int
	for loc := r.crabs[0]; loc <= r.crabs[len(r.crabs)-1]; loc++ {
		f := 0
		for _, c := range r.crabs {
			if c < loc {
				f += fuel[loc-c]
			} else {
				f += fuel[c-loc]
			}
		}

		if loc == r.crabs[0] || f < bestFuel {
			bestFuel = f
		}
	}

	return strconv.Itoa(bestFuel), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.crabs = make([]int, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		for _, s := range strings.Split(scanner.Text(), ",") {
			i, _ := strconv.Atoi(s)
			r.crabs = append(r.crabs, i)
		}
	}

	return nil
}
