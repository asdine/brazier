package fs

import (
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

// newBucket returns a Bucket
func newBucket(path string) *bucket {
	return &bucket{
		path: path,
	}
}

// bucket is a file system implementation of a brazier bucket
type bucket struct {
	path string
}

// Add user data to the bucket. Returns an Iten
func (b *bucket) Add(data []byte, mimeType string, name string) (*brazier.Item, error) {
	if name == "" {
		exts, err := mime.ExtensionsByType(mimeType)
		if err != nil {
			return nil, errors.Wrap(err, "fs.bucket.Add failed to get file extension")
		}
		if len(exts) == 0 {
			return nil, errors.New("fs.bucket.Add unknown mime type")
		}
		name = uniuri.NewLen(10) + exts[0]
	}

	f, err := os.Create(filepath.Join(b.path, name))
	if err != nil {
		return nil, errors.Wrap(err, "fs.bucket.Add failed to create file")
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return nil, errors.Wrap(err, "fs.bucket.Add failed to write file content")
	}

	stats, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "fs.bucket.Add failed to get file informations")
	}

	return &brazier.Item{
		ID:        name,
		Data:      data,
		MimeType:  mimeType,
		CreatedAt: stats.ModTime(),
	}, nil
}

// Get an item by id
func (b *bucket) Get(id string) (*brazier.Item, error) {
	path := filepath.Join(b.path, id)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return nil, store.ErrNotFound
		}
		return nil, errors.Wrap(err, "fs.bucket.Get failed to read file content")
	}

	f, err := os.Open(filepath.Join(b.path, id))
	if err != nil {
		return nil, errors.Wrap(err, "fs.bucket.Get failed to open file")
	}

	defer f.Close()

	stats, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "fs.bucket.Get failed to get file informations")
	}

	return &brazier.Item{
		ID:        id,
		Data:      data,
		MimeType:  mime.TypeByExtension(filepath.Ext(id)),
		CreatedAt: stats.ModTime(),
	}, nil
}

// Delete item from the bucket
func (b *bucket) Delete(id string) error {
	err := os.Remove(filepath.Join(b.path, id))
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") {
			return store.ErrNotFound
		}
		return errors.Wrap(err, "fs.bucket.Delete failed to remove file")
	}

	return nil
}

// Close the session of the bucket
func (b *bucket) Close() error {
	return nil
}
