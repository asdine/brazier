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

func TestCreateItem(t *testing.T) {
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
	require.True(t, b.CloseInvoked)
	item, err := b.Get("b")
	require.NoError(t, err)
	require.Equal(t, []byte(`"my value"`), item.Data)
}
