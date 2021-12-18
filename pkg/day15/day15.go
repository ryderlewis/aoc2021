package day15

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"sort"
	"strconv"
)

type Runner struct {
	grid     [][]int
	path     []Coordinate
	origGrid [][]int
}

type Coordinate struct {
	y, x int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.embiggify(1)
	risk := r.getRisk()
	r.print()

	return strconv.Itoa(risk), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.embiggify(5)
	risk := r.getRisk()
	r.print()

	return strconv.Itoa(risk), nil
}

func (r *Runner) embiggify(n int) {
	cols := len(r.grid[0])
	rows := len(r.grid)

	newGrid := make([][]int, rows*n)
	for i := 0; i < rows*n; i++ {
		row := make([]int, cols*n)
		newGrid[i] = row
	}

	for dx := 0; dx < n; dx++ {
		for dy := 0; dy < n; dy++ {
			delta := dx + dy
			yOffset := rows * dy
			xOffset := cols * dx

			for x := 0; x < cols; x++ {
				for y := 0; y < rows; y++ {
					origVal := r.grid[y][x]
					newVal := origVal + delta
					if newVal > 9 {
						newVal %= 9
					}

					newGrid[y+yOffset][x+xOffset] = newVal
				}
			}
		}
	}

	r.grid = newGrid
	r.origGrid = r.copy(newGrid)
}

func (r *Runner) copy(grid [][]int) [][]int {
	c := make([][]int, len(grid))
	for y := 0; y < len(grid); y++ {
		c[y] = make([]int, len(grid[0]))
	}

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[0]); x++ {
			c[y][x] = grid[y][x]
		}
	}

	return c
}

func (r *Runner) print() {
	coordinates := make(map[Coordinate]bool)
	for _, p := range r.path {
		coordinates[p] = true
	}

	for y := 0; y < len(r.origGrid); y++ {
		for x := 0; x < len(r.origGrid[0]); x++ {
			if _, ok := coordinates[Coordinate{x: x, y: y}]; ok {
				fmt.Printf("\x1b[31m%d\x1b[0m", r.origGrid[y][x])
			} else {
				fmt.Printf("%d", r.origGrid[y][x])
			}
		}
		fmt.Println()
	}
}

func (r *Runner) getRisk() int {
	maxY := len(r.grid) - 1
	maxX := len(r.grid[0]) - 1

	// initialize weights to max int (except for starting point)
	for y := 0; y < len(r.grid); y++ {
		for x := 0; x < len(r.grid[0]); x++ {
			r.grid[y][x] = math.MaxInt
		}
	}
	r.grid[0][0] = 0

	// keep going until the final point is visited
	cur := &Coordinate{x: 0, y: 0}
	visited := make(map[Coordinate]bool)

	precursors := make(map[Coordinate]*Coordinate)
	precursors[*cur] = nil

	unvisited := make(map[Coordinate]bool)

	for {
		visited[*cur] = true

		// find all unvisited neighbors, updating their weights
		neighbors := make([]Coordinate, 0)
		for _, delta := range []Coordinate{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
			c := Coordinate{x: cur.x + delta.x, y: cur.y + delta.y}
			if c.y < 0 || c.y > maxY {
				continue
			}
			if c.x < 0 || c.x > maxX {
				continue
			}
			if _, v := visited[c]; v {
				continue
			}
			neighbors = append(neighbors, c)

			if r.grid[c.y][c.x] > r.grid[cur.y][cur.x]+r.origGrid[c.y][c.x] {
				precursors[c] = cur
				r.grid[c.y][c.x] = r.grid[cur.y][cur.x] + r.origGrid[c.y][c.x]
			}

			unvisited[c] = true
		}

		if len(unvisited) == 0 {
			break
		}

		// find unvisited with lowest weight, and make that the new cur
		x := make([]Coordinate, 0, len(unvisited))
		for c, _ := range unvisited {
			x = append(x, c)
		}
		sort.Slice(x, func(i, j int) bool {
			ni := x[i]
			nj := x[j]
			return r.grid[ni.y][ni.x] < r.grid[nj.y][nj.x]
		})
		cur = &x[0]
		delete(unvisited, x[0])
	}

	r.path = make([]Coordinate, 0)
	cur = &Coordinate{y: maxY, x: maxX}
	for cur != nil {
		r.path = append(r.path, *cur)
		cur = precursors[*cur]
	}

	// reverse the path?

	return r.grid[maxY][maxX]
}

func (r *Runner) readInput(input io.Reader) error {
	r.grid = make([][]int, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		r.grid = append(r.grid, row)

		for i := 0; i < len(line); i++ {
			row[i] = int(line[i] - '0')
		}
	}

	return nil
}
