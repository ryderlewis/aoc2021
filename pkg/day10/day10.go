package day10

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"sort"
	"strconv"
)

type Runner struct {
	lines []string
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	scores := map[rune]int{
		')': 3,
		']': 57,
		'}': 1197,
		'>': 25137,
	}

	score := 0
	for _, line := range r.lines {
		stack := make([]rune, 0)
		for i := 0; i < len(line); i++ {
			r := rune(line[i])
			syntaxError := false

			switch r {
			case '(':
				stack = append(stack, r)
			case '[':
				stack = append(stack, r)
			case '{':
				stack = append(stack, r)
			case '<':
				stack = append(stack, r)
			case ')':
				if len(stack) == 0 || stack[len(stack)-1] != '(' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case ']':
				if len(stack) == 0 || stack[len(stack)-1] != '[' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case '}':
				if len(stack) == 0 || stack[len(stack)-1] != '{' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case '>':
				if len(stack) == 0 || stack[len(stack)-1] != '<' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			}

			if syntaxError {
				score += scores[r]
				break
			}
		}
	}

	return strconv.Itoa(score), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	scores := map[rune]int{
		'(': 1,
		'[': 2,
		'{': 3,
		'<': 4,
	}

	scoreList := make([]int, 0)
	for _, line := range r.lines {
		stack := make([]rune, 0)
		syntaxError := false

		for i := 0; i < len(line); i++ {
			r := rune(line[i])

			switch r {
			case '(':
				stack = append(stack, r)
			case '[':
				stack = append(stack, r)
			case '{':
				stack = append(stack, r)
			case '<':
				stack = append(stack, r)
			case ')':
				if len(stack) == 0 || stack[len(stack)-1] != '(' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case ']':
				if len(stack) == 0 || stack[len(stack)-1] != '[' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case '}':
				if len(stack) == 0 || stack[len(stack)-1] != '{' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			case '>':
				if len(stack) == 0 || stack[len(stack)-1] != '<' {
					syntaxError = true
				} else {
					stack = stack[:len(stack)-1]
				}
			}

			if syntaxError {
				break
			}
		}

		if !syntaxError {
			score := 0
			for len(stack) > 0 {
				lastRune := stack[len(stack)-1]
				score *= 5
				score += scores[lastRune]
				stack = stack[:len(stack)-1]
			}
			scoreList = append(scoreList, score)
		}
	}

	sort.Ints(scoreList)
	score := scoreList[len(scoreList)/2]

	return strconv.Itoa(score), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.lines = make([]string, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		r.lines = append(r.lines, scanner.Text())
	}

	return nil
}
