package brazier

// A Bucket manages a collection of documents.
type Bucket interface {
	Close() error
}

// BucketInfo holds a bucket informations
type BucketInfo struct {
	ID      string
	StoreID string
}

// A Store manages backend specific Buckets
type Store interface {
	ID() string
	Create() (*BucketInfo, error)
	Bucket(id string) Bucket
}

// A Registrar registers bucket informations
type Registrar interface {
	Register(*BucketInfo) error
	Bucket(id string) BucketInfo
}
