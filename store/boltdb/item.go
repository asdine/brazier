package boltdb

import "time"

type item struct {
	ID        int
	PublicID  string `storm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
	MimeType  string
}
