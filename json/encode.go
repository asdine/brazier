package json

import (
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
