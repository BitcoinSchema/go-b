package b

import (
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/bitcoinschema/go-bob"
)

// Prefix is the Bitcom prefix used by B
const Prefix = "19HxigV4QyBv3tHpQVcUEQyq1pzZVdoAut"

// Data is the content portion of the B data
type Data struct {
	Bytes []byte `json:"data,omitempty"`
	UTF8  string `json:"UTF8,omitempty"`
}

// B is B protocol
type B struct {
	Data      Data
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
		fallthrough
	case "binary":
		// Decode base64 data
		data, err := base64.StdEncoding.DecodeString(tape.Cell[1].B)
		if err != nil {
			log.Println("Failed to decode b64 signature", err)
			return err
		}
		b.Data.Bytes = data
	case "utf8":
		fallthrough
	case "utf-8":
		b.Data.UTF8 = tape.Cell[1].S
	}

	if len(tape.Cell) > 4 && len(tape.Cell[4].S) != 0 {
		b.Filename = tape.Cell[4].S
	}

	return nil
}

// BitFsURL is a helper to create a bitfs url to fetch the content over http
func BitFsURL(txid string, outIndex int, scriptChunk int) string {
	return fmt.Sprintf("https://x.bitfs.network/%s.out.%d.%d", txid, outIndex, scriptChunk)
}

// DataURI returns a b64 encoded image that can be set directly. Ex: <img src="b64data" />
func (b *B) DataURI() string {
	// encode raw bytes to b64
	s := base64.StdEncoding.EncodeToString(b.Data.Bytes)
	return fmt.Sprintf("data:%s;base64,%s", b.Encoding, s)
}
