package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func ExistsCommand(str store.Store, args ...string) error {
	if len(args) == 0 {
		return errors.New("wrong number of arguments for 'exists' command")
	}

	keysFound := 0

	for _, key := range args {
		if str.Exists(key) {
			keysFound += 1
		}
	}

	fmt.Fprintln(os.Stdout, "(integer)", keysFound)

	return nil
}
