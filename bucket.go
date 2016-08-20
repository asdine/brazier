package brazier

// A Bucket manages a collection of items.
type Bucket interface {
	Close() error
}

// BucketInfo holds bucket informations
type BucketInfo struct {
	ID    string
	Store string
}

// A Store manages the backend of specific buckets
type Store interface {
	Name() string
	Create(id string) (*BucketInfo, error)
	Bucket(id string) (Bucket, error)
}

// A Registrar registers bucket informations
type Registrar interface {
	Register(*BucketInfo) error
	Bucket(id string) (*BucketInfo, error)
}
