package day11

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Pos struct {
	x, y int
}
type Runner struct {
	energy map[Pos]int
	adjacents map[Pos][]Pos
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	flashCount := 0
	for steps := 0; steps < 100; steps++ {
		flashes := make(map[Pos]bool)

		for pos, e := range r.energy {
			r.energy[pos] = e + 1
		}

		for keepGoing := true; keepGoing; {
			keepGoing = false
			for pos, e := range r.energy {
				if e > 9 {
					if _, ok := flashes[pos]; !ok {
						flashes[pos] = true
						keepGoing = true
						for _, adjacent := range r.adjacents[pos] {
							r.energy[adjacent]++
						}
					}
				}
			}
		}

		flashCount += len(flashes)
		for pos, _ := range flashes {
			r.energy[pos] = 0
		}
	}

	return strconv.Itoa(flashCount), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	steps := 0
	for ; ; steps++ {
		flashes := make(map[Pos]bool)

		for pos, e := range r.energy {
			r.energy[pos] = e + 1
		}

		for keepGoing := true; keepGoing; {
			keepGoing = false
			for pos, e := range r.energy {
				if e > 9 {
					if _, ok := flashes[pos]; !ok {
						flashes[pos] = true
						keepGoing = true
						for _, adjacent := range r.adjacents[pos] {
							r.energy[adjacent]++
						}
					}
				}
			}
		}

		for pos, _ := range flashes {
			r.energy[pos] = 0
		}

		if len(flashes) == 100 {
			break
		}
	}

	return strconv.Itoa(steps+1), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.energy = make(map[Pos]int)
	r.adjacents = make(map[Pos][]Pos)

	scanner := bufio.NewScanner(input)
	y := 0
	for scanner.Scan() {
		nums := scanner.Text()
		for x := 0; x < len(nums); x++ {
			pos := Pos{x: x, y: y}
			e := nums[x] - '0'
			r.energy[pos] = int(e)
			r.adjacents[pos] = make([]Pos, 0)

			for dx := -1; dx <= 1; dx++ {
				for dy := -1; dy <= 1; dy++ {
					if dx == 0 && dy == 0 {
						continue
					}
					nx, ny := x + dx, y + dy
					if nx < 0 || nx > 9 || ny < 0 || ny > 9 {
						continue
					}
					npos := Pos{x: nx, y: ny}
					r.adjacents[pos] = append(r.adjacents[pos], npos)
				}
			}
		}
		y++
	}

	return nil
}
