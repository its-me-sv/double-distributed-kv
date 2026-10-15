package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func GetCommand(str store.Store, args ...string) error {
	if len(args) != 1 {
		return errors.New("wrong number of arguments for 'get' command")
	}

	value, _ := str.Get(args[0])

	fmt.Fprintln(os.Stdout, value)

	return nil
}
