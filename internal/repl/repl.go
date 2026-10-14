package repl

import (
	"bufio"
	"fmt"
	"os"

	"github.com/its-me-sv/double-distributed-kv/internal/config"
	"github.com/its-me-sv/double-distributed-kv/internal/store"
)

func StartRepl(cfg config.Config, str store.Store) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, cfg.AppName, " > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()

		fmt.Println(input)

		if input == "exit" {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read input:", err)
	}
}
