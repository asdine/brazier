package boltdb

import (
	"sync"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/store"
	"github.com/asdine/storm"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// NewStore returns a BoltDB store
func NewStore(path string) *Store {
	return &Store{
		Path:     path,
		sessions: make(map[string][]storm.Node),
	}
}

// Store is a BoltDB store
type Store struct {
	sync.Mutex
	DB       *storm.DB
	Path     string
	sessions map[string][]storm.Node
}

// Create a bucket
func (s *Store) Create(name string) error {
	if len(s.sessions) == 0 {
		err := s.open()
		if err != nil {
			return err
		}
		defer s.close()
	}

	return s.DB.Set("buckets", name, nil)
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(name string) (brazier.Bucket, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.sessions) == 0 {
		err := s.open()
		if err != nil {
			return nil, err
		}
	}

	var str []byte
	err := s.DB.Get("buckets", name, &str)
	if err != nil {
		if err == storm.ErrNotFound {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	node := s.DB.From(name)
	s.sessions[name] = append(s.sessions[name], node)

	return NewBucket(s, name, node), nil
}

// List returns the list of all buckets
func (s *Store) List() ([]string, error) {
	if len(s.sessions) == 0 {
		err := s.open()
		if err != nil {
			return nil, err
		}
		defer s.close()
	}

	var buckets []string
	err := s.DB.Select().Bucket("buckets").RawEach(func(k, v []byte) error {
		buckets = append(buckets, string(k))
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "boltdb.store.List failed to fetch buckets")
	}

	return buckets, nil
}

func (s *Store) open() error {
	var err error

	if s.DB != nil {
		return nil
	}
	s.DB, err = storm.Open(
		s.Path,
		storm.AutoIncrement(),
		storm.Codec(new(rawCodec)),
		storm.BoltOptions(0644, &bolt.Options{
			Timeout: time.Duration(50) * time.Millisecond,
		}),
	)
	return err
}

// Close BoltDB connection
func (s *Store) Close() error {
	s.Lock()
	defer s.Unlock()
	return s.close()
}

func (s *Store) close() error {
	var err error

	s.sessions = make(map[string][]storm.Node)
	if s.DB != nil {
		err = s.DB.Close()
		s.DB = nil
	}

	return err
}

func (s *Store) closeSession(name string) error {
	s.Lock()
	defer s.Unlock()

	list, ok := s.sessions[name]
	if !ok {
		return errors.New("unknown session id")
	}

	if len(list) == 1 {
		delete(s.sessions, name)
		if len(s.sessions) == 0 {
			return s.close()
		}
	} else {
		s.sessions[name] = list[:len(list)-1]
	}

	return nil
}
