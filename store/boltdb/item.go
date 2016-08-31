package boltdb

import "time"

type item struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
	MimeType  string
}
