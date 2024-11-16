package grids

import (
	"strings"
)

type ScanMode uint8

const (
	// Bit 1: i (x) direction scanning
	ScanModePositiveI ScanMode = 0 // Points scan in +i direction
	ScanModeNegativeI ScanMode = 1 // Points scan in -i direction

	// Bit 2: j (y) direction scanning
	ScanModeNegativeJ ScanMode = 0 // Points scan in -j direction
	ScanModePositiveJ ScanMode = 2 // Points scan in +j direction

	// Bit 3: Adjacent points
	ScanModeConsecutiveI ScanMode = 0 // Adjacent points in i direction are consecutive
	ScanModeConsecutiveJ ScanMode = 4 // Adjacent points in j direction are consecutive

	// Bit 4: Row direction
	ScanModeSameDirection ScanMode = 0 // All rows scan in same direction
	ScanModeOppositeRows  ScanMode = 8 // Adjacent rows scan in opposite direction

	// Bit 5: Odd row offset
	ScanModeNoOddOffset ScanMode = 0  // Points within odd rows not offset
	ScanModeOddOffset   ScanMode = 16 // Points within odd rows offset by Di/2

	// Bit 6: Even row offset
	ScanModeNoEvenOffset ScanMode = 0  // Points within even rows not offset
	ScanModeEvenOffset   ScanMode = 32 // Points within even rows offset by Di/2

	// Bit 7: J direction offset
	ScanModeNoJOffset ScanMode = 0  // Points not offset in j direction
	ScanModeJOffset   ScanMode = 64 // Points offset by Dj/2 in j direction

	// Bit 8: Row/Column point counts
	ScanModeRegularPoints ScanMode = 0   // Regular Ni x Nj points
	ScanModeOffsetPoints  ScanMode = 128 // Points may be reduced based on offsets
)

// Helper methods to check individual bits
func (s ScanMode) IsNegativeI() bool {
	return s&1 == 1
}

func (s ScanMode) IsPositiveJ() bool {
	return s&2 == 2
}

func (s ScanMode) IsConsecutiveJ() bool {
	return s&4 == 4
}

func (s ScanMode) HasOppositeRows() bool {
	return s&8 == 8
}

func (s ScanMode) HasOddOffset() bool {
	return s&16 == 16
}

func (s ScanMode) HasEvenOffset() bool {
	return s&32 == 32
}

func (s ScanMode) HasJOffset() bool {
	return s&64 == 64
}

func (s ScanMode) HasOffsetPoints() bool {
	return s&128 == 128
}

// String returns a human readable description of the scan mode
func (s ScanMode) String() string {
	var desc []string

	// i direction
	if s.IsNegativeI() {
		desc = append(desc, "-i scanning")
	} else {
		desc = append(desc, "+i scanning")
	}

	// j direction
	if s.IsPositiveJ() {
		desc = append(desc, "+j scanning")
	} else {
		desc = append(desc, "-j scanning")
	}

	// Consecutive points
	if s.IsConsecutiveJ() {
		desc = append(desc, "consecutive j")
	} else {
		desc = append(desc, "consecutive i")
	}

	return strings.Join(desc, ", ")
}
