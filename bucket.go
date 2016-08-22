package brazier

import "time"

// An Item is what is saved in a bucket. It contains user informations
// and some metadata
type Item struct {
	ID        string
	CreatedAt time.Time
	Data      []byte
}

// A Bucket manages a collection of items.
type Bucket interface {
	Add(data []byte, name string) (*Item, error)
	Get(id string) (*Item, error)
	Delete(id string) error
	Close() error
}

// BucketInfo holds bucket informations
type BucketInfo struct {
	ID        string
	Stores    []string
	CreatedAt time.Time
}

// A Store manages the backend of specific buckets
type Store interface {
	Name() string
	Create(id string) error
	Bucket(id string) (Bucket, error)
}

// A Registrar registers bucket informations
type Registrar interface {
	Create(id string, s Store) (*BucketInfo, error)
	Bucket(id string) (*BucketInfo, error)
}
