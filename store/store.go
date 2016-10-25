package store

import "github.com/asdine/brazier"

// GetBucketOrCreate returns an existing bucket or creates it if it doesn't exist.
func GetBucketOrCreate(r brazier.Registry, path ...string) (brazier.Bucket, error) {
	bucket, err := r.Bucket(path...)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		err = r.Create(path...)
		if err != nil {
			return nil, err
		}
		bucket, err = r.Bucket(path...)
		if err != nil {
			return nil, err
		}
	}

	return bucket, nil
}
