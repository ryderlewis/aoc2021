package day06

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Runner struct {
	fish []int
	memo [257]int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	sum := 0
	for _, f := range r.fish {
		sum += r.calcMemo(80 - f)
	}

	return strconv.Itoa(sum), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	sum := 0
	for _, f := range r.fish {
		sum += r.calcMemo(256 - f)
	}

	return strconv.Itoa(sum), nil
}

func (r *Runner) calcMemo(days int) int {
	if days < 0 {
		return 1
	}

	if r.memo[days] != 0 {
		return r.memo[days]
	}

	if days == 0 {
		r.memo[days] = 1
		return 1
	}

	fishCount := 1 // count the current fish
	for spawnDay := days-1; spawnDay >= 0; spawnDay -= 7 {
		// count the fish that this new fish spawns
		fishCount += r.calcMemo(spawnDay - 8)
	}

	r.memo[days] = fishCount
	return fishCount
}

func (r *Runner) readInput(input io.Reader) error {
	r.fish = make([]int, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), ",")
		for _, f := range fields {
			i, _ := strconv.Atoi(f)
			r.fish = append(r.fish, i)
		}
	}

	return nil
}
