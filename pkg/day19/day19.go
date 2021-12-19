package day19

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

type Coordinate struct {
	x, y, z int
}

type ScannerBeacon struct {
	scanner *Scanner
	sCoord  *Coordinate // coordinate from scanner's point of reference
	gCoord  *Coordinate // coordinate from global point of reference (scanner 0)
}

type Scanner struct {
	id          int
	orientation []int
	gCoord      *Coordinate
	beacons     []*ScannerBeacon
	distances   map[*ScannerBeacon][]float64
}

// the distance function only works for distances that share
// the same scanner orientation (relative direction)
func (c *Coordinate) distance(d *Coordinate) float64 {
	dx, dy, dz := d.x-c.x, d.y-c.y, d.z-c.z
	return math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
}

// matchingCoordinates finds coordinates between two scanners that are the same,
// determined by finding at least 11 overlapping distances to other coordinates within
// the same scanner
func (r *Runner) matchingCoordinates(s1, s2 *Scanner) map[*ScannerBeacon]*ScannerBeacon {
	same := make(map[*ScannerBeacon]*ScannerBeacon)

	for c1, f1 := range s1.distances {
		for c2, f2 := range s2.distances {
			// find overlapping distance counts between c1 and c2
			overlapCount := 0
			for i1, i2 := 0, 0; i1 < len(f1) && i2 < len(f2); {
				if f1[i1] == f2[i2] {
					overlapCount++
					i1++
					i2++
				} else if f1[i1] < f2[i2] {
					i1++
				} else if f2[i2] < f1[i1] {
					i2++
				}
			}
			if overlapCount >= 11 {
				same[c1] = c2
			}
		}
	}

	if len(same) < 12 {
		return nil
	}

	for c1, c2 := range same {
		c2.gCoord = c1.gCoord
	}

	// somewhat ridiculous calculation, try to figure out the actual orientation of
	// scanner 2 based on newly-discovered mapping of relative to global coordinates
	s2.discoverOrientation()
	s2.updateGlobalPositions()

	return same
}

func abs(i int) int {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func (s *Scanner) discoverOrientation() {
	if s.orientation != nil || s.gCoord != nil {
		return
	}

	// find an orientation that produces a consistent global scanner location
	// based on beacons with global coordinates known
	for _, orientation := range orientations {
		var gCoord *Coordinate
		consistent := true

		for _, sb := range s.beacons {
			if sb.gCoord == nil {
				continue
			}

			candidateCoord := Coordinate{}
			dists := []int{0, -sb.sCoord.y, -sb.sCoord.z, -sb.sCoord.x}
			xDist := dists[abs(orientation[0])]
			if orientation[0] < 0 {
				xDist *= -1
			}
			yDist := dists[abs(orientation[1])]
			if orientation[1] < 0 {
				yDist *= -1
			}
			zDist := dists[abs(orientation[2])]
			if orientation[2] < 0 {
				zDist *= -1
			}
			candidateCoord.x = sb.gCoord.x + xDist
			candidateCoord.y = sb.gCoord.y + yDist
			candidateCoord.z = sb.gCoord.z + zDist

			if gCoord == nil {
				gCoord = &candidateCoord
			} else if *gCoord != candidateCoord {
				consistent = false
				break
			}
		}

		if consistent {
			s.orientation = orientation
			s.gCoord = gCoord
			// fmt.Printf("%d: %#v\n", s.id, s.gCoord)
			return
		}
	}

	// should not reach
	panic("Could not find orientation :(")
}

func (s *Scanner) updateGlobalPositions() {
	if s.orientation == nil || s.gCoord == nil {
		// should never happen
		panic("Cannot update global positions :(")
	}

	for _, sb := range s.beacons {
		if sb.gCoord != nil {
			continue
		}

		dists := []int{0, sb.sCoord.y, sb.sCoord.z, sb.sCoord.x}
		xDist := dists[abs(s.orientation[0])]
		if s.orientation[0] < 0 {
			xDist *= -1
		}
		yDist := dists[abs(s.orientation[1])]
		if s.orientation[1] < 0 {
			yDist *= -1
		}
		zDist := dists[abs(s.orientation[2])]
		if s.orientation[2] < 0 {
			zDist *= -1
		}

		sb.gCoord = &Coordinate{
			x: s.gCoord.x + xDist,
			y: s.gCoord.y + yDist,
			z: s.gCoord.z + zDist,
		}
	}
}

func (s *Scanner) calcDistances() {
	s.distances = make(map[*ScannerBeacon][]float64)

	for i, b1 := range s.beacons {
		f := make([]float64, 0, len(s.beacons)-1)

		for j, b2 := range s.beacons {
			if i == j {
				continue
			}
			f = append(f, b1.sCoord.distance(b2.sCoord))
		}

		sort.Float64s(f)
		s.distances[b1] = f
	}
}

type Runner struct {
	scanners     []*Scanner
	lastGlobalId int
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	toCheck := make([]*Scanner, 1)
	toCheck[0] = r.scanners[0]
	checked := make(map[*Scanner]bool)

	for len(toCheck) > 0 {
		s1 := toCheck[0]
		toCheck = toCheck[1:]
		if _, alreadyChecked := checked[s1]; alreadyChecked {
			continue
		}
		checked[s1] = true

		for _, s2 := range r.scanners {
			if _, alreadyChecked := checked[s2]; alreadyChecked {
				continue
			}

			match := r.matchingCoordinates(s1, s2)
			if match == nil {
				continue
			}
			toCheck = append(toCheck, s2)
		}
	}

	// get all the global coordinates of all the beacons
	beacons := make(map[Coordinate]bool)
	for _, s := range r.scanners {
		for _, sb := range s.beacons {
			beacons[*sb.gCoord] = true
		}
	}

	return strconv.Itoa(len(beacons)), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	toCheck := make([]*Scanner, 1)
	toCheck[0] = r.scanners[0]
	checked := make(map[*Scanner]bool)

	for len(toCheck) > 0 {
		s1 := toCheck[0]
		toCheck = toCheck[1:]
		if _, alreadyChecked := checked[s1]; alreadyChecked {
			continue
		}
		checked[s1] = true

		for _, s2 := range r.scanners {
			if _, alreadyChecked := checked[s2]; alreadyChecked {
				continue
			}

			match := r.matchingCoordinates(s1, s2)
			if match == nil {
				continue
			}
			toCheck = append(toCheck, s2)
		}
	}

	maxManhattan := 0
	for _, s1 := range r.scanners {
		for _, s2 := range r.scanners {
			manhattan := abs(s1.gCoord.x - s2.gCoord.x) + abs(s1.gCoord.y - s2.gCoord.y) + abs(s1.gCoord.z - s2.gCoord.z)
			if manhattan > maxManhattan {
				maxManhattan = manhattan
			}
		}
	}

	return strconv.Itoa(maxManhattan), nil
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)

	var s *Scanner
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "---") {
			tokens := strings.Fields(line)
			id, _ := strconv.Atoi(tokens[2])
			s = &Scanner{
				id: id,
			}
			if id == 0 {
				s.gCoord = &Coordinate{0, 0, 0}
			}
			r.scanners = append(r.scanners, s)
		} else if len(line) > 0 {
			tokens := strings.Split(line, ",")
			x, _ := strconv.Atoi(tokens[0])
			y, _ := strconv.Atoi(tokens[1])
			z, _ := strconv.Atoi(tokens[2])
			sb := &ScannerBeacon{
				scanner: s,
				sCoord: &Coordinate{
					x: x,
					y: y,
					z: z,
				},
			}
			if s.id == 0 {
				s.orientation = orientations[0]
				sb.gCoord = sb.sCoord
			}

			s.beacons = append(s.beacons, sb)
		}
	}

	for _, s = range r.scanners {
		s.calcDistances()
	}

	return nil
}

// var orientations
/*
x => x, y => y, z => z    // 3,1,2
x => x, y => z, z => -y   // 3,2,-1
x => x, y => -y, z => -z  // 3,-1,-2
x => x, y => -z, z => y   // 3,-2,1

x => -x, y => y, z => -z  // -3,1,-2
x => -x, y => z, z => y   // -3,2,1
x => -x, y => -y, z => z  // -3,-1,2
x => -x, y => -z, z => -y // -3,-2,-1

x => y, y => -x, z => z   // 1,-3,2
x => y, y => z, z => x    // 1,2,3
x => y, y => x, z => -z   // 1,3,-2
x => y, y => -z, z => -x  // 1,-2,-3

x => -y, y => x, z => z   // -1,3,2
x => -y, y => -z, z => x  // -1,-2,3
x => -y, y => -x, z => -z // -1,-3,-2
x => -y, y => z, z => -x  // -1,2,-3

x => z, y => y, z => -x   // 2,1,-3
x => z, y => x, z => y    // 2,3,1
x => z, y => -y, z => x   // 2,-1,3
x => z, y => -x, z => -y  // 2,-3,-1

x => -z, y => y, z => x   // -2,1,3
x => -z, y => x, z => -y  // -2,3,-1
x => -z, y => -y, z => -x // -2,-1,-3
x => -z, y => -x, z => y  // -2,-3,1
*/
var orientations = [24][]int{
	{3, 1, 2},
	{3, 2, -1},
	{3, -1, -2},
	{3, -2, 1},
	{-3, 1, -2},
	{-3, 2, 1},
	{-3, -1, 2},
	{-3, -2, -1},
	{1, -3, 2},
	{1, 2, 3},
	{1, 3, -2},
	{1, -2, -3},
	{-1, 3, 2},
	{-1, -2, 3},
	{-1, -3, -2},
	{-1, 2, -3},
	{2, 1, -3},
	{2, 3, 1},
	{2, -1, 3},
	{2, -3, -1},
	{-2, 1, 3},
	{-2, 3, -1},
	{-2, -1, -3},
	{-2, -3, 1},
}
