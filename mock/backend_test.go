package mock_test

import (
	"testing"

	"github.com/asdine/brazier/mock"
	"github.com/stretchr/testify/require"
)

func TestBackend(t *testing.T) {
	bck := mock.NewBackend()

	c, err := bck.Bucket("a", "b", "c")
	require.NoError(t, err)
	require.Equal(t, "c", c.(*mock.Bucket).Name)

	a, err := bck.Bucket("a")
	require.NoError(t, err)
	require.Equal(t, "a", a.(*mock.Bucket).Name)

	require.Len(t, bck.Buckets, 1)
}
