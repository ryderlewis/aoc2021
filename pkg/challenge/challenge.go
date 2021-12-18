package challenge

import (
	"io"
)

type DailyChallenge interface {
	Challenge1(input io.Reader) (string, error)
    Challenge2(input io.Reader) (string, error)
}
