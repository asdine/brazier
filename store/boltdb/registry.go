package boltdb

import (
	"strings"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/brazier/store/boltdb/internal"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/asdine/storm/q"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// NewRegistry returns a BoltDB Registry.
func NewRegistry(path string, b brazier.Backend) (*Registry, error) {
	var err error

	db, err := storm.Open(
		path,
		storm.AutoIncrement(),
		storm.Codec(protobuf.Codec),
		storm.BoltOptions(0644, &bolt.Options{
			Timeout: time.Duration(50) * time.Millisecond,
		}),
	)

	if err != nil {
		return nil, errors.Wrap(err, "Can't open database")
	}

	return &Registry{
		DB:      db,
		Backend: b,
	}, nil
}

// Registry is a BoltDB registry.
type Registry struct {
	DB      *storm.DB
	Backend brazier.Backend
}

// Create a bucket in the registry.
func (r *Registry) Create(nodes ...string) error {
	tx, err := r.DB.Begin(true)
	if err != nil {
		return errors.Wrapf(err, "failed to create bucket at path %s", strings.Join(nodes, "/"))
	}
	defer tx.Rollback()

	var path string
	for i, node := range nodes {
		if path != "" {
			path += "/"
		}
		path += node
		err = tx.Save(&internal.Meta{
			Key: path,
		})

		if err != nil && err != storm.ErrAlreadyExists {
			return errors.Wrapf(err, "failed to create bucket at path %s", path)
		}

		// last node must not exist
		if err == storm.ErrAlreadyExists && i == len(nodes)-1 {
			return store.ErrAlreadyExists
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "failed to create bucket at path %s", strings.Join(nodes, "/"))
	}

	return nil
}

// Bucket returns the selected bucket from the Backend.
func (r *Registry) Bucket(nodes ...string) (brazier.Bucket, error) {
	var meta internal.Meta

	path := strings.Join(nodes, "/")

	err := r.DB.One("Key", path, &meta)
	if err == storm.ErrNotFound {
		return nil, store.ErrNotFound
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to fetch bucket at path %s", path)
	}

	return r.Backend.Bucket(nodes...)
}

// Children buckets of the specified path.
func (r *Registry) Children(nodes ...string) ([]brazier.Item, error) {
	var metas []internal.Meta

	path := strings.Join(nodes, "/")

	err := r.DB.Select(
		q.NewFieldMatcher(
			"Key",
			&childMatcher{prefix: path},
		),
	).Find(&metas)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}

		return nil, errors.Wrapf(err, "failed to fetch bucket children at path %s", path)
	}

	tree := keyTree{
		children: make(map[string]*keyTree),
	}
	for _, m := range metas {
		key := strings.TrimPrefix(m.Key, path)
		if key == "" {
			continue
		}
		key = strings.TrimPrefix(key, "/")

		new := childrenToTree(tree.children, strings.Split(key, "/")...)
		if new != nil {
			tree.index = append(tree.index, new)
		}
	}

	return treeToItems(&tree), nil
}

type keyTree struct {
	key      string
	children map[string]*keyTree
	index    []*keyTree
}

func childrenToTree(tree map[string]*keyTree, nodes ...string) *keyTree {
	if len(nodes) == 0 {
		return nil
	}

	var fresh bool

	t, ok := tree[nodes[0]]
	if !ok {
		fresh = true
		t = &keyTree{
			key: nodes[0],
		}
		tree[nodes[0]] = t
	}

	if t.children == nil {
		t.children = make(map[string]*keyTree)
	}

	if len(nodes) > 1 {
		new := childrenToTree(t.children, nodes[1:len(nodes)]...)
		if new != nil {
			t.index = append(t.index, new)
		}
	}

	if fresh {
		return t
	}

	return nil
}

func treeToItems(tree *keyTree) []brazier.Item {
	items := make([]brazier.Item, len(tree.index))

	for i, t := range tree.index {
		items[i].Key = t.key
		if t.children != nil {
			items[i].Children = treeToItems(t)
		}
	}

	return items
}

// Close BoltDB connection
func (r *Registry) Close() error {
	err := r.Backend.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close backend")
	}

	err = r.DB.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close registry")
	}

	return nil
}

type childMatcher struct {
	prefix string
}

func (c *childMatcher) MatchField(v interface{}) (bool, error) {
	key, ok := v.(string)
	if !ok {
		return false, nil
	}

	if !strings.HasPrefix(key, c.prefix) {
		return false, nil
	}

	return len(key) >= len(c.prefix), nil
}
