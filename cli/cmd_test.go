package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCliCreate(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testCreate(t, app)
}

func TestCliSave(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testSave(t, app)
}

func TestCliGet(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testGet(t, app)
}

func TestCliListItems(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testListItems(t, app)
}

func TestCliListBuckets(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testListBuckets(t, app)
}

func TestCliDelete(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testDelete(t, app)
}

func TestCliRPCCreate(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testCreate(t, app)
}

func TestCliRPCSave(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testSave(t, app)
}

func TestCliRPCGet(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testGet(t, app)
}

func TestCliRPCListItems(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testListItems(t, app)
}

func TestCliRPCListBuckets(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testListBuckets(t, app)
}

func TestCliRPCDelete(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testDelete(t, app)
}

func TestCliUse(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	c := NewUseCmd(app)

	err := c.RunE(nil, nil)
	require.Error(t, err)
	require.EqualError(t, err, "Bucket name is missing")

	err = c.RunE(nil, []string{"my bucket"})
	require.Error(t, err)
	require.EqualError(t, err, "Bucket \"my bucket\" not found.\n")

	err = app.Store.Create("my bucket")
	require.NoError(t, err)

	err = c.RunE(nil, []string{"my bucket"})
	require.NoError(t, err)

	db, err := app.settingsDB()
	defer db.Close()
	var to string
	err = db.Get("buckets", "default", &to)
	require.NoError(t, err)
	require.Equal(t, "my bucket", to)
}

func testCreate(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	c := NewCreateCmd(app)

	err := c.RunE(nil, nil)
	require.Error(t, err)
	require.EqualError(t, err, "Bucket name is missing")

	err = c.RunE(nil, []string{"my bucket"})
	require.NoError(t, err)
	require.Equal(t, "Bucket \"my bucket\" successfully created.\n", out.String())

	db, err := app.settingsDB()
	defer db.Close()
	var to string
	err = db.Get("buckets", "default", &to)
	require.NoError(t, err)
	require.Equal(t, "my bucket", to)
}

func testSave(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewSaveCmd(app)

	err := s.RunE(nil, nil)
	require.EqualError(t, err, "Wrong number of arguments")

	err = s.RunE(nil, []string{"my bucket", "my key", "my value"})
	require.NoError(t, err)
	require.Equal(t, "Item \"my key\" successfully saved.\n", out.String())
}

func testGet(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewSaveCmd(app)
	g := NewGetCmd(app)

	tests := map[string][]string{
		"\"abc\"\n":                  []string{"bucket", "string", "abc"},
		"\"bcd\"\n":                  []string{"bucket", "json string", "\"bcd\""},
		"10\n":                       []string{"bucket", "number", "10"},
		"{\"a\":\"b\"}\n":            []string{"bucket", "object", `{"a": "b"}`},
		"[\"a\",10,{\"c\":\"d\"}]\n": []string{"bucket", "array", `["a", 10, {"c": "d"}]`},
	}

	for expected, cmds := range tests {
		err := s.RunE(nil, cmds)
		require.NoError(t, err)
		out.Reset()
		err = g.RunE(nil, cmds[:2])
		require.NoError(t, err)
		require.Equal(t, expected, out.String())
		out.Reset()
	}
}

func testListItems(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewSaveCmd(app)
	l := NewListCmd(app)

	tests := map[string][]string{
		"\"abc\"":                  []string{"bucket", "string", "abc"},
		"\"bcd\"":                  []string{"bucket", "json string", "\"bcd\""},
		"10":                       []string{"bucket", "number", "10"},
		"{\"a\":\"b\"}":            []string{"bucket", "object", `{"a": "b"}`},
		"[\"a\",10,{\"c\":\"d\"}]": []string{"bucket", "array", `["a", 10, {"c": "d"}]`},
	}

	var expected bytes.Buffer

	expected.WriteByte('[')
	first := true
	for output, cmds := range tests {
		err := s.RunE(nil, cmds)
		require.NoError(t, err)
		if !first {
			expected.WriteByte(',')
		} else {
			first = false
		}
		expected.WriteString(`{"data":`)
		expected.WriteString(output)
		expected.WriteString(`,"key":"`)
		expected.WriteString(cmds[1])
		expected.WriteString(`"}`)
	}
	expected.WriteString("]\n")

	out.Reset()
	err := l.RunE(nil, []string{"bucket"})
	require.NoError(t, err)
	require.Equal(t, expected.String(), out.String())
}

func testListBuckets(t *testing.T, app *app) {
	err := app.Store.Create("bucket1")
	require.NoError(t, err)
	err = app.Store.Create("bucket2")
	require.NoError(t, err)

	out := app.Out.(*bytes.Buffer)
	l := NewListCmd(app)
	err = l.RunE(nil, nil)
	require.NoError(t, err)
	require.Equal(t, "bucket1\nbucket2\n", out.String())
}

func testDelete(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewSaveCmd(app)
	d := NewDeleteCmd(app)

	err := s.RunE(nil, []string{"my bucket", "my key", "my value"})
	require.NoError(t, err)
	out.Reset()

	err = d.RunE(nil, []string{"my bucket", "my key"})
	require.NoError(t, err)
	require.Equal(t, "Item \"my key\" successfully deleted.\n", out.String())

	err = d.RunE(nil, []string{"my bucket", "my key"})
	require.Error(t, err)
}
