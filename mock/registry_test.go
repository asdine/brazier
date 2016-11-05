package mock_test

import (
	"fmt"
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

	items, err := r.Children("k")
	require.Equal(t, store.ErrNotFound, err)

	items, err = r.Children("a")
	require.NoError(t, err)
	require.Len(t, items, 2)
	require.Len(t, items[0].Children, 1)
	require.Len(t, items[1].Children, 1)

	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			err := r.Create("z", fmt.Sprintf("b%d", i), fmt.Sprintf("k%d", j))
			require.NoError(t, err)
		}
	}

	items, err = r.Children("z")
	require.NoError(t, err)
	require.Len(t, items, 3)
	require.Equal(t, "b0", items[0].Key)
	require.Len(t, items[0].Children, 5)
	require.Equal(t, "k0", items[0].Children[0].Key)
	require.Len(t, items[1].Children, 5)
}
