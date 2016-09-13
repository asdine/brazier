package cli

import (
	"bytes"
	"testing"

	"github.com/asdine/brazier/store"
	"github.com/stretchr/testify/require"
)

func TestDelete(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	out := app.Out.(*bytes.Buffer)

	s := saveCmd{App: app}
	d := deleteCmd{App: app}

	err := s.Save(nil, []string{"my bucket", "my key", "my value"})
	require.NoError(t, err)
	out.Reset()

	err = d.Delete(nil, []string{"my bucket", "my key"})
	require.NoError(t, err)
	require.Equal(t, "Item \"my key\" successfully deleted.\n", out.String())

	err = d.Delete(nil, []string{"my bucket", "my key"})
	require.Error(t, err)
	require.Equal(t, store.ErrNotFound, err)
}
