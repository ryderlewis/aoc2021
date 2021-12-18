package day05

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Runner struct {
	lines []*Line
}

type Line struct {
	x1, y1, x2, y2 int
}

var _ challenge.DailyChallenge = &Runner{}

func maxInt(vals... int) int {
	m := vals[0]
	for _, x := range vals {
		if x > m {
			m = x
		}
	}

	return m
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	maxVal := 0
	for _, line := range r.lines {
		maxVal = maxInt(line.x1, line.y1, line.x2, line.y2, maxVal)
	}
	rows := maxVal + 1
	cols := maxVal + 1

	grid := make([]int, rows * cols)

	for _, line := range r.lines {
		if line.x1 == line.x2 {
			x := line.x1
			for y := line.y1; y <= line.y2; y++ {
				grid[x + cols * y]++
			}
		} else if line.y1 == line.y2 {
			y := line.y1
			for x := line.x1; x <= line.x2; x++ {
				grid[x + cols * y]++
			}
		}
	}

	count := 0
	for _, val := range grid {
		if val > 1 {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	maxVal := 0
	for _, line := range r.lines {
		maxVal = maxInt(line.x1, line.y1, line.x2, line.y2, maxVal)
	}
	rows := maxVal + 1
	cols := maxVal + 1

	grid := make([]int, rows * cols)

	for _, line := range r.lines {
		if line.x1 == line.x2 {
			x := line.x1
			for y := line.y1; y <= line.y2; y++ {
				grid[x + cols * y]++
			}
		} else if line.y1 == line.y2 {
			y := line.y1
			for x := line.x1; x <= line.x2; x++ {
				grid[x + cols * y]++
			}
		} else {
			x, y := line.x1, line.y1
			dx, dy := 1, 1
			if line.x2 < line.x1 {
				dx = -1
			}
			if line.y2 < line.y1 {
				dy = -1
			}

			grid[x + cols * y]++
			for {
				x += dx
				y += dy
				grid[x + cols * y]++

				if x == line.x2 {
					break
				}
			}
		}
	}

	count := 0
	for _, val := range grid {
		if val > 1 {
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.lines = make([]*Line, 0)

	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		p1 := strings.Split(tokens[0], ",")
		p2 := strings.Split(tokens[2], ",")

		line := &Line{}
		line.x1, _ = strconv.Atoi(p1[0])
		line.y1, _ = strconv.Atoi(p1[1])
		line.x2, _ = strconv.Atoi(p2[0])
		line.y2, _ = strconv.Atoi(p2[1])

		if line.x1 > line.x2 || line.y1 > line.y2 {
			line.x1, line.x2 = line.x2, line.x1
			line.y1, line.y2 = line.y2, line.y1
		}

		r.lines = append(r.lines, line)
	}

	return nil
}
