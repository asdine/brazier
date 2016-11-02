package mock_test

import (
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func TestRegistry(t *testing.T) {
	r := mock.NewRegistry(mock.NewBackend())

	err := r.Create()
	require.Equal(t, store.ErrForbidden, err)

	err = r.Create("a", "b", "c")
	require.NoError(t, err)

	err = r.Create("a", "b", "c")
	require.Equal(t, store.ErrAlreadyExists, err)

	err = r.Create("a", "b")
	require.Equal(t, store.ErrAlreadyExists, err)

	err = r.Create("a", "b", "c", "d")
	require.NoError(t, err)

	err = r.Create("a", "f", "g", "h")
	require.NoError(t, err)

	_, err = r.Bucket("a")
	require.NoError(t, err)

	_, err = r.Bucket()
	require.Equal(t, store.ErrForbidden, err)

	_, err = r.Bucket("a", "b")
	require.NoError(t, err)

	_, err = r.Bucket("a", "b", "k")
	require.Equal(t, store.ErrNotFound, err)
}
