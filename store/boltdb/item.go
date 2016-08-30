package boltdb

import (
	"encoding/json"
	"time"
)

type item struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      json.RawMessage
	MimeType  string
}
