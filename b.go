// Package b is a library for working with B transactions (Bitcoin OP_RETURN protocol) in Go
//
// If you have any suggestions or comments, please feel free to open an issue on
// this GitHub repository!
//
// By BitcoinSchema Organization (https://bitcoinschema.org)
package b

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// Prefix is the Bitcom prefix used by B
const Prefix = "19HxigV4QyBv3tHpQVcUEQyq1pzZVdoAut"

// B is B protocol
type B struct {
	Data      []byte
	MediaType string
	Encoding  string
	Filename  string
}

// MarshalJSON custom json marshaler
func (b *B) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"content-type": b.MediaType,
		"encoding":     b.Encoding,
		"filename":     b.Filename,
	}
	switch EncodingType(strings.ToLower(b.Encoding)) {
	case EncodingUtf8, EncodingUtf8Alt:
		data["content"] = string(b.Data)
	case EncodingBinary, EncodingGzip:
		fallthrough
	default:
		data["content"] = base64.StdEncoding.EncodeToString(b.Data)
	}
	return json.Marshal(data)
}

// UnmarshalJSON custom json unmarshaler
func (b *B) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	b.MediaType = raw["content-type"].(string)
	b.Encoding = raw["encoding"].(string)
	b.Filename = raw["filename"].(string)
	switch EncodingType(strings.ToLower(b.Encoding)) {
	case EncodingUtf8, EncodingUtf8Alt:
		b.Data = []byte(raw["content"].(string))
	case EncodingBinary, EncodingGzip:
		fallthrough
	default:
		var err error
		b.Data, err = base64.StdEncoding.DecodeString(raw["content"].(string))
		if err != nil {
			return err
		}
	}
	return nil
}

// EncodingType is an enum for the possible types of data encoding
type EncodingType string

// Various encoding types
const (
	EncodingBinary  EncodingType = "binary"
	EncodingGzip    EncodingType = "gzip"
	EncodingUtf8    EncodingType = "utf8"
	EncodingUtf8Alt EncodingType = "utf-8"
)

// DataURI returns a b64 encoded image that can be set directly. Ex: <img src="b64data" />
func (b *B) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", b.Encoding, base64.StdEncoding.EncodeToString(b.Data))
}

// BitFsURL is a helper to create a bitfs url to fetch the content over HTTP
func BitFsURL(txID string, outIndex, scriptChunk int) string {
	return fmt.Sprintf("https://x.bitfs.network/%s.out.%d.%d", txID, outIndex, scriptChunk)
}
