package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/storm"
)

const name = "boltdb"

// NewStore returns a BoltDB store
func NewStore(db *storm.DB) *Store {
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
	b := Bucket{
		ID: id,
	}

	err := s.db.Save(&b)

	return &brazier.BucketInfo{
		ID:    b.ID,
		Store: s.Name(),
	}, err
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	var b Bucket

	err := s.db.One("ID", id, &b)
	return &b, err
}
