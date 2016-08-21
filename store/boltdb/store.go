package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/storm"
	"github.com/dchest/uniuri"
)

const name = "boltdb"

// NewStore returns a BoltDB store
func NewStore(db *storm.DB) *Store {
	storm.AutoIncrement()(db)

	return &Store{
		db: db,
	}
}

// Store is a BoltDB store
type Store struct {
	db *storm.DB
}

// Name of the store
func (s *Store) Name() string {
	return name
}

// Create a bucket and return its informations
func (s *Store) Create(id string) (*brazier.BucketInfo, error) {
	if id == "" {
		id = uniuri.NewLen(10)
	}

	return &brazier.BucketInfo{
		ID:        id,
		Store:     s.Name(),
		CreatedAt: time.Now(),
	}, nil
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (*Bucket, error) {
	return NewBucket(s.db.From(id)), nil
}
