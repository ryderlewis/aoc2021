package day13

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Dot struct {
	x, y int
}

type Fold struct {
	alongX bool
	alongY bool
	value int
}

type Runner struct {
	board map[Dot]bool
	width, height int
	folds []Fold
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.fold(r.folds[0])

	return strconv.Itoa(len(r.board)), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	for _, f := range r.folds {
		r.fold(f)
	}

	r.print()

	return strconv.Itoa(0), nil
}

func (r *Runner) fold(fold Fold) {
	delPoints := make([]Dot, 0)
	addPoints := make([]Dot, 0)

	if fold.alongX {
		r.width = (r.width - 1)/2
		for pos, _ := range r.board {
			if pos.x > fold.value {
				delPoints = append(delPoints, pos)
				addPoints = append(addPoints, Dot{
					x: 2 * fold.value - pos.x,
					y: pos.y,
				})
			}
		}
	} else {
		r.height = (r.height - 1)/2

		for pos, _ := range r.board {
			if pos.y > fold.value {
				delPoints = append(delPoints, pos)
				addPoints = append(addPoints, Dot{
					x: pos.x,
					y: 2 * fold.value - pos.y,
				})
			}
		}
	}

	for _, pos := range delPoints {
		delete(r.board, pos)
	}

	for _, pos := range addPoints {
		r.board[pos] = true
	}
}

func (r *Runner) print() {
	buf := make([]rune, r.width)

	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			pos := Dot{x: x, y: y}
			if _, contains := r.board[pos]; contains {
				buf[x] = '#'
			} else {
				buf[x] = ' '
			}
		}
		fmt.Println(string(buf))
	}
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.board = make(map[Dot]bool)
	r.folds = make([]Fold, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, ",") {
			tokens := strings.Split(line, ",")
			x, _ := strconv.Atoi(tokens[0])
			y, _ := strconv.Atoi(tokens[1])
			dot := Dot{
				x: x,
				y: y,
			}
			if x + 1 > r.width {
				r.width = x + 1
			}
			if y + 1 > r.height {
				r.height = y + 1
			}
			r.board[dot] = true
		} else if strings.Contains(line, "fold along") {
			alongX := strings.Contains(line, "x=")
			alongY := strings.Contains(line, "y=")
			tokens := strings.Split(line, "=")
			value, _ := strconv.Atoi(tokens[1])
			fold := Fold{
				value: value,
				alongX: alongX,
				alongY: alongY,
			}
			r.folds = append(r.folds, fold)
		}
	}

	return nil
}
