package cli

import (
	"bytes"
	"encoding/json"
	"strings"
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

func TestCliGetListItems(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	testGetListItems(t, app)
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

func TestCliRPCGetListItems(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testGetListItems(t, app)
}

func TestCliRPCDelete(t *testing.T) {
	app, cleanup := testableAppRPC(t)
	defer cleanup()

	testDelete(t, app)
}

func testCreate(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	c := NewCreateCmd(app)

	err := c.RunE(nil, nil)
	require.Error(t, err)
	require.EqualError(t, err, "Bucket name is missing")

	err = c.RunE(nil, []string{"my bucket/my other bucket/"})
	require.NoError(t, err)
	require.Equal(t, "Bucket \"my bucket/my other bucket/\" successfully created.\n", out.String())
}

func testSave(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewPutCmd(app)

	err := s.RunE(nil, nil)
	require.EqualError(t, err, "Wrong number of arguments")

	err = s.RunE(nil, []string{"my bucket/my key", "my value"})
	require.NoError(t, err)
	require.Equal(t, "Item \"my bucket/my key\" successfully saved.\n", out.String())
}

func testGet(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewPutCmd(app)
	g := NewGetCmd(app, false)

	tests := map[string][]string{
		"\"abc\"\n":                  []string{"checkJson/string", "abc"},
		"\"bcd\"\n":                  []string{"checkJson/json string", "\"bcd\""},
		"10\n":                       []string{"checkJson/number", "10"},
		"{\"a\":\"b\"}\n":            []string{"checkJson/object", `{"a": "b"}`},
		"[\"a\",10,{\"c\":\"d\"}]\n": []string{"checkJson/array", `["a", 10, {"c": "d"}]`},
	}

	for expected, cmds := range tests {
		err := s.RunE(nil, cmds)
		require.NoError(t, err)
		out.Reset()
		err = g.RunE(nil, cmds[:1])
		require.NoError(t, err)
		require.JSONEq(t, expected, out.String())
		out.Reset()
	}
}

func testGetListItems(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewPutCmd(app)
	g := NewGetCmd(app, false)
	gr := NewGetCmd(app, true)

	tests := []map[string][]string{
		{"\"abc\"": []string{"test/checkJson/string", "abc"}},
		{"\"bcd\"": []string{"test/checkJson/json string", "\"bcd\""}},
		{"10": []string{"test/checkJson/number", "10"}},
		{"{\"a\":\"b\"}": []string{"test/checkJson/object", `{"a": "b"}`}},
		{"[\"a\",10,{\"c\":\"d\"}]": []string{"test/checkJson/array", `["a", 10, {"c": "d"}]`}},
	}

	var expected []map[string]interface{}
	for _, test := range tests {
		for output, cmds := range test {
			err := s.RunE(nil, cmds)
			require.NoError(t, err)

			kv := make(map[string]interface{})
			kv["key"] = strings.TrimPrefix(cmds[0], "test/checkJson/")
			v := json.RawMessage(output)
			kv["value"] = &v
			expected = append(expected, kv)
		}
	}

	out.Reset()
	err := g.RunE(nil, []string{"test/checkJson/"})
	require.NoError(t, err)
	rawExpected, err := json.Marshal(expected)
	require.NoError(t, err)
	require.JSONEq(t, string(rawExpected), out.String())

	out.Reset()
	err = g.RunE(nil, []string{"some bucket/"})
	require.Error(t, err)

	out.Reset()
	err = gr.RunE(nil, []string{"test/"})
	require.NoError(t, err)
	var output []interface{}
	err = json.Unmarshal(out.Bytes(), &output)
	require.NoError(t, err)
	require.Len(t, output, 1)
	item := output[0].(map[string]interface{})
	require.Equal(t, "checkJson/", item["key"].(string))
	list := item["value"].([]interface{})
	require.Len(t, list, 5)
}

func testDelete(t *testing.T, app *app) {
	out := app.Out.(*bytes.Buffer)

	s := NewPutCmd(app)
	d := NewDeleteCmd(app)

	err := s.RunE(nil, []string{"a/b/c/d", "my value"})
	require.NoError(t, err)
	out.Reset()

	err = d.RunE(nil, []string{"a/b/c/d"})
	require.NoError(t, err)
	require.Equal(t, "Item \"a/b/c/d\" successfully deleted.\n", out.String())

	err = d.RunE(nil, []string{"a/b/c/d"})
	require.Error(t, err)
}
