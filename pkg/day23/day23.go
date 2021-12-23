package day23

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
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
	DIRT = Amphipod('#')
	A = Amphipod('A')
	B = Amphipod('B')
	C = Amphipod('C')
	D = Amphipod('D')
)

type State [27]Amphipod

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
	cache map[State]int
	maxRecur int
}

// leastEnergy tries a naive recursive implementation, given the current state,
// find the minimum number of moves to get to a final state
func (r *Runner) leastEnergy(s *State, recur int) int {
	if recur > r.maxRecur {
		r.maxRecur = recur
	}
	/*
	fmt.Printf("Calculating leastEnergy, recur=%d:\n", recur)
	r.print(s)
	fmt.Println()
	 */

	if r.isFinal(s) {
		return 0
	}
	if val, ok := r.cache[*s]; ok {
		return val
	}

	var newState State
	newState = *s // makes a copy of s
	bestEnergy := -1

	for i, _ := range s {
		// find legal moves for amphipod a, at position i
		for _, move := range r.validMoves(s, i) {
			// swap values
			newState[move.from], newState[move.to] = newState[move.to], newState[move.from]
			/*
			fmt.Printf("Valid move, recur=%d:\n", recur)
			r.print(newState)
			fmt.Println()
			 */

			x := r.leastEnergy(&newState, recur+1)
			if x >= 0 {
				energy := move.energy + x
				if bestEnergy < 0 || energy < bestEnergy {
					bestEnergy = energy
				}
			}

			// swap back
			newState[move.from], newState[move.to] = newState[move.to], newState[move.from]
		}
	}

	r.cache[*s] = bestEnergy
	return bestEnergy
}

func (r *Runner) validMoves(s *State, i int) []*Move {
	if s[i] == EMPTY || s[i] == DIRT {
		return nil
	}

	a := s[i]
	data := r.amphipodData[s[i]]

	var moves []*Move

	if 0 <= i && i <= 10 {
		// amphipod is in the hall. only legal moves are to go to its final destination
		var dest, destDepth int
		for depth, d := range data.destIndexes {
			switch s[d] {
			case a:
				// this is fine
			case DIRT:
				// this is fine
			case EMPTY:
				dest = d
				destDepth = depth + 1
			default:
				// there is a conflict, cannot move to this destination now
				return nil
			}
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

		steps := hallIndexes[1] - hallIndexes[0] + destDepth
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
		// amphipod is in a room. first, see which room it's in.
		currRoomHallIndex := r.roomHallIndexes[i]
		destRoomHallIndex := data.hallIndex
		alreadyHome := currRoomHallIndex == destRoomHallIndex

		// if amphipod is in the right room, and there are no conflicts deeper in the room, then peace out.
		if alreadyHome {
			// see if there are any other amphipod types deeper in this room. If so, we need to move to let them out
			conflict := false
			for j := i+1; r.roomHallIndexes[j] == currRoomHallIndex; j++ {
				if s[j] != a && s[j] != DIRT {
					conflict = true
					break
				}
			}
			if !conflict {
				return nil // no conflicting amphipods deeper in the room. this one can stay where it's at
			}
		}

		// amphipod should leave the room. make sure it isn't blocked
		sourceDepth := 1
		for j := i-1; r.roomHallIndexes[j] == currRoomHallIndex; j-- {
			if s[j] != EMPTY {
				return nil // amphipod is blocked
			}
			sourceDepth++
		}

		// amphipod is in a room, but it needs to move out for now
		// first of all, see if the amphipod can go straight to its final destination. this is
		// the ideal case
		homeIsGood := true
		var homeDest, destDepth int
		for depth, d := range data.destIndexes {
			switch s[d] {
			case EMPTY:
				homeDest = d
				destDepth = depth + 1
			case DIRT:
				// this is fine
			case a:
				// this is fine
			default:
				// there is a conflict, cannot move to this destination now
				homeIsGood = false
			}
		}

		if homeIsGood {
			// see if there is a clear path to home. If not, then home is not good
			hallIndexes := []int{currRoomHallIndex, destRoomHallIndex}
			sort.Ints(hallIndexes)

			for j := hallIndexes[0]; j <= hallIndexes[1]; j++ {
				if s[j] != EMPTY {
					homeIsGood = false
					break
				}
			}

			if homeIsGood {
				// number of steps is all steps in the hall, plus destDepth steps into the room
				steps := hallIndexes[1] - hallIndexes[0] + sourceDepth + destDepth

				// fmt.Printf("HOME IS GOOD %d => %d!\n", i, homeDest)
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
		for j := currRoomHallIndex-1; j >= 0 && s[j] == EMPTY; j-- {
			if j == 2 || j == 4 || j == 6 || j == 8 {
				// can't stop outside a room
				continue
			}

			steps := currRoomHallIndex - j + sourceDepth
			// fmt.Printf("APPENDING MOVE LEFT %d(%d) => %d\n", i, currRoomHallIndex, j)
			moves = append(moves, &Move{
				amphipod: a,
				from: i,
				to: j,
				distance: steps,
				energy: data.energy * steps,
			})
		}

		for j := currRoomHallIndex+1; j <= 10 && s[j] == EMPTY; j++ {
			if j == 2 || j == 4 || j == 6 || j == 8 {
				// can't stop outside a room
				continue
			}

			steps := j - currRoomHallIndex + sourceDepth
			// fmt.Printf("APPENDING MOVE RIGHT %d(%d) => %d\n", i, currRoomHallIndex, j)
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
	for a, d := range r.amphipodData {
		for _, i := range d.destIndexes {
			if a != s[i] && s[i] != DIRT {
				return false
			}
		}
	}

	return true
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.print(&r.startingState)

	return strconv.Itoa(r.leastEnergy(&r.startingState, 0)), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.startingState[r.amphipodData[A].destIndexes[3]] = r.startingState[r.amphipodData[A].destIndexes[1]]
	r.startingState[r.amphipodData[A].destIndexes[1]] = D
	r.startingState[r.amphipodData[A].destIndexes[2]] = D

	r.startingState[r.amphipodData[B].destIndexes[3]] = r.startingState[r.amphipodData[B].destIndexes[1]]
	r.startingState[r.amphipodData[B].destIndexes[1]] = C
	r.startingState[r.amphipodData[B].destIndexes[2]] = B

	r.startingState[r.amphipodData[C].destIndexes[3]] = r.startingState[r.amphipodData[C].destIndexes[1]]
	r.startingState[r.amphipodData[C].destIndexes[1]] = B
	r.startingState[r.amphipodData[C].destIndexes[2]] = A

	r.startingState[r.amphipodData[D].destIndexes[3]] = r.startingState[r.amphipodData[D].destIndexes[1]]
	r.startingState[r.amphipodData[D].destIndexes[1]] = A
	r.startingState[r.amphipodData[D].destIndexes[2]] = C

	r.print(&r.startingState)

	return strconv.Itoa(r.leastEnergy(&r.startingState, 0)), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.initialize()

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	scanner.Scan()

	scanner.Scan()
	line := scanner.Text()
	r.startingState[11] = Amphipod(line[3])
	r.startingState[15] = Amphipod(line[5])
	r.startingState[19] = Amphipod(line[7])
	r.startingState[23] = Amphipod(line[9])

	scanner.Scan()
	line = scanner.Text()
	r.startingState[12] = Amphipod(line[3])
	r.startingState[16] = Amphipod(line[5])
	r.startingState[20] = Amphipod(line[7])
	r.startingState[24] = Amphipod(line[9])

	r.startingState[13] = DIRT
	r.startingState[14] = DIRT
	r.startingState[17] = DIRT
	r.startingState[18] = DIRT
	r.startingState[21] = DIRT
	r.startingState[22] = DIRT
	r.startingState[25] = DIRT
	r.startingState[26] = DIRT

	return nil
}

func (r *Runner) initialize() {
	r.cache = make(map[State]int)

	r.amphipodData = map[Amphipod]*AmphipodData{
		A: {
			hallIndex: 2,
			destIndexes: []int{11, 12, 13, 14},
			energy: 1,
		},
		B: {
			hallIndex: 4,
			destIndexes: []int{15, 16, 17, 18},
			energy: 10,
		},
		C: {
			hallIndex: 6,
			destIndexes: []int{19, 20, 21, 22},
			energy: 100,
		},
		D: {
			hallIndex: 8,
			destIndexes: []int{23, 24, 25, 26},
			energy: 1000,
		},
	}

	r.roomHallIndexes = make(map[int]int)
	for _, data := range r.amphipodData {
		for _, i := range data.destIndexes {
			r.roomHallIndexes[i] = data.hallIndex
		}
	}
}

func (r *Runner) print(s *State) {
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

	for i := 0; i < len(r.amphipodData[A].destIndexes); i++ {
		if i == 0 {
			fmt.Print("###")
		} else {
			fmt.Print("  #")
		}

		for _, a := range []Amphipod{A, B, C, D} {
			data := r.amphipodData[a]
			fmt.Printf("%s#", s[data.destIndexes[i]])
		}

		if i == 0 {
			fmt.Println("##")
		} else {
			fmt.Println()
		}
	}

	fmt.Println("  #########")
}

var _ challenge.DailyChallenge = &Runner{}
