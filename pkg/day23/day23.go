package day23

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"sort"
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
	roomHallIndexes map[int]int
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
		hallIndexes := []int{i, r.roomHallIndexes[dest]}
		sort.Ints(hallIndexes)
		for j := hallIndexes[0]; j <= hallIndexes[1]; j++ {
			if i == j {
				continue
			}
			if s[j] != EMPTY {
				return nil // path is blocked
			}
		}

		steps := hallIndexes[1] - hallIndexes[0]
		if dest >= 15 {
			steps += 2
		} else {
			steps += 1
		}
		return []*Move{
			{
				amphipod: a,
				from: i,
				to: dest,
				distance: steps,
				energy: data.energy * steps,
			},
		}
	} else {
		// amphipod is in a room. first, see if it's already home.
		if i == data.destIndexes[1] {
			// amphipod is deep into its own room. it's home.
			return nil
		} else if i == data.destIndexes[0] && s[data.destIndexes[1]] == a {
			// amphipod is next to another amphipod, they are both home
			return nil
		}

		// next see if the amphipod is blocked from moving.
		if i >= 15 && s[i-4] != EMPTY {
			// blocked.
			return nil
		}

		// amphipod is in a room, but it needs to move out for now
		roomHallIndex := r.roomHallIndexes[i]
		destHallIndex := data.hallIndex

		// first of all, see if the amphipod can go straight to its final destination. this is
		// the ideal case
		homeIsGood := true
		var homeDest int
		if s[data.destIndexes[0]] == EMPTY {
			switch s[data.destIndexes[1]] {
			case a:
				homeDest = data.destIndexes[0]
			case EMPTY:
				homeDest = data.destIndexes[1]
			default:
				homeIsGood = false
			}
		}

		if homeIsGood {
			// see if there is a clear path to home. If not, then home is not good
			hallIndexes := []int{roomHallIndex, destHallIndex}
			sort.Ints(hallIndexes)

			for j := hallIndexes[0]; j <= hallIndexes[1]; j++ {
				if s[j] != EMPTY {
					homeIsGood = false
					break
				}
			}

			if homeIsGood {
				// number of steps is all steps in the hall, plus either one or two steps
				// to go into the room
				steps := hallIndexes[1] - hallIndexes[0] + 1 + 1
				if homeDest >= 15 {
					steps += 1
				}
				if i >= 15 {
					steps += 1 // starting point is deep, need an extra step to get out
				}
				return []*Move{
					{
						amphipod: a,
						from: i,
						to: homeDest,
						distance: steps,
						energy: data.energy * steps,
					},
				}
			}
		}

		// home is not good. So now we need to figure out where in the hall this amphipod can go.
		for j := roomHallIndex-1; j >= 0 && s[j] == EMPTY; j-- {
			if j == 2 || j == 4 || j == 6 || j == 8 {
				// can't stop outside a room
				continue
			}

			steps := roomHallIndex - j + 1
			if i >= 15 {
				steps += 1 // starting point is deep
			}
			moves = append(moves, &Move{
				amphipod: a,
				from: i,
				to: j,
				distance: steps,
				energy: data.energy * steps,
			})
		}

		for j := roomHallIndex+1; j <= 10 && s[j] == EMPTY; j++ {
			if j == 2 || j == 4 || j == 6 || j == 8 {
				// can't stop outside a room
				continue
			}

			steps := j - roomHallIndex + 1
			if i >= 15 {
				steps += 1 // starting point is deep
			}
			moves = append(moves, &Move{
				amphipod: a,
				from: i,
				to: j,
				distance: steps,
				energy: data.energy * steps,
			})
		}
	}

	return moves
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
		D: {
			hallIndex: 8,
			destIndexes: []int{14, 18},
			energy: 1000,
		},
	}

	r.roomHallIndexes = map[int]int{
		11: 2,
		12: 4,
		13: 6,
		14: 8,
		15: 2,
		16: 4,
		17: 6,
		18: 8,
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
