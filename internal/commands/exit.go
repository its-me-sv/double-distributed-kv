package commands

import (
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func ExitCommand(str store.Store, args ...string) error {
	os.Exit(0)
	return nil
}
