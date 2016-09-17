package store

import "github.com/asdine/brazier"

// GetBucketOrCreate returns an existing bucket or creates it if it doesn't exist
func GetBucketOrCreate(s brazier.Store, name string) (brazier.Bucket, error) {
	bucket, err := s.Bucket(name)
	if err != nil {
		if err != ErrNotFound {
			return nil, err
		}
		err = s.Create(name)
		if err != nil {
			return nil, err
		}
		bucket, err = s.Bucket(name)
		if err != nil {
			return nil, err
		}
	}
	return bucket, nil
}
