package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSave(t *testing.T) {
	app := testableApp(t)

	out := app.Out.(*bytes.Buffer)

	s := saveCmd{App: app}

	err := s.Save(nil, nil)
	require.EqualError(t, err, "Wrong number of arguments")

	err = s.Save(nil, []string{"my bucket", "my key", "my value"})
	require.NoError(t, err)
	require.Equal(t, "Item \"my key\" successfully saved.\n", out.String())
}
