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

// buildInstructionSets returns the lowest and highest possible solutions to the problem
func (r *Runner) buildInstructionSets() []*InstructionSet {
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
	maxTargetZ := 0
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

		for z := 0; z <= 10_000; z++ {
			for inp := 1; inp <= 9; inp++ {
				if vals, err := r.run(instructionSet.instructions, inp, [3]int{0, 0, z}); err == nil {
					if _, ok := instructionSet.targetZVals[vals[2]]; ok {
						nextTargetZVals[z] = true
						if z > maxTargetZ {
							maxTargetZ = z
						}
					}
				}
			}
		}
	}

	fmt.Printf("max target z: %d\n", maxTargetZ)
	return instructionSets
}

func (r *Runner) findFirstSolution(instructionSets []*InstructionSet, vars[3]int, highest bool) string {
	inputsRemaining := len(instructionSets)-1
	solution := ""

	var firstInput, lastInput, delta int
	if highest {
		firstInput, lastInput, delta = 9, 0, -1
	} else {
		firstInput, lastInput, delta = 1, 10, 1
	}

	for inp := firstInput; inp != lastInput; inp += delta {
		if vals, err := r.run(instructionSets[0].instructions, inp, vars); err == nil {
			if _, ok := instructionSets[0].targetZVals[vals[2]]; !ok {
				continue
			}

			if inputsRemaining == 0 {
				// last input
				solution = strconv.Itoa(inp)
				break
			} else {
				// recursively solve
				if recur := r.findFirstSolution(instructionSets[1:], vals, highest); recur != "" {
					solution = strconv.Itoa(inp) + recur
					break
				}
			}
		}
	}

	return solution
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	instructionSets := r.buildInstructionSets()
	return r.findFirstSolution(instructionSets, [3]int{}, true), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	// updated part 1 to return both solutions
	if err := r.readInput(input); err != nil {
		return "", err
	}

	instructionSets := r.buildInstructionSets()
	return r.findFirstSolution(instructionSets, [3]int{}, false), nil
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
