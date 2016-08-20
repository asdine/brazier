package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/pkg/errors"
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
	b := Bucket{
		PublicID: id,
	}

	err := s.db.Save(&b)
	if err != nil {
		if err == storm.ErrAlreadyExists {
			return nil, store.ErrAlreadyExists
		}

		return nil, errors.Wrap(err, "store create bucket failed")
	}

	return &brazier.BucketInfo{
		ID:    b.PublicID,
		Store: s.Name(),
	}, nil
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	var b Bucket

	err := s.db.One("PublicID", id, &b)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}

		return nil, errors.Wrap(err, "store get bucket failed")
	}

	return &b, nil
}
