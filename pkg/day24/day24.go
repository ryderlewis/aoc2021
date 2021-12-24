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

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

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

	candidateZVals := map[int]bool{
		0: true,
	}
	for i := len(instructionSets)-1; i >= 0; i-- {
		instructionSet := instructionSets[i]
		fmt.Printf("instructionSet: %v\n", instructionSet.instructions)

		for z := 0; z <= 1_000_000; z++ {
			for inp := 9; inp >= 1; inp-- {
				if vals, err := r.run(instructionSet.instructions, inp, [3]int{0, 0, z}); err == nil {
					if _, ok := candidateZVals[vals[2]]; ok {
						instructionSet.targetZVals[z] = true
					}
				}
			}
		}

		fmt.Printf("instruction set %d: target zvals: %d\n", i, len(instructionSet.targetZVals))
		candidateZVals = instructionSet.targetZVals
	}

	// now, run through instruction sets one at a time, caching possible outcomes
	return r.findHighest(instructionSets, [3]int{}), nil
}

func (r *Runner) findHighest(instructionSets []*InstructionSet, vars [3]int) string {
	inputsRemaining := len(instructionSets)

	for inp := 9; inp >= 1; inp-- {
		if vals, err := r.run(instructionSets[0].instructions, inp, vars); err == nil {
			if inputsRemaining == 1 { // last input
				if vals[2] == 0 {
					return strconv.Itoa(inp)
				} else {
					return ""
				}
			} else {
				if _, ok := instructionSets[1].targetZVals[vals[2]]; !ok {
					continue
				}

				recur := r.findHighest(instructionSets[1:], vals)
				if recur != "" {
					return strconv.Itoa(inp) + recur
				}
			}
		}
	}

	return "" // no good match
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

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {


		/*
		inp w      (W, X, Y, Z)
		mul x 0    (W, 0, Y, Z)
		add x z    (W, Z, Y, Z)
		mod x 26   (W, Z%26, Y, Z)
		div z 26   (W, Z%26, Y, Z/26)
		add x -2   (W, Z%26-2, Y, Z/26)
		eql x w    (W, (Z%26-2)==W, Y, Z/26)
		eql x 0    (W, (Z%26-2)!=W, Y, Z/26)

		mul y 0    (W, (Z%26-2)!=W, 0, Z/26)
		add y 25   (W, (Z%26-2)!=W, 25, Z/26)
		mul y x    (W, (Z%26-2)!=W, 25*x, Z/26)
		add y 1    (W, x, 25*x+1, Z/26)
		mul z y    (W, x, 25*x+1, Z/26*x)
		mul y 0    (W, x, 0, Z/26*x)
		add y w    (W, x, W, Z/26*x)
		add y 1    (W, x, W+1, Z/26*x)
		mul y x    (W, x, x*(W+1), Z/26*x)
		add z y    (W, x, x*(W+1), Z/26*x + x*(W+1))

		Z/26*x + x*(W+1) == 0

		x * (Z/26 + W + 1) == 0

		 */


		return "", err
	}

	return strconv.Itoa(0), nil
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
