package day08

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Runner struct {
	inputs   [][]string
	outputs  [][]string
	mappings []map[rune]rune
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	count := 0
	for _, output := range r.outputs {
		for _, o := range output {
			switch len(o) {
			case 2:
				count++
			case 3:
				count++
			case 4:
				count++
			case 7:
				count++
			}
		}
	}
	return strconv.Itoa(count), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	digits := map[int]string{
		0: "abcefg",
		1: "cf",
		2: "acdeg",
		3: "acdfg",
		4: "bcdf",
		5: "abdfg",
		6: "abdefg",
		7: "acf",
		8: "abcdefg",
		9: "abcdfg",
	}
	revDigits := make(map[string]int)
	for i, s := range digits {
		revDigits[s] = i
	}

	sum := 0
	for inputIndex, puzzleInput := range r.inputs {
		// find a coherent mapping
		for _, m := range r.mappings {
			isGood := true
			// fmt.Printf("Mapping: %#v\n", m)

			for _, pi := range puzzleInput {
				mappedInput := SortString(strings.Map(func(r rune) rune {
					return m[r]
				}, pi))
				// fmt.Printf("Mapped: %s => %s\n", pi, mappedInput)

				if _, ok := revDigits[mappedInput]; !ok {
					isGood = false
					break
				}
			}

			if isGood {
				output := r.outputs[inputIndex]
				mult := 1000
				for _, o := range output {
					mappedOutput := SortString(strings.Map(func(r rune) rune {
						return m[r]
					}, o))
					sum += mult * revDigits[mappedOutput]
					mult /= 10
				}
				break
			}
		}
	}

	return strconv.Itoa(sum), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	r.inputs = make([][]string, 0)
	r.outputs = make([][]string, 0)

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " | ")
		r.inputs = append(r.inputs, strings.Fields(parts[0]))
		r.outputs = append(r.outputs, strings.Fields(parts[1]))
	}

	r.mappings = make([]map[rune]rune, 0)
	runes := "abcdefg"
	for _, p := range permutations(runes) {
		m := make(map[rune]rune)
		for i := 0; i < len(runes); i++ {
			m[rune(runes[i])] = rune(p[i])
		}
		r.mappings = append(r.mappings, m)
	}

	return nil
}

func join(ins []rune, c rune) (result []string) {
	for i := 0; i <= len(ins); i++ {
		result = append(result, string(ins[:i])+string(c)+string(ins[i:]))
	}
	return
}

func permutations(testStr string) []string {
	var n func(testStr []rune, p []string) []string
	n = func(testStr []rune, p []string) []string {
		if len(testStr) == 0 {
			return p
		} else {
			result := []string{}
			for _, e := range p {
				result = append(result, join([]rune(e), testStr[0])...)
			}
			return n(testStr[1:], result)
		}
	}

	output := []rune(testStr)
	return n(output[1:], []string{string(output[0])})
}

type sortRunes []rune

func (s sortRunes) Less(i, j int) bool {
	return s[i] < s[j]
}

func (s sortRunes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s sortRunes) Len() int {
	return len(s)
}

func SortString(s string) string {
	r := []rune(s)
	sort.Sort(sortRunes(r))
	return string(r)
}
