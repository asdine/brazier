package boltdb

import (
	"github.com/asdine/brazier"
	"github.com/asdine/storm"
)

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
func (r *Registrar) Register(info *brazier.BucketInfo) error {
	return r.db.Save(info)
}

// Bucket returns the bucket informations associated with the given id
func (r *Registrar) Bucket(id string) (*brazier.BucketInfo, error) {
	var info brazier.BucketInfo

	err := r.db.One("ID", id, &info)
	return &info, err
}
