package json

import (
	"encoding/json"
	"io"
	"unicode"
	"unicode/utf8"
)

// IsValid checks if the data is valid JSON
func IsValid(data []byte) bool {
	var i json.RawMessage

	return json.Unmarshal(data, &i) == nil
}

// IsValidReader checks if the data is valid JSON
func IsValidReader(r io.Reader) (bool, []byte) {
	var i json.RawMessage

	dec := json.NewDecoder(r)
	err := dec.Decode(&i)

	return err == nil || err == io.EOF, i
}

// Clean removes unnecessary white space from a JSON value
func Clean(data []byte) []byte {
	to := make([]byte, len(data))

	var inString bool
	start := 0
	i := 0
	l := len(data)
	for start < l {
		wid := 1
		r := rune(data[start])
		if r >= utf8.RuneSelf {
			r, wid = utf8.DecodeRune(data[start:])
		}
		start += wid

		if r == '"' {
			inString = !inString
		} else if unicode.IsSpace(r) {
			if !inString {
				continue
			}
		}
		i += utf8.EncodeRune(to[i:], r)
	}

	return to[:i]
}
