package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/storm"
)

const name = "boltdb"

// NewStore returns a BoltDB store
func NewStore(node storm.Node) *Store {
	return &Store{
		node: node,
	}
}

// Store is a BoltDB store
type Store struct {
	node storm.Node
}

// Name of the store
func (s *Store) Name() string {
	return name
}

// Create a bucket and return its informations
func (s *Store) Create(id string) error {
	return s.node.From(id).Init(&item{})
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	return NewBucket(s.node.From(id)), nil
}
