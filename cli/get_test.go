package cli

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	app, cleanup := testableApp(t)
	defer cleanup()

	out := app.Out.(*bytes.Buffer)

	s := saveCmd{App: app}
	g := getCmd{App: app}

	tests := map[string][]string{
		"\"abc\"\n":                  []string{"bucket", "string", "abc"},
		"\"bcd\"\n":                  []string{"bucket", "json string", "\"bcd\""},
		"10\n":                       []string{"bucket", "number", "10"},
		"{\"a\":\"b\"}\n":            []string{"bucket", "object", `{"a": "b"}`},
		"[\"a\",10,{\"c\":\"d\"}]\n": []string{"bucket", "array", `["a", 10, {"c": "d"}]`},
	}

	for expected, cmds := range tests {
		err := s.Save(nil, cmds)
		require.NoError(t, err)
		out.Reset()
		err = g.Get(nil, cmds[:2])
		require.NoError(t, err)
		require.Equal(t, expected, out.String())
		out.Reset()
	}
}
