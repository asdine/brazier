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

func TestHandler(t *testing.T) {
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
