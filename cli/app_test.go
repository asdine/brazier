package cli

import (
	"os"
	"testing"

	"github.com/asdine/brazier/mock"
)

func testableApp(t *testing.T) *app {
	return &app{Out: os.Stdout, Store: mock.NewStore()}
}
