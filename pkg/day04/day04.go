package day04

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Runner struct {
	values []int
	boards []*Board
}

type Board struct {
	valuePositions map[int]int
	valueCalled    [25]bool
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	for _, val := range r.values {
		fmt.Printf("Calling %d\n", val)
		for _, b := range r.boards {
			if b.Call(val) && b.Winner() {
				return strconv.Itoa(val * b.UncalledSum()), nil
			}
		}
	}

	return "", fmt.Errorf("No winner")
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	lastWin := 0
	for _, val := range r.values {
		fmt.Printf("Calling %d\n", val)
		for _, b := range r.boards {
			if !b.Winner() {
				if b.Call(val) && b.Winner() {
					lastWin = val * b.UncalledSum()
				}
			}
		}
	}

	return strconv.Itoa(lastWin), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.values = make([]int, 0)
	r.boards = make([]*Board, 0)

	scanner.Scan()
	for _, val := range strings.Split(scanner.Text(), ",") {
		i, _ := strconv.Atoi(val)
		r.values = append(r.values, i)
	}

	for scanner.Scan() {
		board := &Board{
			valuePositions: make(map[int]int),
		}
		r.boards = append(r.boards, board)

		for row := 0; row < 5; row++ {
			scanner.Scan()
			nums := strings.Fields(scanner.Text())

			for col := 0; col < 5; col++ {
				i, _ := strconv.Atoi(nums[col])
				board.valuePositions[i] = row*5 + col
			}
		}
	}

	return nil
}

func (b *Board) Call(num int) bool {
	if pos, ok := b.valuePositions[num]; ok {
		b.valueCalled[pos] = true
		return true
	}

	return false
}

func (b *Board) Winner() bool {
	// check horizontals
	for row := 0; row < 5; row++ {
		allTrue := true
		for col := 0; col < 5; col++ {
			if !b.valueCalled[row*5+col] {
				allTrue = false
				break
			}
		}

		if allTrue {
			fmt.Printf("Ryder A\n")
			return true
		}
	}

	// check verticals
	for col := 0; col < 5; col++ {
		allTrue := true
		for row := 0; row < 5; row++ {
			if !b.valueCalled[row*5+col] {
				allTrue = false
				break
			}
		}

		if allTrue {
			fmt.Printf("Ryder B\n")
			return true
		}
	}

	// check diagnals
	/*
	allTrue := true
	for row := 0; row < 5; row++ {
		col := row
		if !b.valueCalled[row*5+col] {
			allTrue = false
			break
		}
	}
	if allTrue {
		fmt.Printf("Ryder C\n")
		return true
	}

	allTrue = true
	for row := 0; row < 5; row++ {
		col := 4 - row
		if !b.valueCalled[row*5+col] {
			allTrue = false
			break
		}
	}
	if allTrue {
		fmt.Printf("Ryder D\n")
		return true
	}
	
	 */

	return false
}

func (b *Board) UncalledSum() int {
	fmt.Printf("Winning board: %#v\n", b)

	sum := 0
	for val, pos := range b.valuePositions {
		if !b.valueCalled[pos] {
			fmt.Printf("Value %d at pos %d uncalled\n", val, pos)
			sum += val
		}
	}

	return sum
}
