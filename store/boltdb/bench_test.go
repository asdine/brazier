package boltdb_test

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"

	"github.com/asdine/brazier/store/boltdb"
)

func BenchmarkBucketSave(b *testing.B) {
	path, cleanup := preparePath(b, "store.db")
	defer cleanup()

	s, err := boltdb.NewBackend(path)
	if err != nil {
		b.Error(err)
	}

	bucket, err := s.Bucket("b1")
	if err != nil {
		b.Error(err)
	}
	defer bucket.Close()

	val := bytes.Repeat([]byte("a"), 64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bucket.Save(fmt.Sprintf("id%d", i), val)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBucketGet(b *testing.B) {
	path, cleanup := preparePath(b, "store.db")
	defer cleanup()

	s, err := boltdb.NewBackend(path)
	if err != nil {
		b.Error(err)
	}

	bucket, err := s.Bucket("b1")
	if err != nil {
		b.Error(err)
	}
	defer bucket.Close()

	val := bytes.Repeat([]byte("a"), 64)
	_, err = bucket.Save("id", val)
	if err != nil {
		b.Error(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bucket.Get("id")
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkBucketPage(b *testing.B) {
	path, cleanup := preparePath(b, "store.db")
	defer cleanup()

	s, err := boltdb.NewBackend(path)
	if err != nil {
		b.Error(err)
	}

	bucket, err := s.Bucket("b1")
	if err != nil {
		b.Error(err)
	}
	defer bucket.Close()

	val := bytes.Repeat([]byte("a"), 64)

	for i := 0; i < 100; i++ {
		_, err = bucket.Save("id"+strconv.Itoa(i), val)
		if err != nil {
			b.Error(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := bucket.Page(1, 100)
		if err != nil {
			b.Error(err)
		}
	}
}
