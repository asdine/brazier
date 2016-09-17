package json

import (
	"bytes"
	"encoding/json"

	"github.com/asdine/brazier"
)

// MarshalList marshals a list of items
func MarshalList(items []brazier.Item) ([]byte, error) {
	list := make([]map[string]interface{}, len(items))

	for i := range items {
		k := items[i].Key
		v := json.RawMessage(items[i].Data)
		list[i] = map[string]interface{}{
			"key":  k,
			"data": &v,
		}
	}

	return json.Marshal(list)
}

// ToValidJSON converts data to a valid JSON payload
func ToValidJSON(data []byte) []byte {
	if IsValid(data) {
		return Clean(data)
	}
	var buffer bytes.Buffer
	buffer.Grow(len(data) + 2)
	buffer.WriteByte('"')
	buffer.Write(data)
	buffer.WriteByte('"')
	return buffer.Bytes()
}
