package b

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/rohenaz/go-bob"
)

// Prefix is the Bitcom prefix used by B
const Prefix = "19HxigV4QyBv3tHpQVcUEQyq1pzZVdoAut"

// B is B protocol
type B struct {
	Data      []byte `json:"data,omitempty"`
	UTF8      string `json:"UTF8,omitempty"`
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

	b.MediaType = tape.Cell[2].S
	b.Encoding = tape.Cell[3].S
	switch strings.ToLower(b.Encoding) {
	case "gzip":
	case "binary":
		// Decode base64 data
		data, err := base64.StdEncoding.DecodeString(tape.Cell[1].B)
		if err != nil {
			log.Println("Failed to decode b64 signature", err)
			return err
		}
		b.Data = data // base 64 decode
	case "utf8":
		fallthrough
	case "utf-8":
		b.UTF8 = tape.Cell[1].S
	}

	if len(tape.Cell[4].S) != 0 {
		b.Filename = tape.Cell[4].S
	}

	return nil
}

// BitFsURL is a helper to create a bitfs url to fetch the content over http
func BitFsURL(txid string, outIndex int, scriptChunk int) string {
	return fmt.Sprintf("%s.out.%d.%d", txid, outIndex, scriptChunk)
}
