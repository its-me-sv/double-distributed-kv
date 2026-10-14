package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func SetCommand(str store.Store, args ...string) error {
	if len(args) != 2 {
		return errors.New("wrong number of arguments for 'set' command")
	}

	wasUpdated := str.Set(args[0], args[1])

	message := "INSERT"
	if wasUpdated {
		message = "UPSERT"
	}
	fmt.Fprintln(os.Stdout, message)

	return nil
}
