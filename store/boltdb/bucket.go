package boltdb

// Bucket is a BoltDB implementation a bucket
type Bucket struct {
	ID string
}

// Close the session of the bucket
func (b *Bucket) Close() error {
	return nil
}
