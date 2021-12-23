package day23

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"strconv"
)

type Amphipod byte

func (a Amphipod) String() string {
	if a == EMPTY {
		return "."
	} else {
		return string(a)
	}
}

const (
	EMPTY = Amphipod(0)
	A = Amphipod('A')
	B = Amphipod('B')
	C = Amphipod('C')
	D = Amphipod('D')
)

type State [19]Amphipod

type Move struct {
	amphipod Amphipod
	from, to, distance, energy int
}

type AmphipodData struct {
	hallIndex int
	destIndexes []int
	energy int
}

type Runner struct {
	startingState State
	amphipodData map[Amphipod]*AmphipodData
}

// leastEnergy tries a naive recursive implementation, given the current state,
// find the minimum number of moves to get to a final state
func (r *Runner) leastEnergy(s *State) int {
	if r.isFinal(s) {
		return 0
	}

	newState := *s // makes a copy of s
	bestEnergy := math.MaxInt
	for i, _ := range s {
		// find legal moves for amphipod a, at position i
		for _, move := range r.validMoves(s, i) {
			// swap values
			newState[move.from], newState[move.to] = newState[move.to], newState[move.from]

			energy := move.energy + r.leastEnergy(&newState)
			if energy < bestEnergy {
				bestEnergy = energy
			}

			// swap back
			newState[move.from], newState[move.to] = newState[move.to], newState[move.from]
		}
	}

	return bestEnergy
}

func (r *Runner) validMoves(s *State, i int) []*Move {
	if s[i] == EMPTY {
		return nil
	}

	a := s[i]
	data := r.amphipodData[s[i]]

	var hallBlocked [11]bool
	var moves []*Move

	if 0 <= i && i <= 10 {
		// amphipod is in the hall. only legal moves are to go to its final destination
		var dest int
		if s[data.destIndexes[1]] == EMPTY && s[data.destIndexes[0]] == EMPTY {
			dest = data.destIndexes[1]
		} else if s[data.destIndexes[1]] == a && s[data.destIndexes[0]] == EMPTY {
			dest = data.destIndexes[0]
		} else {
			// no legal move
			return nil
		}

		// see if there's a clear path from the current location to the destination index
		move := r.move(s, data, i, dest)
		if move != nil {
			moves = append(moves, move)
		}
	} else if 11 <= i && i <= 14 {
		// amphipod is in a room closest to the hall. see if it can move out of the room
		b := s[i+4] // other amphipod in the same room
		if a == b && data.destIndexes[0] == i {
			// amphipod is in its final destination
			return nil
		}

		// amphipod can move to let out its partner, or go to its final home.
		// Note that to do this, we need to see which way this amphipod should
		// move so that it doesn't block any other amphipods in the hall from
		// moving to their destination.
	}

	return moves
}

func (r *Runner) move(s *State, data *AmphipodData, from, to int) *Move {
	return nil
}

func (r *Runner) isFinal(s *State) bool {
	return s[11] == A && s[15] == A &&
		s[12] == B && s[16] == B &&
		s[13] == C && s[17] == C &&
		s[14] == D && s[18] == D
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(r.leastEnergy(&r.startingState)), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(0), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.initialize()

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	scanner.Scan()

	scanner.Scan()
	line := scanner.Text()
	r.startingState[11] = Amphipod(line[3])
	r.startingState[12] = Amphipod(line[5])
	r.startingState[13] = Amphipod(line[7])
	r.startingState[14] = Amphipod(line[9])

	scanner.Scan()
	line = scanner.Text()
	r.startingState[15] = Amphipod(line[3])
	r.startingState[16] = Amphipod(line[5])
	r.startingState[17] = Amphipod(line[7])
	r.startingState[18] = Amphipod(line[9])

	return nil
}

func (r *Runner) initialize() {
	r.amphipodData = map[Amphipod]*AmphipodData{
		A: {
			hallIndex: 2,
			destIndexes: []int{11, 15},
			energy: 1,
		},
		B: {
			hallIndex: 4,
			destIndexes: []int{12, 16},
			energy: 10,
		},
		C: {
			hallIndex: 6,
			destIndexes: []int{13, 17},
			energy: 100,
		},
		C: {
			hallIndex: 8,
			destIndexes: []int{14, 18},
			energy: 1000,
		},
	}
}

func (s State) print() {
	/*
	   #############
	   #...........#
	   ###B#C#B#D###
	     #A#D#C#A#
	     #########
	*/
	fmt.Println("#############")
	fmt.Print("#")
	for _, a := range s[:11] {
		fmt.Printf("%s", a)
	}
	fmt.Println("#")

	fmt.Print("###")
	for _, a := range s[11:15] {
		fmt.Printf("%s#", a)
	}
	fmt.Println("##")

	fmt.Print("  #")
	for _, a := range s[15:] {
		fmt.Printf("%s#", a)
	}
	fmt.Println()
	fmt.Println("  #########")
}

var _ challenge.DailyChallenge = &Runner{}
