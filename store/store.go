package store

import "github.com/asdine/brazier"

// GetBucketOrCreate returns an existing bucket or creates it if it doesn't exist.
func GetBucketOrCreate(r brazier.Registry, name string) (brazier.Bucket, error) {
	bucket, err := r.Bucket(name)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		err = r.Create(name)
		if err != nil {
			return nil, err
		}
		bucket, err = r.Bucket(name)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}
