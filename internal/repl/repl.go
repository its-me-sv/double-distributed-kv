package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/its-me-sv/double-distributed-kv/internal/commands"
	"github.com/its-me-sv/double-distributed-kv/internal/config"
	"github.com/its-me-sv/double-distributed-kv/internal/store"
	"github.com/its-me-sv/double-distributed-kv/internal/utils"
)

func StartRepl(cfg config.Config, str store.Store) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Fprint(os.Stdout, cfg.AppName, "> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		cleanedInput, err := utils.Tokenize(input)
		if err != nil {
			fmt.Fprintln(os.Stdout, "ERR", err.Error())
			continue
		}

		if len(cleanedInput) == 0 {
			continue
		}

		commandName := strings.ToLower(cleanedInput[0])
		arguments := []string{}

		if len(cleanedInput) > 1 {
			arguments = cleanedInput[1:]
		}

		if commandName == "debug" {
			str.Debug()
			continue
		}

		command, commandExists := getCommands()[commandName]
		if !commandExists {
			fmt.Fprintln(os.Stdout, "ERR unknown command")
			continue
		}

		if err := command(str, arguments...); err != nil {
			fmt.Fprintln(os.Stdout, "ERR", err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ERR failed to read input")
	}
}

type commandCallback = func(store.Store, ...string) error

func getCommands() map[string]commandCallback {
	return map[string]commandCallback{
		"exit":   commands.ExitCommand,
		"set":    commands.SetCommand,
		"get":    commands.GetCommand,
		"del":    commands.DelCommand,
		"exists": commands.ExistsCommand,
	}
}
