package boltdb

// Bucket is a BoltDB implementation a bucket
type Bucket struct {
	ID       int
	PublicID string `storm:"unique"`
}

// Close the session of the bucket
func (b *Bucket) Close() error {
	return nil
}
