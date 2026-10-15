package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func SetCommand(str store.Store, args ...string) error {
	if len(args) > 2 {
		return errors.New("syntax error")
	}

	if len(args) != 2 {
		return errors.New("wrong number of arguments for 'set' command")
	}

	str.Set(args[0], args[1])
	fmt.Fprintln(os.Stdout, "OK")

	return nil
}
