package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	brazierHttp "github.com/asdine/brazier/http"
	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func TestCreateItemInvalid(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(nil)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateItemValidJSON(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(` {    " the  key" :   [ 1, "hi" , 45.6    ] }`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	bucket, err := registry.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	require.True(t, b.SaveInvoked)
	require.True(t, b.CloseInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`{" the  key":[1,"hi",45.6]}`), item.Data)
}

func TestCreateItemInvalidJSON(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("PUT", "/a/b", bytes.NewReader([]byte(`my value`)))
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	bucket, err := registry.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	require.True(t, b.SaveInvoked)
	require.True(t, b.CloseInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`"my value"`), item.Data)
}

func TestGetItem(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	err := registry.Create("a")
	require.NoError(t, err)

	bucket, err := registry.Bucket("a")
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
	require.True(t, b.CloseInvoked)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/a/c", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteItem(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	err := registry.Create("a")
	require.NoError(t, err)

	bucket, err := registry.Bucket("a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	_, err = b.Save("b", []byte(`"my value"`))
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("DELETE", "/a/b", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, b.DeleteInvoked)
	require.True(t, b.CloseInvoked)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("DELETE", "/a/b", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("DELETE", "/b/a", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestListItems(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	err := registry.Create("1a")
	require.NoError(t, err)
	bucket, err := registry.Bucket("1a")
	require.NoError(t, err)
	b := bucket.(*mock.Bucket)

	for i := 0; i < 20; i++ {
		_, err = b.Save(fmt.Sprintf("id%d", i), []byte(`"my value"`))
		require.NoError(t, err)
	}

	err = h.Store.CreateBucket("/1a/id20/")
	require.NoError(t, err)

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/1a/", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.True(t, b.PageInvoked)
	require.True(t, b.CloseInvoked)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var list []interface{}
	err = json.Unmarshal(w.Body.Bytes(), &list)
	require.NoError(t, err)

	for i := range list {
		item := list[i].(map[string]interface{})
		require.Equal(t, fmt.Sprintf("id%d", i), item["key"])
	}

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/z", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestTree(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	s := store.NewStore(registry)
	h.Store = s

	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			item, err := s.Put(fmt.Sprintf("/a/b%d/k%d", i, j), []byte(`"Value"`))
			require.NoError(t, err)
			require.NotNil(t, item)
		}
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/a/?recursive=true", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/z?recursive=true", nil)
	h.ServeHTTP(w, r)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestBadRequests(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

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

func TestNestedBuckets(t *testing.T) {
	var h brazierHttp.Handler

	registry := mock.NewRegistry(mock.NewBackend())
	h.Store = store.NewStore(registry)

	var body = []byte(` {    " the  key" :   [ 1, "hi" , 45.6    ] }`)
	var expectedBody = []byte(`{" the  key":[1,"hi",45.6]}`)

	t.Run("put", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("PUT", "/a/b/c/d", bytes.NewReader(body))
		h.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)

		bucket, err := registry.Bucket("a", "b", "c")
		require.NoError(t, err)
		b := bucket.(*mock.Bucket)

		require.True(t, b.SaveInvoked)
		require.True(t, b.CloseInvoked)
		item, err := b.Get("d")
		require.NoError(t, err)
		require.Equal(t, expectedBody, item.Data)
	})

	t.Run("get", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/a/b/c/d", nil)
		h.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Equal(t, expectedBody, w.Body.Bytes())
	})

	t.Run("list", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/a/b/c/", nil)
		h.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
		require.Equal(t, "application/json", w.Header().Get("Content-Type"))
		require.Equal(t, `[{"key":"d","value":`+string(expectedBody)+`}]`, w.Body.String())
	})

	t.Run("delete", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("DELETE", "/a/b/c/d", nil)
		h.ServeHTTP(w, r)
		require.Equal(t, http.StatusOK, w.Code)
	})
}
