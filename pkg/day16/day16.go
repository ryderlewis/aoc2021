package day16

import (
	"bufio"
	"github.com/ryderlewis/aoc2021/pkg/challenge"
	"io"
	"math"
	"strconv"
)

const (
	SUM         = 0
	PRODUCT     = 1
	MINIMUM     = 2
	MAXIMUM     = 3
	LITERAL     = 4
	GREATERTHAN = 5
	LESSTHAN    = 6
	EQUAL       = 7
)

type Packet struct {
	version int
	typeId  int

	// for typeId == 4
	literal int // for typeId 4

	// for typeId != 4
	lengthTypeId       int // 0 == sub-packet bit length, 1 == number of sub-packets
	subPacketBitLength int
	subPacketCount     int

	// raw data
	bits       []int
	subPackets []*Packet
}

type Runner struct {
	packet *Packet
}

var _ challenge.DailyChallenge = &Runner{}

func (r *Runner) Challenge1(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	val := r.sumVersions(r.packet)

	return strconv.Itoa(val), nil
}

func (r *Runner) Challenge2(input io.Reader) (string, error) {
	if err := r.readInput(input); err != nil {
		return "", err
	}

	val := r.doTheMath(r.packet)

	return strconv.Itoa(val), nil
}

func (r *Runner) sumVersions(packet *Packet) int {
	sum := packet.version

	for _, sub := range packet.subPackets {
		sum += r.sumVersions(sub)
	}

	return sum
}

func (r *Runner) doTheMath(packet *Packet) int {
	switch packet.typeId {
	case SUM:
		s := 0
		for _, p := range packet.subPackets {
			s += r.doTheMath(p)
		}
		return s
	case PRODUCT:
		v := 1
		for _, p := range packet.subPackets {
			v *= r.doTheMath(p)
		}
		return v
	case MINIMUM:
		v := math.MaxInt
		for _, p := range packet.subPackets {
			m := r.doTheMath(p)
			if m < v {
				v = m
			}
		}
		return v
	case MAXIMUM:
		v := math.MinInt
		for _, p := range packet.subPackets {
			m := r.doTheMath(p)
			if m > v {
				v = m
			}
		}
		return v
	case GREATERTHAN:
		v := 0
		if r.doTheMath(packet.subPackets[0]) > r.doTheMath(packet.subPackets[1]) {
			v = 1
		}
		return v
	case LESSTHAN:
		v := 0
		if r.doTheMath(packet.subPackets[0]) < r.doTheMath(packet.subPackets[1]) {
			v = 1
		}
		return v
	case EQUAL:
		v := 0
		if r.doTheMath(packet.subPackets[0]) == r.doTheMath(packet.subPackets[1]) {
			v = 1
		}
		return v
	default:
		return packet.literal
	}
}

func (r *Runner) readInput(input io.Reader) error {
	scanner := bufio.NewScanner(input)
	scanner.Scan()
	line := scanner.Text()
	bits := make([]int, len(line)*4)
	for nibble := 0; nibble < len(line); nibble++ {
		val := int(line[nibble] - '0')
		if line[nibble] >= 'A' {
			val = 10 + int(line[nibble]-'A')
		}
		bits[nibble*4] = val >> 3 & 1
		bits[nibble*4+1] = val >> 2 & 1
		bits[nibble*4+2] = val >> 1 & 1
		bits[nibble*4+3] = val & 1
	}

	r.packet, _ = r.packetFromBits(bits)

	return nil
}

func (r *Runner) packetFromBits(bits []int) (*Packet, int) {
	packet := &Packet{
		bits:       bits,
		version:    bits[0]<<2 | bits[1]<<1 | bits[2],
		typeId:     bits[3]<<2 | bits[4]<<1 | bits[5],
		subPackets: make([]*Packet, 0),
	}
	consumed := 6

	if packet.typeId == 4 {
		var c int
		packet.literal, c = r.intFromBits(bits[6:], true)
		consumed += c
	} else {
		packet.lengthTypeId = bits[6]
		consumed++

		if packet.lengthTypeId == 0 {
			// next 15 bits are total length in bits of sub-packets
			packet.subPacketBitLength, _ = r.intFromBits(bits[7:22], false)
			consumed += 15
			consumed += packet.subPacketBitLength

			subPacketBits := bits[22 : 22+packet.subPacketBitLength]
			for len(subPacketBits) > 0 {
				subPacket, bitCount := r.packetFromBits(subPacketBits)
				packet.subPackets = append(packet.subPackets, subPacket)
				subPacketBits = subPacketBits[bitCount:]
			}
		} else {
			// next 11 bits are sub-packet count
			packet.subPacketCount, _ = r.intFromBits(bits[7:18], false)
			consumed += 11

			subPacketBits := bits[18:]
			for i := 0; i < packet.subPacketCount; i++ {
				subPacket, bitCount := r.packetFromBits(subPacketBits)
				consumed += bitCount
				packet.subPackets = append(packet.subPackets, subPacket)
				subPacketBits = subPacketBits[bitCount:]
			}
		}
	}

	return packet, consumed
}

func (r *Runner) intFromBits(bits []int, variable bool) (int, int) {
	ret := 0
	consumed := 0

	if variable {
		cont := 1
		for cont != 0 {
			cont = bits[0]
			for _, b := range bits[1:5] {
				ret <<= 1
				ret |= b & 1
			}
			bits = bits[5:]
			consumed += 5
		}
	} else {
		for _, b := range bits {
			ret <<= 1
			ret |= b & 1
		}
		consumed = len(bits)
	}

	return ret, consumed
}
