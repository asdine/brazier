package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	app := testableApp(t)

	out := app.Out.(*bytes.Buffer)

	c := createCmd{App: app}

	err := c.Create(nil, nil)
	require.Error(t, err)
	require.EqualError(t, err, "Bucket name is missing")

	err = c.Create(nil, []string{"my bucket"})
	require.NoError(t, err)
	require.Equal(t, "Bucket \"my bucket\" successfully created.\n", out.String())
}
