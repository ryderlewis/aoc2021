package day12

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Cave struct {
	name string
	isSmall bool
	neighbors []*Cave
}

type Runner struct {
	caves map[string]*Cave
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// do a dfs through neighbors
	visited := make(map[string]bool)
	curr := r.caves["start"]

	count := r.countPaths(curr, visited)

	return strconv.Itoa(count), nil
}

func (r *Runner) countPaths(curr *Cave, visited map[string]bool) int {
	count := 0

	for _, neighbor := range curr.neighbors {
		if neighbor.name == "end" {
			count++
			continue
		}

		if neighbor.isSmall {
			if _, alreadyVisited := visited[neighbor.name]; alreadyVisited {
				continue
			}
		}

		visited[curr.name] = true
		count += r.countPaths(neighbor, visited)
		delete(visited, curr.name)
	}

	return count
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// do a dfs through neighbors
	visited := make(map[string]bool)
	curr := r.caves["start"]

	count := r.countPaths2(curr, nil, visited)

	return strconv.Itoa(count), nil
}

func (r *Runner) countPaths2(curr, double *Cave, visited map[string]bool) int {
	count := 0

	for _, neighbor := range curr.neighbors {
		if neighbor.name == "end" {
			if double == nil {
				count++
			} else {
				// if double is true, make sure the double was visited a 2nd time, otherwise this is a dup
				if _, doubleVisited := visited[double.name]; doubleVisited || curr == double {
					count++
				}
			}
			continue
		}

		if neighbor.isSmall {
			if _, alreadyVisited := visited[neighbor.name]; alreadyVisited {
				continue
			}
		}

		// try once visiting the current cave
		visited[curr.name] = true
		count += r.countPaths2(neighbor, double, visited)
		delete(visited, curr.name)

		// try again if the current cave is small, and there is not a double, making it the double
		if curr.isSmall && curr.name != "start" && curr.name != "end" && double == nil {
			count += r.countPaths2(neighbor, curr, visited)
		}
	}

	return count
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.caves = make(map[string]*Cave)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "-")
		for _, n := range s {
			if _, exists := r.caves[n]; !exists {
				r.caves[n] = &Cave{
					name: n,
					isSmall: strings.ToLower(n) == n,
					neighbors: make([]*Cave, 0),
				}
			}
		}

		r.caves[s[0]].neighbors = append(r.caves[s[0]].neighbors, r.caves[s[1]])
		r.caves[s[1]].neighbors = append(r.caves[s[1]].neighbors, r.caves[s[0]])
	}

	return nil
}
