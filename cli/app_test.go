package cli

import (
	"bytes"
	"testing"

	"github.com/asdine/brazier/mock"
)

func testableApp(t *testing.T) *app {
	return &app{Out: bytes.NewBuffer([]byte("")), Store: mock.NewStore()}
}
