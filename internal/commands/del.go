package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func DelCommand(str store.Store, args ...string) error {
	if len(args) == 0 {
		return errors.New("wrong number of arguments for 'del' command")
	}

	deletedKeysCount := 0

	for _, key := range args {
		if str.Del(key) {
			deletedKeysCount += 1
		}
	}

	fmt.Fprintln(os.Stdout, "(integer)", deletedKeysCount)

	return nil
}
