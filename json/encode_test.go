package json_test

import (
	"testing"

	"github.com/asdine/brazier"
	"github.com/asdine/brazier/json"
	"github.com/stretchr/testify/require"
)

func TestMarshalList(t *testing.T) {
	items := []brazier.Item{
		{Key: "k1", Data: []byte(`"Data1"`)},
		{Key: "k2", Data: []byte(`"Data2"`)},
		{Key: "k3", Data: []byte(`"Data3"`)},
	}

	expected := `[{"data":"Data1","key":"k1"},{"data":"Data2","key":"k2"},{"data":"Data3","key":"k3"}]`
	out, err := json.MarshalList(items)
	require.NoError(t, err)
	require.Equal(t, expected, string(out))
}

func TestToValidJSON(t *testing.T) {
	tests := map[string]string{
		`invalid è`:                    `"invalid è"`,
		`{"invalid": "json"`:           `"{\"invalid\": \"json\""`,
		`"valid"`:                      `"valid"`,
		`{"dirty"      :   "json" }  `: `{"dirty":"json"}`,
		`5`: `5`,
	}

	for in, out := range tests {
		res := json.ToValidJSON([]byte(in))
		require.Equal(t, out, string(res))
	}
}

func BenchmarkToValidJSON(b *testing.B) {
	invalidJSON := []byte(`

		{
								"the name"       : "  &éà"      , "another     key"   : [ 1,  		10,9, "    str  " ]


		`)

	for i := 0; i < b.N; i++ {
		json.ToValidJSON(invalidJSON)
	}
}
