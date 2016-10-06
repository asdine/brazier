package brazier

// An Item is a key value pair saved in a bucket.
type Item struct {
	Key  string
	Data []byte
}

// A Bucket manages a collection of items.
type Bucket interface {
	// Save a key value pair. It returns the created item.
	Save(key string, data []byte) (*Item, error)
	// Get an item from the bucket.
	Get(key string) (*Item, error)
	// Delete an item from the bucket.
	Delete(key string) error
	// Get the paginated list of items. perPage can be set to -1 to fetch all the items.
	Page(page int, perPage int) ([]Item, error)
	// Close the bucket. Can be used to close sessions if required.
	Close() error
}

// A Store manages the backend of specific buckets.
type Store interface {
	// Get a Bucket instance of this Store
	Bucket(name string) (Bucket, error)
	// Close the Store connexion.
	Close() error
}

// BucketConfig contains a bucket's configuration.
type BucketConfig struct {
	// Name of the bucket.
	Name string
}

// A Registry manages the buckets, their configuration and their associated Store.
type Registry interface {
	// Create a bucket and registers it to the Registry.
	Create(name string) error
	// Fetch a bucket directly from the associated Store.
	Bucket(name string) (Bucket, error)
	// Get a bucket's configuration.
	BucketConfig(name string) (*BucketConfig, error)
	// List all the buckets.
	List() ([]string, error)
	// Close the registry database.
	Close() error
}
