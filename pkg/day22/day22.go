package day22

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"strconv"
	"strings"
)

type Region struct {
	minX, maxX, minY, maxY, minZ, maxZ int
	on bool
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r Region) volume() int {
	return max(r.maxX-r.minX+1, 0) * max(r.maxY-r.minY+1, 0) * max(r.maxZ-r.minZ+1, 0)
}

// intersects tests to see if region r intersects region o
func (r Region) intersects(o Region) bool {
	if r.maxX < o.minX || r.minX > o.maxX {
		return false
	}
	if r.maxY < o.minY || r.minY > o.maxY {
		return false
	}
	if r.maxZ < o.minZ || r.minZ > o.maxZ {
		return false
	}
	return true
}

// spawn new regions from r that are non overlapping with o
func (r Region) nonOverlapping(o Region) []Region {
	ret := make([]Region, 0)

	if !r.intersects(o) {
		ret = append(ret, r)
		return ret
	}

	overlappingMinX := max(r.minX, o.minX)
	overlappingMaxX := min(r.maxX, o.maxX)
	overlappingMinY := max(r.minY, o.minY)
	overlappingMaxY := min(r.maxY, o.maxY)
	overlappingMinZ := max(r.minZ, o.minZ)
	overlappingMaxZ := min(r.maxZ, o.maxZ)

	// add up to six regions not including the regions that overlap with o
	if r.minX < overlappingMinX {
		ret = append(ret, Region{r.minX, overlappingMinX-1, r.minY, r.maxY, r.minZ, r.maxZ, r.on})
	}
	if r.maxX > overlappingMaxX {
		ret = append(ret, Region{overlappingMaxX+1, r.maxX, r.minY, r.maxY, r.minZ, r.maxZ, r.on})
	}
	if r.minY < overlappingMinY {
		ret = append(ret, Region{overlappingMinX, overlappingMaxX, r.minY, overlappingMinY-1, r.minZ, r.maxZ, r.on})
	}
	if r.maxY > overlappingMaxY {
		ret = append(ret, Region{overlappingMinX, overlappingMaxX, overlappingMaxY+1, r.maxY, r.minZ, r.maxZ, r.on})
	}
	if r.minZ < overlappingMinZ {
		ret = append(ret, Region{overlappingMinX, overlappingMaxX, overlappingMinY, overlappingMaxY, r.minZ, overlappingMinZ-1, r.on})
	}
	if r.maxZ > overlappingMaxZ {
		ret = append(ret, Region{overlappingMinX, overlappingMaxX, overlappingMinY, overlappingMaxY, overlappingMaxZ+1, r.maxZ, r.on})
	}

	return ret
}

type Runner struct {
	regions []Region
}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	// turn off anything outside of -50..50 in each dimension
	r.regions = append(r.regions, Region{math.MinInt, -51, math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, false})
	r.regions = append(r.regions, Region{51, math.MaxInt, math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, false})
	r.regions = append(r.regions, Region{math.MinInt, math.MaxInt, math.MinInt, -51, math.MinInt, math.MaxInt, false})
	r.regions = append(r.regions, Region{math.MinInt, math.MaxInt, 51, math.MaxInt, math.MinInt, math.MaxInt, false})
	r.regions = append(r.regions, Region{math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, math.MinInt, -51, false})
	r.regions = append(r.regions, Region{math.MinInt, math.MaxInt, math.MinInt, math.MaxInt, 51, math.MaxInt, false})

	return strconv.Itoa(r.reboot()), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	return strconv.Itoa(r.reboot()), nil
}

// reboot runs the reboot steps and returns a count of on cubes
func (r *Runner) reboot() int {
	on := make([]Region, 0)

	for _, reg := range r.regions {
		nextOn := make([]Region, 0)

		for _, reg2 := range on {
			for _, reg3 := range reg2.nonOverlapping(reg) {
				nextOn = append(nextOn, reg3)
			}
		}

		if reg.on {
			nextOn = append(nextOn, reg)
		}

		on = nextOn
	}

	onCount := 0
	for _, reg := range on {
		onCount += reg.volume()
	}

	return onCount
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		tokens := strings.Fields(scanner.Text())
		coords := strings.Split(tokens[1], ",")

		xeq := strings.Split(strings.Split(coords[0], "=")[1], "..")
		yeq := strings.Split(strings.Split(coords[1], "=")[1], "..")
		zeq := strings.Split(strings.Split(coords[2], "=")[1], "..")

		minX, _ := strconv.Atoi(xeq[0])
		maxX, _ := strconv.Atoi(xeq[1])
		minY, _ := strconv.Atoi(yeq[0])
		maxY, _ := strconv.Atoi(yeq[1])
		minZ, _ := strconv.Atoi(zeq[0])
		maxZ, _ := strconv.Atoi(zeq[1])

		r.regions = append(r.regions, Region{
			minX: minX,
			maxX: maxX,
			minY: minY,
			maxY: maxY,
			minZ: minZ,
			maxZ: maxZ,
			on: tokens[0] == "on",
		})
	}

	return nil
}

var _ challenge.DailyChallenge = &Runner{}