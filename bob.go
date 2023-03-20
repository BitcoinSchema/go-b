package b

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/bitcoinschema/go-bpu"
)

// NewFromTape will create a new AIP object from a bob.Tape
// Using the FromTape() alone will prevent validation (data is needed via SetData to enable)
func NewFromTape(tape bpu.Tape) (b *B, err error) {
	b = new(B)
	err = b.FromTape(tape)
	return
}

// NewFromTapes will create a new B object from a []bob.Tape
// Using the FromTapes() alone will prevent validation (data is needed via SetData to enable)
func NewFromTapes(tapes []bpu.Tape) (b *B, err error) {
	// Loop tapes -> cells (only supporting 1 sig right now)
	for _, t := range tapes {
		for _, cell := range t.Cell {
			if cell.S != nil && *cell.S == Prefix {
				b = new(B)
				err = b.FromTape(t)
				// b.SetDataFromTapes(tapes)
				return
			}
		}
	}
	err = errors.New("no b tape found")
	return
}

// todo: SetDataFromTapes()

// FromTape takes a BOB Tape and returns a B data structure
func (b *B) FromTape(tape bpu.Tape) (err error) {
	if len(tape.Cell) < 3 { // B only requires 3 elements at minimum
		err = fmt.Errorf("invalid B tx Only %d pushdatas", len(tape.Cell))
		return
	}

	// Loop to find start of B
	var startIndex int
	for i, cell := range tape.Cell {
		if cell.S != nil && *cell.S == Prefix {
			startIndex = i
			break
		}
	}

	// Media type is after data
	b.MediaType = *tape.Cell[startIndex+2].S

	// Optional Encoding is after media
	if len(tape.Cell) > startIndex+3 && *tape.Cell[startIndex+3].S != "" {
		b.Encoding = *tape.Cell[startIndex+3].S
	} else {
		// default encoding is binary
		b.Encoding = string(EncodingBinary)
	}

	switch EncodingType(strings.ToLower(b.Encoding)) {
	case EncodingGzip:
		fallthrough
	case EncodingBinary:
		// Decode base64 data
		if tape.Cell[startIndex+1].B != nil {
			bStr := *tape.Cell[startIndex+1].B
			if b.Data.Bytes, err = base64.StdEncoding.DecodeString(bStr); err != nil {
				return
			}
		}
	case EncodingUtf8:
		fallthrough
	case EncodingUtf8Alt:
		if tape.Cell[startIndex+1].S != nil {
			b.Data.UTF8 = *tape.Cell[startIndex+1].S
		} else {
			b.Data.UTF8 = ""
		}
	}

	// Filename is optional and last
	if len(tape.Cell) > startIndex+4 && tape.Cell[startIndex+4].S != nil {
		b.Filename = *tape.Cell[startIndex+4].S
	}

	return
}
