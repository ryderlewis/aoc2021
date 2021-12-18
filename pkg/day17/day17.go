package day17

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"strconv"
	"strings"
)

type Runner struct {
	minX, maxX int
	minY, maxY int
}

type XVel struct {
	startVel int
	finalVal int
	minSteps int
	maxSteps int
}

type YVel struct {
	startVel   int
	highestVal int
	minSteps   int
	maxSteps   int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// get x velocities
	xvels := r.getXVels()

	// now figure out all the starting y velocities that will get us there.
	yvels := r.getYVels()

	style := 0
	for _, xvel := range xvels {
		for _, yvel := range yvels {
			if xvel.maxSteps < yvel.minSteps || xvel.minSteps > yvel.maxSteps {
				continue
			}
			if style < yvel.highestVal {
				style = yvel.highestVal
			}
		}
	}

	return strconv.Itoa(style), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// get x velocities
	xvels := r.getXVels()

	// now figure out all the starting y velocities that will get us there.
	yvels := r.getYVels()

	count := 0
	for _, xvel := range xvels {
		for _, yvel := range yvels {
			if xvel.maxSteps < yvel.minSteps || xvel.minSteps > yvel.maxSteps {
				continue
			}
			count++
		}
	}

	return strconv.Itoa(count), nil
}

func (r *Runner) getXVels() map[int]*XVel {
	xvels := make(map[int]*XVel)

	for v := 0; v <= r.maxX; v++ {
		finalVal := v * (v + 1) / 2
		if finalVal < r.minX {
			continue
		}

		minSteps := 0
		maxSteps := 0

		currX := 0
		for i := 0; i <= v; i++ {
			if currX >= r.minX && currX <= r.maxX {
				if minSteps == 0 {
					minSteps = i
				}
				maxSteps = i
			}

			currX += v - i
		}

		if maxSteps > 0 {
			if finalVal >= r.minX && finalVal <= r.maxX {
				maxSteps = math.MaxInt
			}

			xvels[v] = &XVel{
				startVel: v,
				finalVal: finalVal,
				minSteps: minSteps,
				maxSteps: maxSteps,
			}
		}
	}

	return xvels
}

func (r *Runner) getYVels() map[int]*YVel {
	// let's find all the possible steps that a particular x velocity can put us in the target range.
	yvels := make(map[int]*YVel)

	for v := r.minY; v <= -r.minY; v++ {
		highestVal := 0
		if v > 0 {
			highestVal = v * (v + 1) / 2
		}

		minSteps := 0
		maxSteps := 0

		currY := 0
		for i := 0; currY >= r.minY; i++ {
			if currY >= r.minY && currY <= r.maxY {
				if minSteps == 0 {
					minSteps = i
				}
				maxSteps = i
			}

			currY += v - i
		}

		if maxSteps > 0 {
			yvels[v] = &YVel{
				startVel: v,
				highestVal: highestVal,
				minSteps: minSteps,
				maxSteps: maxSteps,
			}
		}
	}

	return yvels
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	scanner.Scan()
	line := scanner.Text()

	ytokens := strings.Split(line, ", y=")
	yequals := strings.Split(ytokens[1], "..")
	xtokens := strings.Split(ytokens[0], ": x=")
	xequals := strings.Split(xtokens[1], "..")

	r.minX, _ = strconv.Atoi(xequals[0])
	r.maxX, _ = strconv.Atoi(xequals[1])
	r.minY, _ = strconv.Atoi(yequals[0])
	r.maxY, _ = strconv.Atoi(yequals[1])

	return nil
}
