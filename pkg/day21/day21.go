package day21

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Player struct {
	pos int
	score int
}

type DeterministicDie struct {
	rolls int
	nextVal int
}

func (d *DeterministicDie) Roll(times int) int {
	val := 0
	for i := 0; i < times; i++ {
		d.rolls++
		val += d.nextVal
		d.nextVal++
		if d.nextVal > 100 {
			d.nextVal = 1
		}
	}
	return val
}

type Runner struct {
	p1, p2 *Player
	board [10]int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	die := &DeterministicDie{
		nextVal: 1,
	}

	for {
		r.p1.pos += die.Roll(3)
		r.p1.pos %= len(r.board)
		r.p1.score += r.board[r.p1.pos]
		if r.p1.score >= 1000 {
			break
		}

		r.p2.pos += die.Roll(3)
		r.p2.pos %= len(r.board)
		r.p2.score += r.board[r.p2.pos]
		if r.p2.score >= 1000 {
			break
		}
	}

	losingScore := r.p1.score
	if losingScore >= 1000 {
		losingScore = r.p2.score
	}

	return strconv.Itoa(die.rolls * losingScore), nil
}

type Universe struct {
	positions [2]int
	scores [2]int
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// universes is a count of possible universes
	universes := make(map[Universe]uint64)

	// and the starting universe is a single universe where both players have a score of 0
	universes[Universe{
		positions: [2]int{r.p1.pos, r.p2.pos},
	}] = 1
	currentPlayer := 0
	var wins [2]uint64

	// dirac die can have the following rolls:
	rolls := map[int]int{
		3: 1,
		4: 3,
		5: 6,
		6: 7,
		7: 6,
		8: 3,
		9: 1,
	}

	for len(universes) > 0 {
		nextUniverses := make(map[Universe]uint64)

		for u, c := range universes {
			for roll, count := range rolls {
				nextU := u
				nextU.positions[currentPlayer] += roll
				nextU.positions[currentPlayer] %= len(r.board)
				nextU.scores[currentPlayer] += r.board[nextU.positions[currentPlayer]]

				if nextU.scores[currentPlayer] >= 21 {
					wins[currentPlayer] += c * uint64(count)
				} else {
					nextUniverses[nextU] += c * uint64(count)
				}
			}
		}

		universes = nextUniverses
		currentPlayer++
		currentPlayer %= 2
	}

	bestScore := wins[0]
	if wins[1] > bestScore {
		bestScore = wins[1]
	}

	return strconv.FormatUint(bestScore, 10), nil
}

func (r *Runner) readInput(input io.Reader) error {
	r.p1 = &Player{}
	r.p2 = &Player{}

	scanner := bufio.NewScanner(input)
	scanner.Scan()
	p1Fields := strings.Fields(scanner.Text())
	r.p1.pos, _ = strconv.Atoi(p1Fields[len(p1Fields)-1])
	r.p1.pos--

	scanner.Scan()
	p2Fields := strings.Fields(scanner.Text())
	r.p2.pos, _ = strconv.Atoi(p2Fields[len(p2Fields)-1])
	r.p2.pos--

	for i := 0; i < 10; i++ {
		r.board[i] = i+1
	}

	return nil
}
