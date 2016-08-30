package brazier

import "time"

// An Item is what is saved in a bucket. It contains user informations
// and some metadata
type Item struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      []byte
}

// A Bucket manages a collection of items.
type Bucket interface {
	Save(id string, data []byte) (*Item, error)
	Get(id string) (*Item, error)
	Delete(id string) error
	Close() error
}

// A Store manages the backend of specific buckets
type Store interface {
	Create(id string) error
	Bucket(id string) (Bucket, error)
}
