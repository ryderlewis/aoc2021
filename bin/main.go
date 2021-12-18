package main

import (
	"flag"
	"fmt"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"github.com/ryderlewis/aoc2021/pkg/day01"
	"github.com/ryderlewis/aoc2021/pkg/day02"
	"github.com/ryderlewis/aoc2021/pkg/day03"
	"github.com/ryderlewis/aoc2021/pkg/day04"
	"github.com/ryderlewis/aoc2021/pkg/day05"
	"github.com/ryderlewis/aoc2021/pkg/day06"
	"github.com/ryderlewis/aoc2021/pkg/day07"
	"github.com/ryderlewis/aoc2021/pkg/day08"
	"github.com/ryderlewis/aoc2021/pkg/day09"
	"github.com/ryderlewis/aoc2021/pkg/day10"
	"github.com/ryderlewis/aoc2021/pkg/day11"
	"github.com/ryderlewis/aoc2021/pkg/day12"
	"github.com/ryderlewis/aoc2021/pkg/day13"
	"github.com/ryderlewis/aoc2021/pkg/day14"
	"github.com/ryderlewis/aoc2021/pkg/day15"
	"github.com/ryderlewis/aoc2021/pkg/day16"
	"github.com/ryderlewis/aoc2021/pkg/day17"
	"github.com/ryderlewis/aoc2021/pkg/day18"
	"github.com/ryderlewis/aoc2021/pkg/day19"
	"github.com/ryderlewis/aoc2021/pkg/day20"
	"github.com/ryderlewis/aoc2021/pkg/day21"
	"github.com/ryderlewis/aoc2021/pkg/day22"
	"github.com/ryderlewis/aoc2021/pkg/day23"
	"github.com/ryderlewis/aoc2021/pkg/day24"
	"github.com/ryderlewis/aoc2021/pkg/day25"
	"io"
	"os"
)

func main() {
	day := flag.Int("day", 0, "Day number, 1 through 25")
	chnum := flag.Int("challenge", 0, "Challenge number, 1 or 2")
	fname := flag.String("filename", "", "File with input values")

	flag.Parse()

	if *day < 1 || *day > 25 || *chnum < 1 || *chnum > 2 {
		flag.Usage()
	}

	var input io.Reader
	if *fname == "" || *fname == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(*fname)
		if err != nil {
			fmt.Printf("Couldn't open %v: %v\n", *fname, err)
			flag.Usage()
		}

		input = f
		defer f.Close()
	}

	var dc challenge.DailyChallenge
	var answer string

	switch *day {
	case 1:
		dc = &day01.Runner{}
	case 2:
		dc = &day02.Runner{}
	case 3:
		dc = &day03.Runner{}
	case 4:
		dc = &day04.Runner{}
	case 5:
		dc = &day05.Runner{}
	case 6:
		dc = &day06.Runner{}
	case 7:
		dc = &day07.Runner{}
	case 8:
		dc = &day08.Runner{}
	case 9:
		dc = &day09.Runner{}
	case 10:
		dc = &day10.Runner{}
	case 11:
		dc = &day11.Runner{}
	case 12:
		dc = &day12.Runner{}
	case 13:
		dc = &day13.Runner{}
	case 14:
		dc = &day14.Runner{}
	case 15:
		dc = &day15.Runner{}
	case 16:
		dc = &day16.Runner{}
	case 17:
		dc = &day17.Runner{}
	case 18:
		dc = &day18.Runner{}
	case 19:
		dc = &day19.Runner{}
	case 20:
		dc = &day20.Runner{}
	case 21:
		dc = &day21.Runner{}
	case 22:
		dc = &day22.Runner{}
	case 23:
		dc = &day23.Runner{}
	case 24:
		dc = &day24.Runner{}
	case 25:
		dc = &day25.Runner{}
	}

	if dc == nil {
		fmt.Printf("Day not implemented: %d\n", *day)
		flag.Usage()
	}

	var err error
	if *chnum == 1 {
		answer, err = dc.Challenge1(input)
	} else {
		answer, err = dc.Challenge2(input)
	}

	if err == nil {
		fmt.Printf("Day %d, challenge %d: %s\n", *day, *chnum, answer)
	} else {
		fmt.Printf("Error: %v\n", err)
	}
}
