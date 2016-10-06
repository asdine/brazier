package brazier

// An Item is what is saved in a bucket. It contains user informations
type Item struct {
	Key  string
	Data []byte
}

// A Bucket manages a collection of items.
type Bucket interface {
	Save(key string, data []byte) (*Item, error)
	Get(key string) (*Item, error)
	Delete(key string) error
	Page(page int, perPage int) ([]Item, error)
	Close() error
}

// A Store manages the backend of specific buckets
type Store interface {
	Bucket(name string) (Bucket, error)
	Close() error
}

// BucketInfo contains a bucket informations
type BucketInfo struct {
	Name string
}

// A Registry holds bucket informations
type Registry interface {
	Create(name string) error
	Bucket(name string) (Bucket, error)
	BucketInfo(name string) (*BucketInfo, error)
	List() ([]string, error)
	Close() error
}
