package day14

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	val  rune
	next *Node
}
type Runner struct {
	head  *Node
	rules map[string]rune
	pairs map[string]int
	first rune
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	counts := make(map[rune]int)
	for i := 0; i < 10; i++ {
		r.polymerize()
	}

	for curr := r.head; curr != nil; curr = curr.next {
		counts[curr.val]++
	}

	vals := make([]int, 0)
	for _, count := range counts {
		vals = append(vals, count)
	}
	sort.Ints(vals)

	return strconv.Itoa(vals[len(vals)-1] - vals[0]), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	counts := make(map[rune]int)
	counts[r.first] = 1
	for i := 0; i < 40; i++ {
		r.polymerize2()
	}

	for p, c := range r.pairs {
		counts[rune(p[1])] += c
	}

	vals := make([]int, 0)
	for val, count := range counts {
		fmt.Printf("%c: %d\n", val, count)
		if count > 0 {
			vals = append(vals, count)
		}
	}
	sort.Ints(vals)

	return strconv.Itoa(vals[len(vals)-1] - vals[0]), nil
}

func (r *Runner) polymerize() {
	var curr, next *Node
	var buf strings.Builder
	buf.Grow(2)

	curr = r.head
	next = curr.next

	for next != nil {
		buf.Reset()
		buf.WriteRune(curr.val)
		buf.WriteRune(next.val)

		if insert, ok := r.rules[buf.String()]; ok {
			mid := &Node{
				val: insert,
				next: next,
			}
			curr.next = mid
		}

		curr = next
		next = curr.next
	}
}

func (r *Runner) polymerize2() {
	add := make(map[string]int)
	sub := make(map[string]int)

	for p, c := range r.pairs {
		if insert, ok := r.rules[p]; ok {
			sub[p] += c
			add[fmt.Sprintf("%c%c", p[0], insert)] += c
			add[fmt.Sprintf("%c%c", insert, p[1])] += c
		}
	}

	for p, c := range add {
		r.pairs[p] += c
	}

	for p, c := range sub {
		r.pairs[p] -= c
	}
}

func (r *Runner) readInput(input io.Reader) error {
	r.rules = make(map[string]rune)
	r.pairs = make(map[string]int)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()

		if r.head == nil {
			var curr, last *Node
			for i := 0; i < len(line); i++ {
				curr = &Node{val: rune(line[i])}
				if i == 0 {
					r.head = curr
				} else {
					last.next = curr
				}
				last = curr
			}

			for i := 0; i < len(line)-1; i++ {
				pair := line[i:i+2]
				r.pairs[pair]++
			}
			r.first = rune(line[0])
		} else if strings.Contains(line, "->") {
			tokens := strings.Split(line, " -> ")
			r.rules[tokens[0]] = rune(tokens[1][0])
		}
	}

	return nil
}
