package boltdb

import (
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

type bucketInfo struct {
	ID        int
	PublicID  string `storm:"unique"`
	Stores    []string
	CreatedAt time.Time
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

// Create a new bucket in the registrar
func (r *Registrar) Create(id string, s brazier.Store) (*brazier.BucketInfo, error) {
	if id == "" {
		id = uniuri.NewLen(10)
	}

	i := bucketInfo{
		PublicID:  id,
		Stores:    []string{s.Name()},
		CreatedAt: time.Now(),
	}

	err := r.db.Save(&i)
	if err != nil {
		if err == storm.ErrAlreadyExists {
			return nil, store.ErrAlreadyExists
		}
		return nil, errors.Wrap(err, "boltdb.registrar.Create failed saving bucket")
	}

	err = s.Create(id)
	if err != nil {
		return nil, errors.Wrap(err, "boltdb.registrar.Create failed creating bucket in the store")
	}

	return &brazier.BucketInfo{
		ID:        i.PublicID,
		Stores:    i.Stores,
		CreatedAt: i.CreatedAt,
	}, nil
}

// Bucket returns the bucket informations associated with the given id
func (r *Registrar) Bucket(id string) (*brazier.BucketInfo, error) {
	var i bucketInfo

	err := r.db.One("PublicID", id, &i)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "boltdb.registrar.Bucket failed getting bucket")
	}

	return &brazier.BucketInfo{
		ID:        i.PublicID,
		CreatedAt: i.CreatedAt,
		Stores:    i.Stores,
	}, nil
}
