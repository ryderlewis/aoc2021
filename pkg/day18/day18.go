package day18

import (
	"bufio"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"strconv"
	"strings"
)

type Pair struct {
	parent *Pair

	leftVal *int
	leftPair *Pair

	rightVal *int
	rightPair *Pair
}

type Runner struct {
	inputPairs []*Pair
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	var p *Pair
	for _, q := range r.inputPairs {
		if p == nil {
			p = q
		} else {
			p = p.combine(q)
		}
		// fmt.Printf("%s => ", p)
		p.reduce()
		// fmt.Printf("%s\n", p)
	}

	return strconv.Itoa(p.magnitude()), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	maxMagnitude := 0

	// find results of all pairwise additions
	for i, p1 := range r.inputPairs {
		for j, p2 := range r.inputPairs {
			if i == j {
				continue
			}
			p1c := p1.copy()
			p2c := p2.copy()
			c := p1c.combine(p2c)

			c.reduce()
			m := c.magnitude()
			if m > maxMagnitude {
				maxMagnitude = m
				// fmt.Printf("%d: %s + %s = %s\n", m, p1, p2, c)
			}
		}
	}

	return strconv.Itoa(maxMagnitude), nil
}

func (p *Pair) combine(q *Pair) *Pair {
	n := &Pair{
		leftPair: p,
		rightPair: q,
	}

	p.parent = n
	q.parent = n

	return n
}

func (p *Pair) copy() *Pair {
	return parsePair([]rune(p.String()))
}

func (p *Pair) reduce() {
	for {
		if p.explode() {
			continue
		}
		if p.split() {
			continue
		}
		break
	}
}

func (p *Pair) magnitude() int {
	mag := 0

	if p.leftVal != nil {
		mag += *p.leftVal * 3
	} else {
		mag += p.leftPair.magnitude() * 3
	}

	if p.rightVal != nil {
		mag += *p.rightVal * 2
	} else {
		mag += p.rightPair.magnitude() * 2
	}

	return mag
}

func (p *Pair) isLeft() bool {
	return p.parent != nil && p.parent.leftPair == p
}

func (p *Pair) isRight() bool {
	return p.parent != nil && p.parent.rightPair == p
}

func (p *Pair) isRoot() bool {
	return p.parent == nil
}

func (p *Pair) nestCount() int {
	if p.isRoot() {
		return 0
	}

	return p.parent.nestCount() + 1
}

type TraverseInfo struct {
	nodes []*TraverseNode
}

type TraverseNode struct {
	val  *int
	pair *Pair
	nestCount int
	isLeft, isRight bool
}

func (p *Pair) traverse(info *TraverseInfo) {
	if p.leftVal != nil {
		info.nodes = append(info.nodes, &TraverseNode{
			val: p.leftVal,
			pair: p,
			nestCount: p.nestCount(),
			isLeft: true,
		})
	} else {
		p.leftPair.traverse(info)
	}

	if p.rightVal != nil {
		info.nodes = append(info.nodes, &TraverseNode{
			val: p.rightVal,
			pair: p,
			nestCount: p.nestCount(),
			isRight: true,
		})
	} else {
		p.rightPair.traverse(info)
	}
}

func (p *Pair) explode() bool {
	// find first pair in in-order traversal that is more than 4 levels nested
	traverseInfo := &TraverseInfo{}
	p.traverse(traverseInfo)

	for i, t := range traverseInfo.nodes {
		if t.nestCount == 4 {
			// this one needs to be exploded!
			if i > 0 {
				*traverseInfo.nodes[i-1].val += *t.val
			}
			if i + 2 < len(traverseInfo.nodes) {
				*traverseInfo.nodes[i+2].val += *traverseInfo.nodes[i+1].val
			}
			zero := 0
			parent := t.pair.parent
			if t.pair.isLeft() {
				parent.leftPair = nil
				parent.leftVal = &zero
			} else {
				parent.rightPair = nil
				parent.rightVal = &zero
			}
			return true
		}
	}

	return false
}

func (p *Pair) split() bool {
	// find first pair in in-order traversal that is more than 9
	traverseInfo := &TraverseInfo{}
	p.traverse(traverseInfo)

	for _, t := range traverseInfo.nodes {
		if *t.val > 9 {
			lval := *t.val/2
			rval := (*t.val+1)/2
			newPair := &Pair{
				parent: t.pair,
				leftVal: &lval,
				rightVal: &rval,
			}
			if t.isLeft {
				t.pair.leftVal = nil
				t.pair.leftPair = newPair
			} else {
				t.pair.rightVal = nil
				t.pair.rightPair = newPair
			}
			return true
		}
	}

	return false
}

func (r *Runner) readInput(input io.Reader) error {
	r.inputPairs = make([]*Pair, 0)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		r.inputPairs = append(r.inputPairs, parsePair([]rune(scanner.Text())))
	}

	return nil
}

func (p *Pair) String() string {
	var b strings.Builder
	b.WriteRune('[')

	if p.leftVal != nil {
		_, _ = fmt.Fprintf(&b, "%d,", *p.leftVal)
	} else {
		_, _ = fmt.Fprintf(&b, "%s,", p.leftPair)
	}

	if p.rightVal != nil {
		_, _ = fmt.Fprintf(&b, "%d]", *p.rightVal)
	} else {
		_, _ = fmt.Fprintf(&b, "%s]", p.rightPair)
	}

	return b.String()
}

func parsePair(runes []rune) *Pair {
	// [[1,2],3]
	p := &Pair{}

	// parse left side
	var commaPos int
	if runes[1] == '[' {
		lpos := 1
		rpos := endBracketPos(runes, lpos)
		p.leftPair = parsePair(runes[lpos:rpos+1])
		p.leftPair.parent = p
		commaPos = rpos+1
	} else {
		leftVal := int(runes[1]-'0')
		p.leftVal = &leftVal
		commaPos = 2
	}

	// parse right side
	if runes[commaPos+1] == '[' {
		lpos := commaPos+1
		rpos := endBracketPos(runes, lpos)
		p.rightPair = parsePair(runes[lpos:rpos+1])
		p.rightPair.parent = p
	} else {
		rightVal := int(runes[commaPos+1]-'0')
		p.rightVal = &rightVal
	}

	return p
}

func endBracketPos(runes []rune, startBracketPos int) int {
	brackets := 1
	var rpos int
	for i := startBracketPos+1; i < len(runes); i++ {
		if runes[i] == '[' {
			brackets++
		} else if runes[i] == ']' {
			brackets--
			if brackets == 0 {
				rpos = i
				break
			}
		}
	}

	return rpos
}
