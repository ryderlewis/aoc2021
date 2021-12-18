package day09

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"sort"
	"strconv"
)

type Runner struct {
	grid [][]int
}

type Pair struct {
	row, col int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	risk := 0
	for _, pair := range r.lowPairs() {
		risk += r.grid[pair.row][pair.col] + 1
	}

	return strconv.Itoa(risk), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	basinSizes := make([]int, 0)

	for _, pair := range r.lowPairs() {
		// go outward from basin
		basin := make(map[Pair]bool)
		toTry := []Pair{pair}

		for len(toTry) > 0 {
			p := toTry[0]
			toTry = toTry[1:]

			basin[p] = true

			for dr := -1; dr <= 1; dr++ {
				for dc := -1; dc <= 1; dc++ {
					if dr == 0 && dc == 0 {
						continue
					}
					if dr != 0 && dc != 0 {
						continue
					}
					if r.posValid(p.row+dr, p.col+dc) {
						posPair := Pair{row: p.row+dr, col: p.col+dc}
						if r.grid[p.row+dr][p.col+dc] > r.grid[p.row][p.col] && r.grid[p.row+dr][p.col+dc] != 9 {
							if _, tried := basin[posPair]; !tried {
								toTry = append(toTry, posPair)
							}
							basin[posPair] = true
						}
					}
				}
			}
		}

		basinLen := 0
		for _, ok := range basin {
			if ok {
				basinLen++
			}
		}
		basinSizes = append(basinSizes, basinLen)
		fmt.Printf("%#v: %d\n", pair, basinLen)
	}

	sort.Ints(basinSizes)
	ret := 1
	for _, v := range basinSizes[len(basinSizes)-3:] {
		ret *= v
	}
	return strconv.Itoa(ret), nil
}

func (r *Runner) pairValid(p Pair) bool {
	if p.row < 0 || p.row > len(r.grid) - 1 {
		return false
	}
	if p.col < 0 || p.col > len(r.grid[0]) - 1 {
		return false
	}
	return true
}

func (r *Runner) posValid(row, col int) bool {
	return r.pairValid(Pair{row: row, col: col})
}

func (r *Runner) lowPairs() []Pair {
	pairs := make([]Pair, 0)

	for row := 0; row < len(r.grid); row++ {
		for col := 0; col < len(r.grid[0]); col++ {
			isLow := true
			if r.posValid(row-1, col) && r.grid[row][col] >= r.grid[row-1][col] {
				isLow = false
			}
			if r.posValid(row+1, col) && r.grid[row][col] >= r.grid[row+1][col] {
				isLow = false
			}
			if r.posValid(row, col-1) && r.grid[row][col] >= r.grid[row][col-1] {
				isLow = false
			}
			if r.posValid(row, col+1) && r.grid[row][col] >= r.grid[row][col+1] {
				isLow = false
			}
			if isLow {
				pairs = append(pairs, Pair{row: row, col: col})
			}
		}
	}

	return pairs
}

func (r *Runner) readInput(input io.Reader) error {
	r.grid = make([][]int, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		vals := make([]int, len(line))
		for i := 0; i < len(line); i++ {
			val, _ := strconv.Atoi(line[i:i+1])
			vals[i] = val
		}
		r.grid = append(r.grid, vals)
	}

	return nil
}
