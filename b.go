package b

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/rohenaz/go-bob"
)

// Prefix is the Bitcom prefix used by B
const Prefix = "19HxigV4QyBv3tHpQVcUEQyq1pzZVdoAut"

// B is B protocol
type B struct {
	Data      []byte `json:"data"`
	MediaType string `json:"mediaType"`
	Encoding  string `json:"encoding"`
	Filename  string `json:"filename,omitempty"`
}

// New creates a new B struct
func New() *B {
	return &B{}
}

// FromTape takes a BOB Tape and returns a B data structure
func (b *B) FromTape(tape bob.Tape) error {
	if len(tape.Cell) < 4 || tape.Cell[0].S != Prefix {
		return fmt.Errorf("Invalid B tx Only %d pushdatas", len(tape.Cell))
	}

	// Decode base64 data
	data, err := base64.StdEncoding.DecodeString(tape.Cell[1].B)
	if err != nil {
		log.Println("Failed to decode b64 signature", err)
		return err
	}
	b.Data = data // base 64 decode
	b.MediaType = tape.Cell[2].S
	b.Encoding = tape.Cell[3].S
	b.Filename = tape.Cell[4].S

	return nil
}
