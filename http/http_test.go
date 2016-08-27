package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	brazierHttp "github.com/asdine/brazier/http"
	"github.com/asdine/brazier/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBucket(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(`{"name": "mybucket"}`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusCreated, w.Code)
	require.True(t, s.CreateInvoked)
	_, ok := s.Buckets["mybucket"]
	require.True(t, ok)
}

func TestCreateItem(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	err := s.Create("a")
	require.NoError(t, err)

	bucket, err := s.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(`my value`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, b.SaveInvoked)
	require.True(t, b.CloseInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`"my value"`), item.Data)
}
