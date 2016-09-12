package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListItems(t *testing.T) {
	app := testableApp(t)

	out := app.Out.(*bytes.Buffer)

	s := saveCmd{App: app}
	l := listCmd{App: app}

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
		err := s.Save(nil, cmds)
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
	err := l.List(nil, []string{"bucket"})
	require.NoError(t, err)
	require.Equal(t, expected.String(), out.String())
}

func TestListBuckets(t *testing.T) {
	app := testableApp(t)

	err := app.Store.Create("bucket1")
	require.NoError(t, err)
	err = app.Store.Create("bucket2")
	require.NoError(t, err)

	out := app.Out.(*bytes.Buffer)
	l := listCmd{App: app}
	err = l.List(nil, nil)
	require.NoError(t, err)
	require.Equal(t, "bucket1\nbucket2\n", out.String())
}
