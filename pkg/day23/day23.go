package day23

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
)

type Amphipod rune

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

type Room struct {
	hallIndex int
	amphipods [2]Amphipod
}

type Runner struct {
	hall [11]Amphipod
	rooms [4]Room
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	r.print()

	return strconv.Itoa(0), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(0), nil
}

func (r *Runner) print() {
/*
   #############
   #...........#
   ###B#C#B#D###
     #A#D#C#A#
     #########
*/
	fmt.Println("#############")
	fmt.Print("#")
	for _, a := range r.hall {
		fmt.Printf("%s", a)
	}
	fmt.Println("#")

	fmt.Print("###")
	for _, room := range r.rooms {
		fmt.Printf("%s#", room.amphipods[0])
	}
	fmt.Println("##")

	fmt.Print("  #")
	for _, room := range r.rooms {
		fmt.Printf("%s#", room.amphipods[1])
	}
	fmt.Println()
	fmt.Println("  #########")
}

func (r *Runner) readInput(input io.Reader) error {
	r.rooms[0].hallIndex = 2
	r.rooms[1].hallIndex = 4
	r.rooms[2].hallIndex = 6
	r.rooms[3].hallIndex = 8

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	scanner.Scan()

	scanner.Scan()
	line := scanner.Text()
	r.rooms[0].amphipods[0] = Amphipod(line[3])
	r.rooms[1].amphipods[0] = Amphipod(line[5])
	r.rooms[2].amphipods[0] = Amphipod(line[7])
	r.rooms[3].amphipods[0] = Amphipod(line[9])

	scanner.Scan()
	line = scanner.Text()
	r.rooms[0].amphipods[1] = Amphipod(line[3])
	r.rooms[1].amphipods[1] = Amphipod(line[5])
	r.rooms[2].amphipods[1] = Amphipod(line[7])
	r.rooms[3].amphipods[1] = Amphipod(line[9])

	return nil
}
