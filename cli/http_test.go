package cli

import (
	"testing"

	"github.com/asdine/brazier"
	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	app := testableApp(t)
	app.Config.HTTP.Port = 55898

	h := httpCmd{App: app, ServeFunc: func(s brazier.Store, port int) error {
		require.Equal(t, 55898, port)
		return nil
	}}

	err := h.Serve(nil, nil)
	require.NoError(t, err)
}
