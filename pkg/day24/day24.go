package day24

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Var byte

const (
	INP = "inp"
	ADD = "add"
	MUL = "mul"
	DIV = "div"
	MOD = "mod"
	EQL = "eql"
	W = Var('w')
	X = Var('x')
	Y = Var('y')
	Z = Var('z')
	LITERAL = Var(0)
)

var (
	DIVERR = errors.New("div by zero")
	MODERR = errors.New("mod by zero or neg")
)

type Instruction struct {
	operator string
	operand1 Var
	operand2 Var
	literal int
}

type InstructionSet struct {
	instructions []*Instruction
	targetZVals map[int]bool
}

type State struct {
	remaining, z int
}

type Runner struct {
	instructions []*Instruction
}

var _ challenge.DailyChallenge = &Runner{}

// solve returns the lowest and highest possible solutions to the problem
func (r *Runner) solve() (string, string) {
	// work instructions one input at a time
	instructionSets := make([]*InstructionSet, 0)
	startIndex := 0
	for i := 1; i < len(r.instructions); i++ {
		if r.instructions[i].operator == INP {
			instructionSets = append(instructionSets, &InstructionSet{
				instructions: r.instructions[startIndex:i],
				targetZVals: make(map[int]bool),
			})
			startIndex = i
		}
	}
	// assumes final instruction is not an input
	instructionSets = append(instructionSets, &InstructionSet{
		instructions: r.instructions[startIndex:],
		targetZVals: make(map[int]bool),
	})

	// work backwards. The final instruction set wants an output of z==0.
	// in order to do this, we want to know which input values of z will produce
	// the desired output This informs us of which possible output z values from each
	// step are acceptable towards building a solution.
	var nextTargetZVals map[int]bool
	for i := len(instructionSets)-1; i >= 0; i-- {
		instructionSet := instructionSets[i]

		if i == len(instructionSets)-1 {
			instructionSet.targetZVals = map[int]bool{
				0: true,
			}
		} else {
			instructionSet.targetZVals = nextTargetZVals
		}
		nextTargetZVals = make(map[int]bool)

		for z := 0; z <= 1_000_000; z++ {
			for inp := 1; inp <= 9; inp++ {
				if vals, err := r.run(instructionSet.instructions, inp, [3]int{0, 0, z}); err == nil {
					if _, ok := instructionSet.targetZVals[vals[2]]; ok {
						nextTargetZVals[z] = true
					}
				}
			}
		}
	}

	possibleSolutions := r.findSolutions(instructionSets, [3]int{})
	fmt.Printf("possible solutions: %d\n", len(possibleSolutions))

	return possibleSolutions[0], possibleSolutions[len(possibleSolutions)-1]
}

func (r *Runner) findSolutions(instructionSets []*InstructionSet, vars[3]int) []string {
	solutions := make([]string, 0)
	inputsRemaining := len(instructionSets)-1

	for inp := 1; inp <= 9; inp++ {
		if vals, err := r.run(instructionSets[0].instructions, inp, vars); err == nil {
			if _, ok := instructionSets[0].targetZVals[vals[2]]; !ok {
				continue
			}

			if inputsRemaining == 0 {
				// last input
				solutions = append(solutions, strconv.Itoa(inp))
			} else {
				// recursively solve
				for _, recur := range r.findSolutions(instructionSets[1:], vals) {
					solutions = append(solutions, strconv.Itoa(inp) + recur)
				}
			}
		}
	}

	return solutions
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	p1, p2 := r.solve()

	// now, run through instruction sets one at a time, caching possible outcomes
	return fmt.Sprintf("%s - %s", p1, p2), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	// updated part 1 to return both solutions
	return r.Challenge1(input)
}

func (r *Runner) run(instructions []*Instruction, input int, v [3]int) ([3]int, error) {
	vars := map[Var]int{
		W: 0,
		X: v[0],
		Y: v[1],
		Z: v[2],
	}

	for _, inst := range instructions {
		var val int
		if inst.operand2 == LITERAL {
			val = inst.literal
		} else {
			val = vars[inst.operand2]
		}

		switch inst.operator {
		case INP:
			vars[inst.operand1] = input
		case ADD:
			vars[inst.operand1] += val
		case MUL:
			vars[inst.operand1] *= val
		case DIV:
			if val == 0 {
				// fmt.Printf("DIVERR: %d / %d\n", vars[inst.operand1], val)
				return [3]int{}, DIVERR
			}
			vars[inst.operand1] /= val
		case MOD:
			if val <= 0 {
				// fmt.Printf("MODERR: %d %% %d\n", vars[inst.operand1], val)
				return [3]int{}, MODERR
			}
			if vars[inst.operand1] < 0 {
				// fmt.Printf("MODERR: %d %% %d\n", vars[inst.operand1], val)
				return [3]int{}, MODERR
			}
			vars[inst.operand1] %= val
		case EQL:
			if vars[inst.operand1] == val {
				vars[inst.operand1] = 1
			} else {
				vars[inst.operand1] = 0
			}
		default:
			panic(fmt.Sprintf("bad operator: %s", inst.operator))
		}
	}

	return [3]int{vars[X], vars[Y], vars[Z]}, nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		instruction := &Instruction{
			operator: tokens[0],
			operand1: Var(tokens[1][0]),
		}
		if len(tokens) > 2 {
			switch tokens[2] {
			case "w", "x", "y", "z":
				instruction.operand2 = Var(tokens[2][0])
			default:
				instruction.operand2 = LITERAL
				instruction.literal, _ = strconv.Atoi(tokens[2])
			}
		}
		r.instructions = append(r.instructions, instruction)
	}

	return nil
}
