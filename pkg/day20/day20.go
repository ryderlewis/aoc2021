package day20

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Coordinate struct {
	x, y int
}

type Runner struct {
	algo [512]int
	grid map[Coordinate]int
	width int
	height int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	for i := 0; i < 2; i++ {
		r.enhance((i%2) & r.algo[0])
	}

	return strconv.Itoa(len(r.grid)), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	for i := 0; i < 50; i++ {
		r.enhance((i%2) & r.algo[0])
	}

	return strconv.Itoa(len(r.grid)), nil
}

func (r *Runner) val(x, y, blink int) int {
	// if the x or y value is "outside" the image, its value will be the "blink" value
	if x < 0 || y < 0 || x >= r.width || y >= r.height {
		return blink
	}

	return r.grid[Coordinate{x, y}]
}

func (r *Runner) enhance(blink int) {
	grid := make(map[Coordinate]int)

	for y := -1; y < r.height + 1; y++ {
		for x := -1; x < r.width + 1; x++ {
			// make an algo index from the grid
			index := r.val(x-1, y-1, blink) << 8 |
				r.val(x, y-1, blink) << 7 |
				r.val(x+1, y-1, blink) << 6 |
				r.val(x-1, y, blink) << 5 |
				r.val(x, y, blink) << 4 |
				r.val(x+1, y, blink) << 3 |
				r.val(x-1, y+1, blink) << 2 |
				r.val(x, y+1, blink) << 1 |
				r.val(x+1, y+1, blink) & 1

			val := r.algo[index]
			if val > 0 {
				grid[Coordinate{x+1, y+1}] = val
			}
		}
	}

	r.grid = grid
	r.width += 2
	r.height += 2
}

func (r *Runner) print() {
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			char := '.'
			if r.grid[Coordinate{x, y}] > 0 {
				char = '#'
			}
			fmt.Printf("%c", char)
		}
		fmt.Println()
	}
}

func (r *Runner) readInput(input io.Reader) error {
	r.grid = make(map[Coordinate]int)

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	algoLine := scanner.Text()
	for i, c := range []rune(algoLine) {
		if c == '#' {
			r.algo[i] = 1
		}
	}

	scanner.Scan() // blank line

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if y == 0 {
			r.width = len(line)
		}

		for x, c := range []rune(line) {
			if c == '#' {
				r.grid[Coordinate{x, y}] = 1
			}
		}
		y++
	}
	r.height = y

	return nil
}
