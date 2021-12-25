package day25

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Square byte

const (
	// values at start of each round
	EAST = Square('>') // east-facing cucumber
	SOUTH = Square('v') // south-facing cucumber
	EMPTY = Square('.') // empty square

	// temporary values
	_ Square = iota
	MOVEDEAST // occupied by a cucumber that has already moved east one step
	MOVEDSOUTH // occupied by a cucumber that has already moved south one step
	VACATEDEAST // square that was occupied by an east-facing cucumber
	VACATEDSOUTH // square that was occupied by a south-facing cucumber
)

type Runner struct {
	board [][]Square
	width, height int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	steps := 1
	for r.move() {
		steps++
	}

	return strconv.Itoa(steps), nil
}

func (r *Runner) move() bool {
	moved := false

	// move east-facing cucumbers
	for y, row := range r.board {
		for x, square := range row {
			// see if an east-facing cucumber in this square moves
			if square == EAST {
				rightY, rightX := y, (x+1)%r.width
				if r.board[rightY][rightX] == EMPTY {
					r.board[y][x] = VACATEDEAST
					r.board[rightY][rightX] = MOVEDEAST
					moved = true
				}
			}
		}
	}

	// move south-facing cucumbers
	for y, row := range r.board {
		for x, square := range row {
			// see if a south-facing cucumber in this square moves
			if square == SOUTH {
				// see if a south-facing cucumber moves down
				downY, downX := (y+1)%r.height, x
				switch r.board[downY][downX] {
				case EMPTY, VACATEDEAST:
					r.board[y][x] = VACATEDSOUTH
					r.board[downY][downX] = MOVEDSOUTH
					moved = true
				}
			}
		}
	}

	// clean up the board
	for y, row := range r.board {
		for x, square := range row {
			switch square {
			case VACATEDEAST, VACATEDSOUTH:
				r.board[y][x] = EMPTY
			case MOVEDEAST:
				r.board[y][x] = EAST
			case MOVEDSOUTH:
				r.board[y][x] = SOUTH
			}
		}
	}

	return moved
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
		row := []Square(scanner.Text())
		r.width = len(row)
		r.height++

		r.board = append(r.board, row)
	}

	return nil
}
