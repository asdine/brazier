package cli

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func testableApp(t *testing.T) (*app, func()) {
	dir, err := ioutil.TempDir(os.TempDir(), "brazier")
	if err != nil {
		t.Error(err)
	}

	a := app{Out: &bytes.Buffer{}, Path: path.Join(dir, "brazier.db")}
	return &a, func() {
		os.RemoveAll(dir)
	}
}
