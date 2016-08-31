package boltdb

import (
	"errors"
	"sync"
	"time"

	"github.com/asdine/brazier"
	"github.com/asdine/storm"
	"github.com/asdine/storm/codec/protobuf"
	"github.com/boltdb/bolt"
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

// Create a bucket and return its informations
func (s *Store) Create(id string) error {
	bucket, err := s.Bucket(id)
	if err != nil {
		return err
	}
	defer bucket.Close()

	b := bucket.(*Bucket)
	return b.node.Init(&Item{})
}

// Bucket returns the bucket associated with the given id
func (s *Store) Bucket(id string) (brazier.Bucket, error) {
	s.Lock()
	defer s.Unlock()

	if len(s.sessions) == 0 {
		err := s.open()
		if err != nil {
			return nil, err
		}
	}

	node := s.DB.From(id)
	s.sessions[id] = append(s.sessions[id], node)

	return NewBucket(s, id, node), nil
}

func (s *Store) open() error {
	var err error

	if s.DB != nil {
		return nil
	}
	s.DB, err = storm.Open(
		s.Path,
		storm.AutoIncrement(),
		storm.Codec(protobuf.Codec),
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

func (s *Store) closeSession(id string) error {
	s.Lock()
	defer s.Unlock()

	list, ok := s.sessions[id]
	if !ok {
		return errors.New("unknown session id")
	}

	if len(list) == 1 {
		delete(s.sessions, id)
		if len(s.sessions) == 0 {
			return s.close()
		}
	} else {
		s.sessions[id] = list[:len(list)-1]
	}

	return nil
}
