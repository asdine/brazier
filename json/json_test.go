package json_test

import (
	"strings"
	"testing"

	"github.com/asdine/brazier/json"
	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	require.False(t, json.IsValid([]byte("")))
	require.False(t, json.IsValid([]byte("   ")))
	require.True(t, json.IsValid([]byte(`"string"`)))
	require.True(t, json.IsValid([]byte(`10.6`)))
	require.True(t, json.IsValid([]byte(`{"key": "value"}`)))
	require.False(t, json.IsValid([]byte(`{"bad": "format"`)))
	require.False(t, json.IsValid([]byte(`something else`)))
}

func TestIsValidFromReader(t *testing.T) {
	ok, data := json.IsValidReader(strings.NewReader(`"string"`))
	require.True(t, ok)
	require.Equal(t, `"string"`, string(data))
	ok, data = json.IsValidReader(strings.NewReader(`10.6`))
	require.True(t, ok)
	require.Equal(t, `10.6`, string(data))
	ok, data = json.IsValidReader(strings.NewReader(`{"key": "value"}`))
	require.True(t, ok)
	require.Equal(t, `{"key": "value"}`, string(data))
	ok, data = json.IsValidReader(strings.NewReader(`{"bad": "format"`))
	require.False(t, ok)
	require.Nil(t, data)
	ok, data = json.IsValidReader(strings.NewReader(`"something else`))
	require.False(t, ok)
	require.Nil(t, data)
}

func TestClean(t *testing.T) {
	require.Equal(t, []byte(``), json.Clean([]byte(``)))
	require.Equal(t, []byte(`"a b    c"`), json.Clean([]byte(`"a b    c"`)))
	require.Equal(t, []byte(`"a b    c"`), json.Clean([]byte(`   "a b    c"  `)))
	require.Equal(t, []byte(`10`), json.Clean([]byte("10\n")))
	require.Equal(t, []byte(`{"the name":"  &éà","another     key":[1,10,9,"    str  "]}`), json.Clean([]byte(`

		{
								"the name"       : "  &éà"      , "another     key"   : [ 1,  		10,9, "    str  " ]   }


		`)))
}

func BenchmarkClean(b *testing.B) {
	data := []byte(`

		{
								"the name"       : "  &éà"      , "another     key"   : [ 1,  		10,9, "    str  " ]   }


		`)

	for i := 0; i < b.N; i++ {
		json.Clean(data)
	}
}
