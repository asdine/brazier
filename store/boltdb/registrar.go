package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/pkg/errors"
)

type info struct {
	ID       int
	BucketID string `storm:"unique"`
	Store    string
}

// NewRegistrar returns a Registrar
func NewRegistrar(db *storm.DB) *Registrar {
	return &Registrar{
		db: db,
	}
}

// A Registrar registers bucket informations
type Registrar struct {
	db *storm.DB
}

// Register a bucket in the registrar
func (r *Registrar) Register(bucketInfo *brazier.BucketInfo) error {
	i := info{
		BucketID: bucketInfo.ID,
		Store:    bucketInfo.Store,
	}

	err := r.db.Save(&i)
	if err != nil {
		if err == storm.ErrAlreadyExists {
			return store.ErrAlreadyExists
		}
		return errors.Wrap(err, "registrar register bucket failed")
	}
	return nil
}

// Bucket returns the bucket informations associated with the given id
func (r *Registrar) Bucket(id string) (*brazier.BucketInfo, error) {
	var i info

	err := r.db.One("BucketID", id, &i)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "registrar get bucket failed")
	}

	return &brazier.BucketInfo{
		ID:    i.BucketID,
		Store: i.Store,
	}, nil
}
