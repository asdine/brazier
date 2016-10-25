package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// NewStore returns a BoltDB store
func NewStore(path string) (*Store, error) {
	var err error

	db, err := storm.Open(
		path,
		storm.AutoIncrement(),
		storm.Codec(protobuf.Codec),
		storm.BoltOptions(0644, &bolt.Options{
			Timeout: time.Duration(50) * time.Millisecond,
		}),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open database")
	}

	return &Store{
		DB: db,
	}, nil
}

// Store is a BoltDB store
type Store struct {
	DB *storm.DB
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(path ...string) (brazier.Bucket, error) {
	return NewBucket(s, s.DB.From(path...), path...), nil
}

// Close BoltDB connection
func (s *Store) Close() error {
	return s.DB.Close()
}
