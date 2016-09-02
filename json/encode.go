package json

import (
	"encoding/json"

	"github.com/asdine/brazier"
)

// MarshalList marshals a list of items
func MarshalList(items []brazier.Item) ([]byte, error) {
	list := make([]map[string]*json.RawMessage, len(items))

	for i := range items {
		d := json.RawMessage(items[i].Data)
		list[i] = map[string]*json.RawMessage{
			items[i].ID: &d,
		}
	}

	return json.Marshal(list)
}
