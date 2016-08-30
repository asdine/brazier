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

func TestCreateItemValidJSON(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(` {    " the  key" :   [ 1, "hi" , 45.6    ] }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	bucket, err := s.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	require.True(t, b.SaveInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`{" the  key":[1,"hi",45.6]}`), item.Data)
}

func TestCreateItemInvalidJSON(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(`my value`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	bucket, err := s.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	require.True(t, b.SaveInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`"my value"`), item.Data)
}

func TestGetItem(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	err := s.Create("a")
	require.NoError(t, err)

	bucket, err := s.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	item, err := b.Save("b", []byte(`"my value"`))
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/a/b", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.Equal(t, item.Data, w.Body.Bytes())

	require.True(t, b.GetInvoked)
}

func TestBadRequests(t *testing.T) {
	var h brazierHttp.Handler

	s := mock.NewStore()
	h.Store = s

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/a/b", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/a/b/c", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/a/b", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
