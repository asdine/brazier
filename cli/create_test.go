package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	out := app.Out.(*bytes.Buffer)

	c := createCmd{App: app}

	err := c.Create(nil, nil)
	assert.Error(t, err)
	assert.EqualError(t, err, "Bucket name is missing")

	err = c.Create(nil, []string{"my bucket"})
	assert.NoError(t, err)
	assert.Equal(t, "Bucket \"my bucket\" successfully created.\n", out.String())

	err = c.Create(nil, []string{"my bucket"})
	assert.Error(t, err)
	assert.EqualError(t, err, "The bucket \"my bucket\" already exists.\n")
}
