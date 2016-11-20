package json

import (
	"bytes"
	"encoding/json"

	"github.com/asdine/brazier"
)

// MarshalList marshals a list of items
func MarshalList(items []brazier.Item) ([]byte, error) {
	return json.Marshal(marshalList(items))
}

// MarshalListPretty marshals a list of items
func MarshalListPretty(items []brazier.Item) ([]byte, error) {
	return json.MarshalIndent(marshalList(items), "", "  ")
}

func marshalList(items []brazier.Item) []map[string]interface{} {
	list := make([]map[string]interface{}, len(items))

	for i := range items {
		k := items[i].Key
		list[i] = map[string]interface{}{
			"key": k,
		}
		if items[i].Children != nil {
			list[i]["value"] = marshalList(items[i].Children)

			continue
		}

		v := json.RawMessage(items[i].Data)
		list[i]["value"] = &v
	}

	return list
}

// ToValidJSON converts data to a valid JSON payload
func ToValidJSON(data []byte) []byte {
	if IsValid(data) {
		return Clean(data)
	}

	count := bytes.Count(data, []byte(`"`))
	out := make([]byte, len(data)+count+2)

	out[0] = '"'
	j := 1
	for _, b := range data {
		if b == '"' {
			out[j] = '\\'
			j++
		}
		out[j] = b
		j++
	}
	out[j] = '"'
	return out
}
